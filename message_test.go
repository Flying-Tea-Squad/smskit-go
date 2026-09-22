package smskit

import (
	"errors"
	"testing"
)

func TestSMSMessageValidate(t *testing.T) {
	if err := (SMSMessage{Body: "hello"}).Validate(); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}
	if err := (SMSMessage{To: "+254700000000"}).Validate(); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected invalid request, got %v", err)
	}
	if err := (SMSMessage{To: "+254700000000", Body: "hello"}).Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}
