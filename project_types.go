package scheduler0_go_client

// Project represents a project
type Project struct {
	ID           int64   `json:"id"`
	AccountID    int64   `json:"accountId"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	DateCreated  string  `json:"dateCreated"`
	DateModified *string `json:"dateModified"`
	CreatedBy    string  `json:"createdBy"`
	ModifiedBy   *string `json:"modifiedBy"`
	DeletedBy    *string `json:"deletedBy"`
}

// ProjectResponse represents the response for a single project
type ProjectResponse struct {
	Success bool    `json:"success"`
	Data    Project `json:"data"`
}

// ProjectRequestBody represents the request body for creating a project
type ProjectRequestBody struct {
	AccountID   int64  `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   string `json:"createdBy"`
}

// ProjectUpdateRequestBody represents the request body for updating a project
type ProjectUpdateRequestBody struct {
	AccountID   int64  `json:"-"`
	Description string `json:"description"`
	ModifiedBy  string `json:"modifiedBy"`
}

// ProjectDeleteRequestBody represents the request body for deleting a project
type ProjectDeleteRequestBody struct {
	AccountID int64  `json:"-"`
	DeletedBy string `json:"deletedBy"`
}

// PaginatedProjectsResponse represents a paginated list of projects
type PaginatedProjectsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total    int       `json:"total"`
		Offset   int       `json:"offset"`
		Limit    int       `json:"limit"`
		Projects []Project `json:"projects"`
	} `json:"data"`
}

// ListProjectsParams represents parameters for listing projects
type ListProjectsParams struct {
	AccountID        int64  // Account ID override (0 to use client default)
	Limit            int    // Maximum number of items to return
	Offset           int    // Number of items to skip
	OrderBy          string // Field to order by (e.g., "date_created", "date_modified")
	OrderByDirection string // Direction to order ("asc" or "desc")
}
