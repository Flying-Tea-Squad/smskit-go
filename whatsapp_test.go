package smskit

import (
	"errors"
	"testing"
)

func TestWhatsAppText(t *testing.T) {
	msg := WhatsAppText{To: "+254700000000", Body: "hi"}
	if err := msg.Validate(); err != nil {
		t.Error("a full message should pass")
	}

	msg = WhatsAppText{Body: "hi"}
	if !errors.Is(msg.Validate(), ErrInvalidRequest) {
		t.Error("a message with no recipient should fail")
	}

	msg = WhatsAppText{To: "+254700000000"}
	if !errors.Is(msg.Validate(), ErrInvalidRequest) {
		t.Error("a message with no body should fail")
	}
}

func TestWhatsAppTemplate(t *testing.T) {
	msg := WhatsAppTemplate{To: "+254700000000", Name: "welcome", Language: "en"}
	if err := msg.Validate(); err != nil {
		t.Error("a full template should pass")
	}

	msg = WhatsAppTemplate{Name: "welcome", Language: "en"}
	if !errors.Is(msg.Validate(), ErrInvalidRequest) {
		t.Error("a template with no recipient should fail")
	}

	msg = WhatsAppTemplate{To: "+254700000000", Language: "en"}
	if !errors.Is(msg.Validate(), ErrInvalidRequest) {
		t.Error("a template with no name should fail")
	}

	msg = WhatsAppTemplate{To: "+254700000000", Name: "welcome"}
	if !errors.Is(msg.Validate(), ErrInvalidRequest) {
		t.Error("a template with no language should fail")
	}
}

func TestWhatsAppMedia(t *testing.T) {
	msg := WhatsAppMedia{To: "+254700000000", URL: "https://example.com/a.png"}
	if err := msg.Validate(); err != nil {
		t.Error("full media should pass")
	}

	msg = WhatsAppMedia{URL: "https://example.com/a.png"}
	if !errors.Is(msg.Validate(), ErrInvalidRequest) {
		t.Error("media with no recipient should fail")
	}

	msg = WhatsAppMedia{To: "+254700000000"}
	if !errors.Is(msg.Validate(), ErrInvalidRequest) {
		t.Error("media with no link should fail")
	}
}
