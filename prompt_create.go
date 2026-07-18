package scheduler0_go_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// promptJobsEnvelope mirrors the API response envelope:
//
//	{"success":true,"data":{"providers":[...],"classification":{...}}}
type promptJobsEnvelope struct {
	Success bool         `json:"success"`
	Data    PromptResult `json:"data"`
}

// promptSkippedEnvelope mirrors the 422 response body for a guardrail rejection:
//
//	{"success":false,"data":{"message":"...","classification":{...}}}
type promptSkippedEnvelope struct {
	Success bool `json:"success"`
	Data    struct {
		Message        string                `json:"message"`
		Classification *IntentClassification `json:"classification,omitempty"`
	} `json:"data"`
}

// classifyEnvelope mirrors POST /prompt/classify 200 response:
//
//	{"success":true,"data":{"classification":{...}}}
type classifyEnvelope struct {
	Success bool `json:"success"`
	Data    struct {
		Classification IntentClassification `json:"classification"`
	} `json:"data"`
}

// doPrompt executes an HTTP request against the prompt endpoint and handles the
// structured 422 PromptSkippedError response in addition to generic errors.
func (c *Client) doPrompt(req *http.Request) (*PromptResult, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode == http.StatusUnprocessableEntity {
		var skipped promptSkippedEnvelope
		if jsonErr := json.Unmarshal(body, &skipped); jsonErr == nil && skipped.Data.Message != "" {
			return nil, &PromptSkippedError{
				Message:        skipped.Data.Message,
				Classification: skipped.Data.Classification,
			}
		}
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var envelope promptJobsEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("parse prompt response: %w", err)
	}
	result := envelope.Data
	return &result, nil
}

// CreateJobFromPrompt creates job configurations from an AI prompt.
// It returns a PromptResult containing all provider results and the full
// intent classification. An optional accountIDOverride can be supplied to
// set the X-Account-ID header.
func (c *Client) CreateJobFromPrompt(body *PromptJobRequest, accountIDOverride ...string) (*PromptResult, error) {
	req, err := c.newRequest("POST", "/ai/prompt", body, accountIDOverride...)
	if err != nil {
		return nil, err
	}
	return c.doPrompt(req)
}

// ClassifyPrompt runs only the intent guardrail against a prompt. No AI model
// is invoked and no credits are consumed. Returns an error when the classifier
// is not configured on the server (503).
func (c *Client) ClassifyPrompt(body *ClassifyPromptRequest, accountIDOverride ...string) (*IntentClassification, error) {
	req, err := c.newRequest("POST", "/ai/prompt/classify", body, accountIDOverride...)
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

	var envelope classifyEnvelope
	if err := json.Unmarshal(rawBody, &envelope); err != nil {
		return nil, fmt.Errorf("parse classify response: %w", err)
	}
	result := envelope.Data.Classification
	return &result, nil
}

// IsPromptSkippedError reports whether err is a *PromptSkippedError.
func IsPromptSkippedError(err error) bool {
	var t *PromptSkippedError
	return errors.As(err, &t)
}
