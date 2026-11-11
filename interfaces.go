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
