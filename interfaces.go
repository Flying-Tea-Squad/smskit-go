package smskit

import "context"

// SMSSender sends one SMS message.
type SMSSender interface {
	SendSMS(context.Context, SMSMessage) (*SendResult, error)
}

// WhatsAppTextSender sends a WhatsApp text message.
type WhatsAppTextSender interface {
	SendWhatsAppText(context.Context, WhatsAppText) (*SendResult, error)
}

// WhatsAppTemplateSender sends an approved WhatsApp template message.
type WhatsAppTemplateSender interface {
	SendWhatsAppTemplate(context.Context, WhatsAppTemplate) (*SendResult, error)
}

// WhatsAppMediaSender sends supported WhatsApp media.
type WhatsAppMediaSender interface {
	SendWhatsAppMedia(context.Context, WhatsAppMedia) (*SendResult, error)
}

// WhatsAppText is the provider-independent representation of a WhatsApp text.
type WhatsAppText struct {
	To   string
	Body string
}

// WhatsAppTemplate is the provider-independent representation of a template.
type WhatsAppTemplate struct {
	To         string
	Name       string
	Language   string
	Parameters []string
}

// WhatsAppMedia is the provider-independent representation of supported media.
type WhatsAppMedia struct {
	To          string
	URL         string
	ContentType string
	Caption     string
}
