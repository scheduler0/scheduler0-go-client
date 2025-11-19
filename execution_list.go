package scheduler0_go_client

import "fmt"

// ListExecutions retrieves job executions with query parameters
func (c *Client) ListExecutions(startDate, endDate string, projectID, jobID, accountID int64, limit, offset int) (*PaginatedExecutionsResponse, error) {
	queryParams := map[string]string{
		"startDate": startDate,
		"endDate":   endDate,
		"limit":     fmt.Sprintf("%d", limit),
		"offset":    fmt.Sprintf("%d", offset),
	}

	if projectID > 0 {
		queryParams["projectId"] = fmt.Sprintf("%d", projectID)
	}
	if jobID > 0 {
		queryParams["jobId"] = fmt.Sprintf("%d", jobID)
	}

	var accountIDOverride string
	if accountID > 0 {
		accountIDOverride = fmt.Sprintf("%d", accountID)
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
