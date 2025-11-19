package scheduler0_go_client

// Feature represents a feature
type Feature struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	DateCreated  string  `json:"dateCreated"`
	DateModified *string `json:"dateModified"`
}

// FeatureRequest represents a request to add/remove a feature
type FeatureRequest struct {
	AccountID int64 `json:"-"`
	FeatureID int64 `json:"featureId"`
}

// FeatureRequestResponse represents the response for feature operations
type FeatureRequestResponse struct {
	Success bool           `json:"success"`
	Data    FeatureRequest `json:"data"`
}

// FeaturesResponse represents the response for listing features
type FeaturesResponse struct {
	Success bool      `json:"success"`
	Data    []Feature `json:"data"`
}
