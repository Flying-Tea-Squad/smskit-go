package smskit

import "encoding/json"

// MessageStatus is the conservative normalized status of a message.
type MessageStatus string

const (
	// StatusUnknown means the provider's status was not recognized.
	StatusUnknown MessageStatus = "unknown"
	// StatusAccepted means the provider accepted the message for processing.
	StatusAccepted MessageStatus = "accepted"
	// StatusQueued means the message is waiting to be sent.
	StatusQueued MessageStatus = "queued"
	// StatusSent means the provider has sent the message onward.
	StatusSent MessageStatus = "sent"
	// StatusDelivered means the provider confirmed delivery to the handset.
	StatusDelivered MessageStatus = "delivered"
	// StatusRead means the recipient has read the message.
	StatusRead MessageStatus = "read"
	// StatusFailed means the message could not be sent or delivered.
	StatusFailed MessageStatus = "failed"
)

// SendResult contains common send information and the provider's exact status.
type SendResult struct {
	MessageID      string
	Status         MessageStatus
	Provider       string
	ProviderStatus string
	Raw            json.RawMessage
}
