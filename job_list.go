package scheduler0_go_client

import "fmt"

// ListJobs retrieves all jobs with optional query parameters
func (c *Client) ListJobs(params ListJobsParams) (*PaginatedJobsResponse, error) {
	queryParams := map[string]string{
		"limit":  fmt.Sprintf("%d", params.Limit),
		"offset": fmt.Sprintf("%d", params.Offset),
	}

	// Only add projectId if it's not empty
	if params.ProjectID != "" {
		queryParams["projectId"] = params.ProjectID
	}

	if params.OrderBy != "" {
		queryParams["orderBy"] = params.OrderBy
	}
	if params.OrderByDirection != "" {
		queryParams["orderByDirection"] = params.OrderByDirection
	}

	var accountIDOverride string
	if params.AccountID > 0 {
		accountIDOverride = fmt.Sprintf("%d", params.AccountID)
	}
	req, err := c.newRequestWithQuery("GET", "/jobs", nil, queryParams, accountIDOverride)
	if err != nil {
		return nil, err
	}

	var result PaginatedJobsResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

