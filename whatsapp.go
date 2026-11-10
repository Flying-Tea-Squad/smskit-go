package smskit

import (
	"fmt"
	"strings"
)

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

// Validate makes sure WhatsApp media has a recipient and a link to the file.
func (m WhatsAppMedia) Validate() error {
	if err := requireRecipient("whatsapp media", m.To); err != nil {
		return err
	}
	if strings.TrimSpace(m.URL) == "" {
		return fmt.Errorf("whatsapp media url is required: %w", ErrInvalidRequest)
	}
	return nil
}
