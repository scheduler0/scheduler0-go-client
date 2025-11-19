package scheduler0_go_client

import "fmt"

// ListJobs retrieves all jobs with optional query parameters
// projectID can be empty string to list all jobs for the account
// accountIDOverride is optional - if provided, overrides the client's default account ID
func (c *Client) ListJobs(projectID string, limit, offset int, orderBy, orderByDirection string, accountIDOverride ...string) (*PaginatedJobsResponse, error) {
	queryParams := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
	}

	// Only add projectId if it's not empty
	if projectID != "" {
		queryParams["projectId"] = projectID
	}

	if orderBy != "" {
		queryParams["orderBy"] = orderBy
	}
	if orderByDirection != "" {
		queryParams["orderByDirection"] = orderByDirection
	}

	var accountID string
	if len(accountIDOverride) > 0 {
		accountID = accountIDOverride[0]
	}
	req, err := c.newRequestWithQuery("GET", "/jobs", nil, queryParams, accountID)
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

