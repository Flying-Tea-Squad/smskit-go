package smskit

import (
	"errors"
	"testing"
)

func TestWhatsAppValidate(t *testing.T) {
	tests := []struct {
		name    string
		msg     interface{ Validate() error }
		wantErr bool
	}{
		{"text valid", WhatsAppText{To: "+254700000000", Body: "hi"}, false},
		{"text missing recipient", WhatsAppText{Body: "hi"}, true},
		{"text missing body", WhatsAppText{To: "+254700000000"}, true},
		{"text spaces only", WhatsAppText{To: "  ", Body: "\t"}, true},

		{"template valid", WhatsAppTemplate{To: "+254700000000", Name: "welcome", Language: "en"}, false},
		{"template valid with parameters", WhatsAppTemplate{To: "+254700000000", Name: "welcome", Language: "en", Parameters: []string{"Sharleen"}}, false},
		{"template missing recipient", WhatsAppTemplate{Name: "welcome", Language: "en"}, true},
		{"template missing name", WhatsAppTemplate{To: "+254700000000", Language: "en"}, true},
		{"template missing language", WhatsAppTemplate{To: "+254700000000", Name: "welcome"}, true},
		{"template spaces only name", WhatsAppTemplate{To: "+254700000000", Name: " ", Language: "en"}, true},
		{"template empty parameter", WhatsAppTemplate{To: "+254700000000", Name: "welcome", Language: "en", Parameters: []string{""}}, true},

		{"media valid", WhatsAppMedia{To: "+254700000000", URL: "https://example.com/a.png"}, false},
		{"media missing recipient", WhatsAppMedia{URL: "https://example.com/a.png"}, true},
		{"media missing url", WhatsAppMedia{To: "+254700000000"}, true},
		{"media spaces only url", WhatsAppMedia{To: "+254700000000", URL: "  "}, true},
		{"media not a url", WhatsAppMedia{To: "+254700000000", URL: "not a url"}, true},
		{"media file url", WhatsAppMedia{To: "+254700000000", URL: "file:///etc/passwd"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.Validate()
			if tt.wantErr && !errors.Is(err, ErrInvalidRequest) {
				t.Errorf("want ErrInvalidRequest, got %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
