package scheduler0_go_client

import (
	"encoding/json"
	"fmt"
	"io"
)

// analyzeSuggestionsEnvelope mirrors POST /suggestions/analyze 200 response:
//
//	{"success":true,"data":{...}}
type analyzeSuggestionsEnvelope struct {
	Success bool                     `json:"success"`
	Data    AnalyzeSuggestionsResult `json:"data"`
}

// AnalyzeSuggestions analyzes a conversation and returns structured, explainable
// suggestions (commitments, requests, deadlines, follow-ups, etc.) plus the
// obligations they map to. An optional accountIDOverride sets the X-Account-ID
// header.
//
// English only: a non-en* locale in Options is rejected by the API with
// UNSUPPORTED_LOCALE (400).
func (c *Client) AnalyzeSuggestions(body *AnalyzeSuggestionsRequest, accountIDOverride ...string) (*AnalyzeSuggestionsResult, error) {
	req, err := c.newRequest("POST", "/ai/suggestions/analyze", body, accountIDOverride...)
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

	var envelope analyzeSuggestionsEnvelope
	if err := json.Unmarshal(rawBody, &envelope); err != nil {
		return nil, fmt.Errorf("parse analyze suggestions response: %w", err)
	}
	result := envelope.Data
	return &result, nil
}
