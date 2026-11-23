package wasiliana

import (
	"context"
	"encoding/json"
	"fmt"

	smskit "github.com/Flying-Tea-Squad/smskit-go"
)

const sendPath = "/api/v1/send/sms"

var _ smskit.SMSSender = (*Client)(nil)

// SendSMS sends one SMS through a short code or sender ID. Wasiliana requires
// msg.From. The API returns no message ID, so SendResult.MessageID is empty.
// SendSMS never retries.
func (c *Client) SendSMS(ctx context.Context, msg smskit.SMSMessage) (*smskit.SendResult, error) {
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	if msg.From == "" {
		return nil, fmt.Errorf("wasiliana: sender (From) is required: %w", smskit.ErrInvalidRequest)
	}
	status, body, err := c.postJSON(ctx, sendPath, sendRequest{
		Recipients: []string{msg.To}, From: msg.From, Message: msg.Body,
	})
	if err != nil {
		return nil, err
	}
	return parseSendResponse(status, body)
}

func parseSendResponse(status int, body []byte) (*smskit.SendResult, error) {
	if status/100 != 2 {
		return nil, toError(status, body)
	}
	var r apiResponse
	if err := json.Unmarshal(body, &r); err != nil {
		// The message may have been dispatched: not retryable.
		return nil, &smskit.ProviderError{
			Provider: providerName, Kind: smskit.ErrorProvider,
			Message: "malformed response", HTTPStatus: status,
		}
	}
	switch r.Status {
	case "success":
		return &smskit.SendResult{
			Status:         smskit.StatusAccepted,
			Provider:       providerName,
			ProviderStatus: r.Status,
			Raw:            json.RawMessage(body),
		}, nil
	case "failed":
		return nil, toError(status, body)
	default:
		return nil, &smskit.ProviderError{
			Provider: providerName, Kind: smskit.ErrorProvider,
			Message: "unrecognized response status", HTTPStatus: status,
			Raw: json.RawMessage(body),
		}
	}
}
