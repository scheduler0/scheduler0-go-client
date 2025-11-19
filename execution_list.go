package scheduler0_go_client

import "fmt"

// ListExecutions retrieves job executions with query parameters
func (c *Client) ListExecutions(params ListExecutionsParams) (*PaginatedExecutionsResponse, error) {
	queryParams := map[string]string{
		"startDate": params.StartDate,
		"endDate":   params.EndDate,
		"limit":     fmt.Sprintf("%d", params.Limit),
		"offset":    fmt.Sprintf("%d", params.Offset),
	}

	if params.ProjectID > 0 {
		queryParams["projectId"] = fmt.Sprintf("%d", params.ProjectID)
	}
	if params.JobID > 0 {
		queryParams["jobId"] = fmt.Sprintf("%d", params.JobID)
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
