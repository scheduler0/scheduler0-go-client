package scheduler0_go_client

import "fmt"

// ListCredentials retrieves all credentials with optional query parameters
func (c *Client) ListCredentials(limit, offset int, orderBy, orderByDirection string) (*PaginatedCredentialsResponse, error) {
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

	req, err := c.newRequestWithQuery("GET", "/credentials", nil, queryParams)
	if err != nil {
		return nil, err
	}

	var result PaginatedCredentialsResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

