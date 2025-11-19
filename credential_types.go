package scheduler0_go_client

// Credential represents a credential
type Credential struct {
	ID           int64   `json:"id"`
	AccountID    int64   `json:"accountId"`
	Archived     bool    `json:"archived"`
	APIKey       string  `json:"apiKey"`
	APISecret    string  `json:"apiSecret"`
	DateCreated  string  `json:"dateCreated"`
	DateModified *string `json:"dateModified"`
	DateDeleted  *string `json:"dateDeleted"`
	CreatedBy    string  `json:"createdBy"`
	ModifiedBy   *string `json:"modifiedBy"`
	DeletedBy    *string `json:"deletedBy"`
	ArchivedBy   *string `json:"archivedBy"`
}

// CredentialResponse represents the response for a single credential
type CredentialResponse struct {
	Success bool       `json:"success"`
	Data    Credential `json:"data"`
}

// PaginatedCredentialsResponse represents a paginated list of credentials
type PaginatedCredentialsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total       int          `json:"total"`
		Offset      int          `json:"offset"`
		Limit       int          `json:"limit"`
		Credentials []Credential `json:"credentials"`
	} `json:"data"`
}

// ListCredentialsParams represents parameters for listing credentials
type ListCredentialsParams struct {
	Limit            int    // Maximum number of items to return
	Offset           int    // Number of items to skip
	OrderBy          string // Field to order by (e.g., "date_created", "date_modified")
	OrderByDirection string // Direction to order ("asc" or "desc")
}

// CredentialCreateRequestBody represents the request body for creating a credential
type CredentialCreateRequestBody struct {
	AccountID int64  `json:"-"`
	Archived  bool   `json:"archived,omitempty"`
	CreatedBy string `json:"createdBy"`
}

// CredentialUpdateRequestBody represents the request body for updating a credential
type CredentialUpdateRequestBody struct {
	AccountID  int64  `json:"-"`
	Archived   bool   `json:"archived,omitempty"`
	ModifiedBy string `json:"modifiedBy"`
}

// CredentialDeleteRequestBody represents the request body for deleting a credential
type CredentialDeleteRequestBody struct {
	AccountID int64  `json:"-"`
	DeletedBy string `json:"deletedBy"`
}

// CredentialArchiveRequestBody represents the request body for archiving a credential
type CredentialArchiveRequestBody struct {
	AccountID  int64  `json:"-"`
	ArchivedBy string `json:"archivedBy"`
}
