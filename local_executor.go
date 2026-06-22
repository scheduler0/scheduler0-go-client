package scheduler0_go_client

import "fmt"

// RegisterLocalExecutor creates a new local-type executor on the server and returns its ID.
//
// POST /api/v1/local-executors
func (c *Client) RegisterLocalExecutor(body *LocalExecutorRegisterRequest) (*LocalExecutorRegisterResponse, error) {
	req, err := c.newRequest("POST", "/local-executors", body)
	if err != nil {
		return nil, err
	}

	var result LocalExecutorRegisterResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PullLocalExecutorJobs returns all active jobs assigned to the given local executor.
//
// GET /api/v1/local-executors/{id}/jobs
func (c *Client) PullLocalExecutorJobs(executorID int64) (*LocalExecutorJobsResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/local-executors/%d/jobs", executorID), nil)
	if err != nil {
		return nil, err
	}

	var result LocalExecutorJobsResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ReportLocalExecutions sends a batch of execution events to the server to be committed.
//
// POST /api/v1/local-executors/{id}/executions
func (c *Client) ReportLocalExecutions(executorID int64, reports []LocalExecutionReport) (*ReportLocalExecutionsResponse, error) {
	req, err := c.newRequest("POST", fmt.Sprintf("/local-executors/%d/executions", executorID), reports)
	if err != nil {
		return nil, err
	}

	var result ReportLocalExecutionsResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
