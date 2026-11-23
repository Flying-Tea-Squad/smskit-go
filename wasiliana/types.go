package wasiliana

import "encoding/json"

// sendRequest is the single-send payload. The optional linkid, message_uid
// and is_otp fields are not exposed yet.
type sendRequest struct {
	Recipients []string `json:"recipients"`
	From       string   `json:"from"`
	Message    string   `json:"message"`
}

// apiResponse is the documented envelope: {"status":"success|failed","data":"..."}.
// data is kept raw because the docs only show strings.
type apiResponse struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
}
