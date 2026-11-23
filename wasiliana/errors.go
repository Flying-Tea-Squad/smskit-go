package wasiliana

import (
	"encoding/json"
	"net/http"
	"strings"

	smskit "github.com/Flying-Tea-Squad/smskit-go"
)

// knownMessages are the provider messages documented at
// https://docs.wasiliana.com. They cannot contain user data, so they are safe
// to surface. Any other provider text is dropped from Message (it stays in Raw).
var knownMessages = map[string]bool{
	"Wrong sender id sent":                 true,
	"App does not match your api key":      true,
	"Short code does not exist":            true,
	"App does not exist":                   true,
	"ApiKey must be present in the header": true,
}

// toError translates a failed HTTP status or a "status":"failed" body.
func toError(httpStatus int, body []byte) error {
	e := &smskit.ProviderError{Provider: providerName, HTTPStatus: httpStatus}
	if json.Valid(body) {
		e.Raw = json.RawMessage(body)
	}
	switch {
	case httpStatus == http.StatusBadRequest || httpStatus == http.StatusNotAcceptable:
		e.Kind = smskit.ErrorInvalidRequest
	case httpStatus == http.StatusUnauthorized:
		e.Kind = smskit.ErrorAuthentication
	case httpStatus == http.StatusForbidden:
		e.Kind = smskit.ErrorAuthorization
	case httpStatus == http.StatusUnprocessableEntity:
		e.Kind = smskit.ErrorRejected // live API: unknown app / bad key answers 422
	case httpStatus == http.StatusTooManyRequests:
		e.Kind, e.Retryable = smskit.ErrorRateLimited, true
	case httpStatus == http.StatusInternalServerError || httpStatus == http.StatusServiceUnavailable:
		e.Kind, e.Retryable = smskit.ErrorTemporary, true
	case httpStatus/100 == 2:
		e.Kind = smskit.ErrorRejected // 2xx carrying "status":"failed"
	default:
		e.Kind = smskit.ErrorProvider
	}

	e.Message = http.StatusText(httpStatus)
	if httpStatus/100 == 2 {
		e.Message = "request rejected"
	}
	var r struct {
		Status  string `json:"status"`
		Data    string `json:"data"`
		Message string `json:"message"`
		Code    any    `json:"code"`
	}
	if json.Unmarshal(body, &r) == nil {
		for _, m := range []string{r.Data, r.Message} {
			m = strings.TrimSpace(m) // live API sends e.g. "App does not match your api key "
			if knownMessages[m] {
				e.Message = m
				break
			}
		}
	}
	return e
}
