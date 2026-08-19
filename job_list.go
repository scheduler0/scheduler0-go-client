package scheduler0_go_client

import "fmt"

// ListJobs retrieves all jobs with optional query parameters
func (c *Client) ListJobs(params ListJobsParams) (*PaginatedJobsResponse, error) {
	queryParams := map[string]string{
		"offset": fmt.Sprintf("%d", params.Offset),
	}
	// Omit a zero limit so the server applies its authoritative default page size (and its
	// maximum cap) rather than being asked for an explicit zero-length page.
	if params.Limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", params.Limit)
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

