package scheduler0_go_client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// scheduleEnvelope mirrors POST /ai/schedule 201 response:
//
//	{"success":true,"data":{...}}
type scheduleEnvelope struct {
	Success bool           `json:"success"`
	Data    ScheduleResult `json:"data"`
}

// ScheduleFromPrompt turns a natural-language prompt into scheduled jobs: the server runs
// the prompt pipeline (intent guardrail + generation), resolves or creates a project, picks
// the executor whose description/tags best match the prompt (or uses a pinned/only executor),
// and creates the jobs synchronously. A rejected prompt returns a *PromptSkippedError (HTTP
// 422); use IsPromptSkippedError to detect it.
func (c *Client) ScheduleFromPrompt(body *SchedulePromptRequest) (*ScheduleResult, error) {
	req, err := c.newRequest("POST", "/ai/schedule", body)
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

	if resp.StatusCode == http.StatusUnprocessableEntity {
		var skipped promptSkippedEnvelope
		if jsonErr := json.Unmarshal(rawBody, &skipped); jsonErr == nil && skipped.Data.Message != "" {
			return nil, &PromptSkippedError{
				Message:        skipped.Data.Message,
				Classification: skipped.Data.Classification,
			}
		}
		return nil, fmt.Errorf("API error: %s", string(rawBody))
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s", string(rawBody))
	}

	var envelope scheduleEnvelope
	if err := json.Unmarshal(rawBody, &envelope); err != nil {
		return nil, fmt.Errorf("parse schedule response: %w", err)
	}
	result := envelope.Data
	return &result, nil
}
