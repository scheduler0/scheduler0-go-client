package scheduler0_go_client

import "fmt"

// ListExecutions retrieves job executions with query parameters
func (c *Client) ListExecutions(params ListExecutionsParams) (*PaginatedExecutionsResponse, error) {
	queryParams := map[string]string{
		"limit":  fmt.Sprintf("%d", params.Limit),
		"offset": fmt.Sprintf("%d", params.Offset),
	}

	// Only add date parameters if provided (optional)
	if params.StartDate != "" {
		queryParams["startDate"] = params.StartDate
	}
	if params.EndDate != "" {
		queryParams["endDate"] = params.EndDate
	}

	if params.ProjectID > 0 {
		queryParams["projectId"] = fmt.Sprintf("%d", params.ProjectID)
	}
	if params.JobID > 0 {
		queryParams["jobId"] = fmt.Sprintf("%d", params.JobID)
	}
	if params.State != "" {
		queryParams["state"] = params.State
	}
	if params.OrderBy != "" {
		queryParams["orderBy"] = params.OrderBy
	}
	if params.OrderDirection != "" {
		queryParams["orderDirection"] = params.OrderDirection
	}

	var accountIDOverride string
	if params.AccountID > 0 {
		accountIDOverride = fmt.Sprintf("%d", params.AccountID)
	}

	req, err := c.newRequestWithQuery("GET", "/executions", nil, queryParams, accountIDOverride)
	if err != nil {
		return nil, err
	}

	var result PaginatedExecutionsResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetExecutionMinuteBuckets retrieves execution counts grouped by minute
func (c *Client) GetExecutionMinuteBuckets(params GetExecutionMinuteBucketsParams) ([]ExecutionMinuteBucket, error) {
	queryParams := map[string]string{
		"startDate": params.StartDate,
		"endDate":   params.EndDate,
	}

	if params.JobID > 0 {
		queryParams["jobId"] = fmt.Sprintf("%d", params.JobID)
	}

	var accountIDOverride string
	if params.AccountID > 0 {
		accountIDOverride = fmt.Sprintf("%d", params.AccountID)
	}

	req, err := c.newRequestWithQuery("GET", "/executions/minute-buckets", nil, queryParams, accountIDOverride)
	if err != nil {
		return nil, err
	}

	var result ExecutionMinuteBucketsResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}
