package scheduler0_go_client

// Account represents an account
type Account struct {
	ID           int64            `json:"id"`
	Name         string           `json:"name"`
	Features     []AccountFeature `json:"features"`
	DateCreated  string           `json:"dateCreated"`
	DateModified *string          `json:"dateModified"`
}

// AccountFeature represents a feature associated with an account
type AccountFeature struct {
	AccountID int64  `json:"accountId"`
	FeatureID int64  `json:"featureId"`
	Feature   string `json:"feature"`
}

// AccountCreateRequestBody represents the request body for creating an account
type AccountCreateRequestBody struct {
	Name string `json:"name"`
}

// AccountUpdateRequestBody represents the request body for updating an account
type AccountUpdateRequestBody struct {
	Name string `json:"name"`
}

// AccountResponse represents the response for a single account
type AccountResponse struct {
	Success bool    `json:"success"`
	Data    Account `json:"data"`
}

// AccountJobExecutionsCount represents the execution count for an account
type AccountJobExecutionsCount struct {
	ID             int64  `json:"id"`
	AccountID      int64  `json:"accountId"`
	ExecutionCount int64  `json:"executionCount"`
	Tokens         int64  `json:"tokens"`
	DateCreated    string `json:"dateCreated"`
	DateModified   string `json:"dateModified"`
	NextResetDate  string `json:"nextResetDate"`
}

// AccountTokensResponse represents the response for getting account token balance
type AccountTokensResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Tokens int64 `json:"tokens"`
	} `json:"data"`
}

// AccountTokensAddRequest represents the request body for adding tokens
type AccountTokensAddRequest struct {
	Amount int64 `json:"amount"`
}

// AccountTokensAddResponse represents the response for adding tokens to an account
type AccountTokensAddResponse struct {
	Success bool `json:"success"`
	Data    struct {
		NewBalance int64 `json:"newBalance"`
	} `json:"data"`
}

// AccountExecutionCountResponse represents the response for account execution count
type AccountExecutionCountResponse struct {
	Success bool                      `json:"success"`
	Data    AccountJobExecutionsCount `json:"data"`
}

// AccountExecutionCountIncreaseResponse represents the response for increasing account execution count
type AccountExecutionCountIncreaseResponse struct {
	Success bool `json:"success"`
	Data    struct {
		NewExecutionCount uint64 `json:"newExecutionCount"`
	} `json:"data"`
}

// AIUsageDimension reports one AI quota dimension (prompt or classify) for an account: the
// feature-derived monthly Limit, the number of successful requests Used in the current period,
// and the Remaining allowance. Usage is log-derived on the backend (counted from the request
// logs since the period start), so these are the authoritative values to render.
type AIUsageDimension struct {
	Limit     int64 `json:"limit"`
	Used      int64 `json:"used"`
	Remaining int64 `json:"remaining"`
}

// AIUsage is the log-derived view of an account's AI request usage for the current period.
type AIUsage struct {
	AccountID     int64            `json:"accountId"`
	PeriodStart   string           `json:"periodStart"`
	NextResetDate string           `json:"nextResetDate"`
	Prompt        AIUsageDimension `json:"prompt"`
	Classify      AIUsageDimension `json:"classify"`
}

// AIUsageResponse represents the response for account AI usage.
type AIUsageResponse struct {
	Success bool    `json:"success"`
	Data    AIUsage `json:"data"`
}
