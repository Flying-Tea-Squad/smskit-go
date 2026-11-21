package smskit

import "testing"

func TestNormalizeStatus(t *testing.T) {
	mapping := map[string]MessageStatus{
		"Accepted": StatusAccepted,
		"queued":   StatusQueued,
		"sent":     StatusSent,
	}

	tests := []struct {
		name           string
		providerStatus string
		mapping        map[string]MessageStatus
		want           MessageStatus
	}{
		{name: "known status", providerStatus: "Accepted", mapping: mapping, want: StatusAccepted},
		{name: "another known status", providerStatus: "queued", mapping: mapping, want: StatusQueued},
		{name: "unrecognized status", providerStatus: "delivering", mapping: mapping, want: StatusUnknown},
		{name: "empty status", providerStatus: "", mapping: mapping, want: StatusUnknown},
		{name: "nil mapping", providerStatus: "sent", mapping: nil, want: StatusUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeStatus(tt.providerStatus, tt.mapping); got != tt.want {
				t.Fatalf("NormalizeStatus(%q) = %q, want %q", tt.providerStatus, got, tt.want)
			}
		})
	}
}
