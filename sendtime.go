package scheduler0_go_client

import (
	"encoding/json"
	"fmt"
	"io"
)

// sendTimeSuggestionsEnvelope mirrors POST /ai/suggestions/time 200 response:
//
//	{"success":true,"data":{...}}
type sendTimeSuggestionsEnvelope struct {
	Success bool                      `json:"success"`
	Data    SendTimeSuggestionsResult `json:"data"`
}

// SendTimeSuggestions recommends suitable future send times for a message given
// sender/recipient time zones, working hours, quiet hours, weekends, priority,
// and coverage rules. The engine is deterministic and does not send the message
// or create a job. An optional accountIDOverride sets the X-Account-ID header.
func (c *Client) SendTimeSuggestions(body *SendTimeSuggestionsRequest, accountIDOverride ...string) (*SendTimeSuggestionsResult, error) {
	req, err := c.newRequest("POST", "/ai/suggestions/time", body, accountIDOverride...)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s", string(rawBody))
	}

	var envelope sendTimeSuggestionsEnvelope
	if err := json.Unmarshal(rawBody, &envelope); err != nil {
		return nil, fmt.Errorf("parse send time suggestions response: %w", err)
	}
	result := envelope.Data
	return &result, nil
}
