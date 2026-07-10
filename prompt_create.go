package scheduler0_go_client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// promptJobsEnvelope mirrors the API response envelope:
//
//	{"success":true,"data":{"providers":[...],"classification":{...}}}
type promptJobsEnvelope struct {
	Success bool `json:"success"`
	Data    struct {
		Providers      []PromptProviderResult `json:"providers"`
		Classification *PromptClassification  `json:"classification,omitempty"`
	} `json:"data"`
}

// promptSkippedEnvelope mirrors the 422 guardrail-rejection body:
//
//	{"success":false,"message":"...","classification":{...}}
type promptSkippedEnvelope struct {
	Success        bool                  `json:"success"`
	Message        string                `json:"message"`
	Classification *PromptClassification `json:"classification,omitempty"`
}

// CreateJobFromPrompt creates job configurations from an AI prompt.
// An optional accountIDOverride can be supplied to set the X-Account-ID header.
// If the intent guardrail rejects the prompt the error is of type *PromptSkippedError.
func (c *Client) CreateJobFromPrompt(body *PromptJobRequest, accountIDOverride ...string) (*PromptResponse, error) {
	req, err := c.newRequest("POST", "/prompt", body, accountIDOverride...)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnprocessableEntity {
		var skipped promptSkippedEnvelope
		if decErr := json.NewDecoder(resp.Body).Decode(&skipped); decErr != nil {
			return nil, &PromptSkippedError{Message: "prompt rejected by guardrail"}
		}
		return nil, &PromptSkippedError{
			Message:        skipped.Message,
			Classification: skipped.Classification,
		}
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("prompt API returned status %d", resp.StatusCode)
	}

	var result promptJobsEnvelope
	if decErr := json.NewDecoder(resp.Body).Decode(&result); decErr != nil {
		return nil, decErr
	}
	if !result.Success {
		return nil, fmt.Errorf("prompt API returned success=false")
	}

	return &PromptResponse{
		Providers:      result.Data.Providers,
		Classification: result.Data.Classification,
	}, nil
}
