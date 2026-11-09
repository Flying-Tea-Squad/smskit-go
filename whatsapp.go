package smskit

import "fmt"

// Validate makes sure a WhatsApp text has a recipient and a message.
func (m WhatsAppText) Validate() error {
	if m.To == "" {
		return fmt.Errorf("who is this for? to is empty: %w", ErrInvalidRequest)
	}
	if m.Body == "" {
		return fmt.Errorf("the message is empty: %w", ErrInvalidRequest)
	}
	return nil
}

// Validate makes sure a WhatsApp template has a recipient, a name and a language.
func (m WhatsAppTemplate) Validate() error {
	if m.To == "" {
		return fmt.Errorf("who is this for? to is empty: %w", ErrInvalidRequest)
	}
	if m.Name == "" {
		return fmt.Errorf("the template name is empty: %w", ErrInvalidRequest)
	}
	if m.Language == "" {
		return fmt.Errorf("the template language is empty: %w", ErrInvalidRequest)
	}
	return nil
}

// Validate makes sure WhatsApp media has a recipient and a link to the file.
func (m WhatsAppMedia) Validate() error {
	if m.To == "" {
		return fmt.Errorf("who is this for? to is empty: %w", ErrInvalidRequest)
	}
	if m.URL == "" {
		return fmt.Errorf("the media link is empty: %w", ErrInvalidRequest)
	}
	return nil
}
