package smskit

import (
	"fmt"
	"net/url"
	"strings"
)

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

// requireRecipient checks that the recipient is not empty or only spaces.
func requireRecipient(channel, to string) error {
	if strings.TrimSpace(to) == "" {
		return fmt.Errorf("%s recipient is required: %w", channel, ErrInvalidRequest)
	}
	return nil
}

// Validate makes sure a WhatsApp text has a recipient and a message.
func (m WhatsAppText) Validate() error {
	if err := requireRecipient("whatsapp text", m.To); err != nil {
		return err
	}
	if strings.TrimSpace(m.Body) == "" {
		return fmt.Errorf("whatsapp text body is required: %w", ErrInvalidRequest)
	}
	return nil
}

// Validate makes sure a WhatsApp template has a recipient, a name and a language.
// Parameters are optional, but none of them can be blank.
func (m WhatsAppTemplate) Validate() error {
	if err := requireRecipient("whatsapp template", m.To); err != nil {
		return err
	}
	if strings.TrimSpace(m.Name) == "" {
		return fmt.Errorf("whatsapp template name is required: %w", ErrInvalidRequest)
	}
	if strings.TrimSpace(m.Language) == "" {
		return fmt.Errorf("whatsapp template language is required: %w", ErrInvalidRequest)
	}
	for i, p := range m.Parameters {
		if strings.TrimSpace(p) == "" {
			return fmt.Errorf("whatsapp template parameter %d is empty: %w", i, ErrInvalidRequest)
		}
	}
	return nil
}

// Validate makes sure WhatsApp media has a recipient and a web link (http or https).
func (m WhatsAppMedia) Validate() error {
	if err := requireRecipient("whatsapp media", m.To); err != nil {
		return err
	}
	if strings.TrimSpace(m.URL) == "" {
		return fmt.Errorf("whatsapp media url is required: %w", ErrInvalidRequest)
	}
	u, err := url.Parse(m.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return fmt.Errorf("whatsapp media url must be a http or https link: %w", ErrInvalidRequest)
	}
	return nil
}
