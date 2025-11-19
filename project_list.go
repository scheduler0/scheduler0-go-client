package scheduler0_go_client

import "fmt"

// ListProjects retrieves all projects with optional query parameters
func (c *Client) ListProjects(accountID int64, limit, offset int, orderBy, orderByDirection string) (*PaginatedProjectsResponse, error) {
	queryParams := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
	}

	if orderBy != "" {
		queryParams["orderBy"] = orderBy
	}
	if orderByDirection != "" {
		queryParams["orderByDirection"] = orderByDirection
	}

	var accountIDOverride string
	if accountID > 0 {
		accountIDOverride = fmt.Sprintf("%d", accountID)
	}

	req, err := c.newRequestWithQuery("GET", "/projects", nil, queryParams, accountIDOverride)
	if err != nil {
		return nil, err
	}

	var result PaginatedProjectsResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
