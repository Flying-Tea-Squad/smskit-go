package smskit

import "encoding/json"

// MessageStatus is the conservative normalized status of a message.
type MessageStatus string

const (
	StatusUnknown   MessageStatus = "unknown"
	StatusAccepted  MessageStatus = "accepted"
	StatusQueued    MessageStatus = "queued"
	StatusSent      MessageStatus = "sent"
	StatusDelivered MessageStatus = "delivered"
	StatusRead      MessageStatus = "read"
	StatusFailed    MessageStatus = "failed"
)

// SendResult contains common send information and the provider's exact status.
type SendResult struct {
	MessageID      string
	Status         MessageStatus
	Provider       string
	ProviderStatus string
	Raw            json.RawMessage
}
