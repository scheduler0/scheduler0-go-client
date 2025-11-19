package scheduler0_go_client

// CreateJob creates a new job
// Note: This is a convenience method that wraps a single job in an array.
// The API always expects an array and returns 202 Accepted with a request ID for async tracking.
// For better control, use BatchCreateJobs directly.
// accountIDOverride is optional - if provided, overrides the client's default account ID
func (c *Client) CreateJob(body *JobRequestBody, accountIDOverride ...string) (*BatchJobResponse, error) {
	return c.BatchCreateJobs([]JobRequestBody{*body}, accountIDOverride...)
}

// BatchCreateJobs creates multiple jobs in a single request
// accountIDOverride is optional - if provided, overrides the client's default account ID
func (c *Client) BatchCreateJobs(jobs []JobRequestBody, accountIDOverride ...string) (*BatchJobResponse, error) {
	var accountID string
	if len(accountIDOverride) > 0 {
		accountID = accountIDOverride[0]
	}
	req, err := c.newRequest("POST", "/jobs", jobs, accountID)
	if err != nil {
		return nil, err
	}

	var result BatchJobResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

