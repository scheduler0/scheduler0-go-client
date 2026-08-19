package scheduler0_go_client

import "fmt"

// ListProjects retrieves all projects with optional query parameters
func (c *Client) ListProjects(params ListProjectsParams) (*PaginatedProjectsResponse, error) {
	queryParams := map[string]string{
		"offset": fmt.Sprintf("%d", params.Offset),
	}
	// Omit a zero limit so the server applies its authoritative default page size (and its
	// maximum cap) rather than being asked for an explicit zero-length page.
	if params.Limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", params.Limit)
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
