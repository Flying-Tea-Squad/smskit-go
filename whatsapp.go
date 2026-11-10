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
		return fmt.Errorf("the message is empty: %w", ErrInvalidRequest)
	}
	return nil
}

// Validate makes sure a WhatsApp template has a recipient, a name and a language.
func (m WhatsAppTemplate) Validate() error {
	if err := requireRecipient("whatsapp template", m.To); err != nil {
		return err
	}
	if strings.TrimSpace(m.Name) == "" {
		return fmt.Errorf("the template name is empty: %w", ErrInvalidRequest)
	}
	if strings.TrimSpace(m.Language) == "" {
		return fmt.Errorf("the template language is empty: %w", ErrInvalidRequest)
	}
	return nil
}

// Validate makes sure WhatsApp media has a recipient and a link to the file.
func (m WhatsAppMedia) Validate() error {
	if err := requireRecipient("whatsapp media", m.To); err != nil {
		return err
	}
	if strings.TrimSpace(m.URL) == "" {
		return fmt.Errorf("the media link is empty: %w", ErrInvalidRequest)
	}
	return nil
}
