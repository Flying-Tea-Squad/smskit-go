package wasiliana

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	smskit "github.com/Flying-Tea-Squad/smskit-go"
)

const testKey = "secret-key-123"

func fixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func newClient(t *testing.T, h http.HandlerFunc, mut func(*Config)) (*Client, *int32) {
	t.Helper()
	var n int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&n, 1)
		h(w, r)
	}))
	t.Cleanup(srv.Close)
	cfg := Config{APIKey: testKey, BaseURL: srv.URL}
	if mut != nil {
		mut(&cfg)
	}
	c, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return c, &n
}

func msg() smskit.SMSMessage {
	return smskit.SMSMessage{To: "254700000000", From: "company-x", Body: "hello"}
}

func TestNewRejectsEmptyKey(t *testing.T) {
	if _, err := New(Config{}); !errors.Is(err, smskit.ErrInvalidRequest) {
		t.Fatalf("err = %v, want ErrInvalidRequest", err)
	}
}

func TestNewDoesNotModifyHTTPClient(t *testing.T) {
	hc := &http.Client{Timeout: 7 * time.Second}
	if _, err := New(Config{APIKey: testKey, HTTPClient: hc}); err != nil {
		t.Fatal(err)
	}
	if hc.Timeout != 7*time.Second || hc.Transport != nil {
		t.Errorf("caller client modified: %+v", hc)
	}
}

func TestSendSuccessContract(t *testing.T) {
	var raw map[string]any
	c, _ := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/send/sms" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("apiKey") != testKey {
			t.Errorf("apiKey header = %q", r.Header.Get("apiKey"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("content type = %q", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("User-Agent") != "test-agent" {
			t.Errorf("user agent = %q", r.Header.Get("User-Agent"))
		}
		_ = json.NewDecoder(r.Body).Decode(&raw)
		_, _ = w.Write([]byte(fixture(t, "send_success.json")))
	}, func(c *Config) { c.UserAgent = "test-agent" })

	res, err := c.SendSMS(context.Background(), msg())
	if err != nil {
		t.Fatal(err)
	}
	rec, ok := raw["recipients"].([]any)
	if !ok || len(rec) != 1 || rec[0] != "254700000000" {
		t.Errorf("recipients = %#v", raw["recipients"])
	}
	if raw["from"] != "company-x" || raw["message"] != "hello" {
		t.Errorf("body = %#v", raw)
	}
	if res.Status != smskit.StatusAccepted || res.ProviderStatus != "success" ||
		res.Provider != "wasiliana" || res.MessageID != "" || len(res.Raw) == 0 {
		t.Errorf("result = %+v", res)
	}
}

func TestSendErrors(t *testing.T) {
	cases := []struct {
		name      string
		status    int
		body      string
		kind      smskit.ErrorKind
		retryable bool
		message   string
	}{
		{"failed sender id", 200, "send_failed_sender_id.json", smskit.ErrorRejected, false, "Wrong sender id sent"},
		{"failed app mismatch", 200, "send_failed_app.json", smskit.ErrorRejected, false, "App does not match your api key"},
		{"failed short code", 200, "send_failed_short_code.json", smskit.ErrorRejected, false, "Short code does not exist"},
		{"failed app missing", 200, "send_failed_app_missing.json", smskit.ErrorRejected, false, "App does not exist"},
		{"failed unknown text drops message", 200, `{"status":"failed","data":"x 254700000000 y"}`, smskit.ErrorRejected, false, "request rejected"},
		{"400", 400, `{}`, smskit.ErrorInvalidRequest, false, ""},
		{"406", 406, `{}`, smskit.ErrorInvalidRequest, false, ""},
		{"401", 401, "error_401.json", smskit.ErrorAuthentication, false, "ApiKey must be present in the header"},
		{"live 401 trailing space", 401, `{"message":"App does not match your api key ","code":401}`, smskit.ErrorAuthentication, false, "App does not match your api key"},
		{"live 422 bad key", 422, `{"message":"App does not exist","code":422}`, smskit.ErrorRejected, false, "App does not exist"},
		{"403", 403, `{}`, smskit.ErrorAuthorization, false, ""},
		{"404", 404, `{}`, smskit.ErrorProvider, false, ""},
		{"405", 405, `{}`, smskit.ErrorProvider, false, ""},
		{"429", 429, `{}`, smskit.ErrorRateLimited, true, ""},
		{"500", 500, `oops`, smskit.ErrorTemporary, true, ""},
		{"503", 503, ``, smskit.ErrorTemporary, true, ""},
		{"malformed json", 200, `{not json`, smskit.ErrorProvider, false, "malformed response"},
		{"empty body", 200, ``, smskit.ErrorProvider, false, "malformed response"},
		{"unknown status", 200, `{"status":"maybe"}`, smskit.ErrorProvider, false, "unrecognized response status"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body := tc.body
			if strings.HasSuffix(body, ".json") {
				body = fixture(t, body)
			}
			c, _ := newClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(body))
			}, nil)
			res, err := c.SendSMS(context.Background(), msg())
			var pe *smskit.ProviderError
			if res != nil || !errors.As(err, &pe) {
				t.Fatalf("res=%v err=%v", res, err)
			}
			if pe.Kind != tc.kind || pe.Retryable != tc.retryable || pe.Provider != "wasiliana" || pe.HTTPStatus != tc.status {
				t.Errorf("got kind=%s retry=%v http=%d", pe.Kind, pe.Retryable, pe.HTTPStatus)
			}
			if tc.message != "" && pe.Message != tc.message {
				t.Errorf("message = %q, want %q", pe.Message, tc.message)
			}
			for _, s := range []string{testKey, "254700000000", "hello"} {
				if strings.Contains(err.Error(), s) {
					t.Errorf("error leaks %q: %v", s, err)
				}
			}
		})
	}
}

func TestInvalidInputMakesNoRequest(t *testing.T) {
	c, n := newClient(t, func(w http.ResponseWriter, r *http.Request) {}, nil)
	for _, m := range []smskit.SMSMessage{
		{From: "x", Body: "b"},
		{To: "254700000000", From: "x"},
		{To: "254700000000", Body: "b"},
	} {
		if _, err := c.SendSMS(context.Background(), m); !errors.Is(err, smskit.ErrInvalidRequest) {
			t.Errorf("%+v: err = %v", m, err)
		}
	}
	if atomic.LoadInt32(n) != 0 {
		t.Errorf("%d requests sent", *n)
	}
}

func TestCancelledContext(t *testing.T) {
	c, _ := newClient(t, func(w http.ResponseWriter, r *http.Request) {}, nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := c.SendSMS(ctx, msg())
	var pe *smskit.ProviderError
	if !errors.As(err, &pe) || pe.Kind != smskit.ErrorTransport || pe.Retryable {
		t.Fatalf("err = %v", err)
	}
}

func TestTimeoutNotRetried(t *testing.T) {
	c, n := newClient(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
	}, func(c *Config) { c.HTTPClient = &http.Client{Timeout: 50 * time.Millisecond} })
	_, err := c.SendSMS(context.Background(), msg())
	var pe *smskit.ProviderError
	if !errors.As(err, &pe) || pe.Kind != smskit.ErrorTransport || pe.Retryable {
		t.Fatalf("err = %v", err)
	}
	if strings.Contains(err.Error(), testKey) {
		t.Errorf("error leaks key: %v", err)
	}
	if got := atomic.LoadInt32(n); got != 1 {
		t.Errorf("requests = %d, want exactly 1 (no retry)", got)
	}
}
