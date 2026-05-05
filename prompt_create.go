package scheduler0_go_client

// promptJobsEnvelope mirrors the API response envelope:
//
//	{"success":true,"data":[{"provider":"openai","jobs":[...],...},...]}
type promptJobsEnvelope struct {
	Success bool                   `json:"success"`
	Data    []PromptProviderResult `json:"data"`
}

// CreateJobFromPrompt creates job configurations from an AI prompt.
// It returns the flattened list of jobs across all providers.
// An optional accountIDOverride can be supplied to set the X-Account-ID header.
func (c *Client) CreateJobFromPrompt(body *PromptJobRequest, accountIDOverride ...string) ([]PromptJobResponse, error) {
	req, err := c.newRequest("POST", "/prompt", body, accountIDOverride...)
	if err != nil {
		return nil, err
	}

	var result promptJobsEnvelope
	if err = c.do(req, &result); err != nil {
		return nil, err
	}

	// Flatten jobs from all provider results into a single slice.
	var jobs []PromptJobResponse
	for _, pr := range result.Data {
		jobs = append(jobs, pr.Jobs...)
	}
	return jobs, nil
}

