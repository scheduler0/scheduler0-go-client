package scheduler0_go_client

// Credential represents a credential returned by the Scheduler0 API.
//
// API key / secret model:
//   - APIKey is a stable plaintext identifier stored as-is in the DB. It is
//     always present in API responses.
//   - APISecret (the AES-GCM encrypted form) is never included in API responses.
//     Clients must send the original plaintext secret they received at creation.
//   - PlaintextSecret is returned exactly once, in the 201 create response.
//     It is never stored server-side and cannot be retrieved afterwards.
type Credential struct {
	ID              int64    `json:"id"`
	AccountID       int64    `json:"accountId"`
	Archived        bool     `json:"archived"`
	APIKey          string   `json:"apiKey"`
	PlaintextSecret string   `json:"plaintextSecret,omitempty"` // Returned once on creation only
	DateCreated     string   `json:"dateCreated"`
	DateModified    *string  `json:"dateModified"`
	CreatedBy       string   `json:"createdBy"`
	ModifiedBy      *string  `json:"modifiedBy"`
	DeletedBy       *string  `json:"deletedBy"`
	ArchivedBy      *string  `json:"archivedBy"`
	ExpiresAt       *string  `json:"expiresAt,omitempty"`
	Scopes          []string `json:"scopes,omitempty"`
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
	AccountID        int64  // Optional: Account ID override (0 uses client default)
	Limit            int    // Maximum number of items to return
	Offset           int    // Number of items to skip
	OrderBy          string // Field to order by (e.g., "date_created", "date_modified")
	OrderByDirection string // Direction to order ("asc" or "desc")
}

// CredentialCreateRequestBody represents the request body for creating a credential.
// Scopes must be a subset of {"read","write","execute"} and is required by the API.
type CredentialCreateRequestBody struct {
	AccountID int64    `json:"-"`
	Archived  bool     `json:"archived,omitempty"`
	CreatedBy string   `json:"createdBy"`
	Scopes    []string `json:"scopes"`
	// ExpiresInSeconds optionally requests a shorter credential lifetime. The
	// server clamps it to its allowed range; nil uses the server default expiry.
	ExpiresInSeconds *int64 `json:"expiresInSeconds,omitempty"`
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
