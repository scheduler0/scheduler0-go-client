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

// GetDateRangeAnalytics retrieves execution counts grouped by minute buckets for a date range
// All dates and times should be in UTC (timezone conversion should be done on frontend)
func (c *Client) GetDateRangeAnalytics(params GetDateRangeAnalyticsParams) (*DateRangeAnalyticsAPIResponse, error) {
	queryParams := map[string]string{
		"startDate": params.StartDate,
		"startTime": params.StartTime,
	}

	var accountIDOverride string
	if params.AccountID > 0 {
		accountIDOverride = fmt.Sprintf("%d", params.AccountID)
	}

	req, err := c.newRequestWithQuery("GET", "/executions/analytics", nil, queryParams, accountIDOverride)
	if err != nil {
		return nil, err
	}

	var result DateRangeAnalyticsAPIResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetExecutionTotals retrieves total counts of scheduled, success, and failed executions for an account
func (c *Client) GetExecutionTotals(accountID int64) (*ExecutionTotalsAPIResponse, error) {
	var accountIDOverride string
	if accountID > 0 {
		accountIDOverride = fmt.Sprintf("%d", accountID)
	}

	req, err := c.newRequestWithQuery("GET", "/executions/totals", nil, nil, accountIDOverride)
	if err != nil {
		return nil, err
	}

	var result ExecutionTotalsAPIResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CleanupOldExecutionLogs cleans up old execution logs for an account based on retention period
// accountIDOverride is optional - if provided, overrides the client's default account ID
func (c *Client) CleanupOldExecutionLogs(accountID string, retentionMonths int, accountIDOverride ...string) (*CleanupOldLogsResponse, error) {
	requestBody := CleanupOldLogsRequestBody{
		AccountID:       accountID,
		RetentionMonths: retentionMonths,
	}

	var accountIDHeader string
	if len(accountIDOverride) > 0 && accountIDOverride[0] != "" {
		accountIDHeader = accountIDOverride[0]
	} else {
		accountIDHeader = accountID
	}

	req, err := c.newRequest("POST", "/executions/cleanup-old-logs", requestBody, accountIDHeader)
	if err != nil {
		return nil, err
	}

	var result CleanupOldLogsResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
