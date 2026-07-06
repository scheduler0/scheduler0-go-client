package scheduler0_go_client

import "fmt"

// PromptRequest mirrors the AccountPromptRequest model from scheduler0-private: a single
// persisted AI prompt execution.
type PromptRequest struct {
	ID               uint64  `json:"id"`
	AccountID        uint64  `json:"account_id"`
	Prompt           string  `json:"prompt"`
	Provider         string  `json:"provider"`
	Model            string  `json:"model"`
	Output           string  `json:"output"`
	InputTokens      uint64  `json:"input_tokens"`
	OutputTokens     uint64  `json:"output_tokens"`
	TotalTokens      uint64  `json:"total_tokens"`
	DurationMs       uint64  `json:"duration_ms"`
	EstimatedCostUSD float64 `json:"estimated_cost_usd"`
	Status           string  `json:"status"`
	Error            string  `json:"error,omitempty"`
	DateCreated      string  `json:"date_created"`
}

// ListPromptRequestsParams narrows and paginates the prompt-request log query.
type ListPromptRequestsParams struct {
	AccountID uint64
	Provider  string
	Model     string
	Status    string
	Search    string
	StartDate string // RFC3339
	EndDate   string // RFC3339
	Order     string // ASC or DESC
	Limit     uint64
	Offset    uint64
}

// PromptRequestsData is the payload returned by the prompt-requests endpoint.
type PromptRequestsData struct {
	Requests []PromptRequest `json:"requests"`
	Total    uint64          `json:"total"`
	Limit    uint64          `json:"limit"`
	Offset   uint64          `json:"offset"`
}

// PromptRequestsResponse wraps the standard API envelope for the prompt-request log.
type PromptRequestsResponse struct {
	Success bool               `json:"success"`
	Data    PromptRequestsData `json:"data"`
}

// ListPromptRequests retrieves the account's AI prompt-request log with filters/pagination.
func (c *Client) ListPromptRequests(params ListPromptRequestsParams) (*PromptRequestsResponse, error) {
	queryParams := map[string]string{
		"limit":  fmt.Sprintf("%d", params.Limit),
		"offset": fmt.Sprintf("%d", params.Offset),
	}
	if params.Provider != "" {
		queryParams["provider"] = params.Provider
	}
	if params.Model != "" {
		queryParams["model"] = params.Model
	}
	if params.Status != "" {
		queryParams["status"] = params.Status
	}
	if params.Search != "" {
		queryParams["search"] = params.Search
	}
	if params.StartDate != "" {
		queryParams["start"] = params.StartDate
	}
	if params.EndDate != "" {
		queryParams["end"] = params.EndDate
	}
	if params.Order != "" {
		queryParams["order"] = params.Order
	}

	var accountIDOverride string
	if params.AccountID > 0 {
		accountIDOverride = fmt.Sprintf("%d", params.AccountID)
	}

	req, err := c.newRequestWithQuery("GET", "/prompt-requests", nil, queryParams, accountIDOverride)
	if err != nil {
		return nil, err
	}

	var result PromptRequestsResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
