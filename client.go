package scheduler0_go_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	APIKey     string
	APISecret  string
	Version    string
	// Basic Auth for peer communication
	Username string
	Password string
	// Account ID for most endpoints
	AccountID string
}

func NewClient(baseURL, version string, options ...ClientOption) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	client := &Client{
		BaseURL:    u,
		HTTPClient: &http.Client{},
		Version:    version,
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	return client, nil
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithAPIKey sets the API key and secret for authentication
func WithAPIKey(apiKey, apiSecret string) ClientOption {
	return func(c *Client) {
		c.APIKey = apiKey
		c.APISecret = apiSecret
	}
}

// WithAccountID sets the account ID for requests
func WithAccountID(accountID string) ClientOption {
	return func(c *Client) {
		c.AccountID = accountID
	}
}

// WithBasicAuth sets the username and password for basic authentication
func WithBasicAuth(username, password string) ClientOption {
	return func(c *Client) {
		c.Username = username
		c.Password = password
	}
}

// Convenience functions for common use cases

// NewAPIClient creates a client with API key authentication
func NewAPIClient(baseURL, version, apiKey, apiSecret string) (*Client, error) {
	return NewClient(baseURL, version, WithAPIKey(apiKey, apiSecret))
}

// NewAPIClientWithAccount creates a client with API key authentication and account ID
func NewAPIClientWithAccount(baseURL, version, apiKey, apiSecret, accountID string) (*Client, error) {
	return NewClient(baseURL, version, WithAPIKey(apiKey, apiSecret), WithAccountID(accountID))
}

// NewBasicAuthClient creates a client with basic authentication for peer communication
func NewBasicAuthClient(baseURL, version, username, password string) (*Client, error) {
	return NewClient(baseURL, version, WithBasicAuth(username, password))
}

func (c *Client) newRequest(method, endpoint string, body interface{}) (*http.Request, error) {
	versionPrefix := fmt.Sprintf("/api/%s/", c.Version)

	rel := &url.URL{Path: path.Join(fmt.Sprintf("%s%s", c.BaseURL.Path, versionPrefix), endpoint)}
	u := c.BaseURL.ResolveReference(rel)

	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Set authentication based on client type
	if c.Username != "" && c.Password != "" {
		// Basic Auth for peer communication
		req.SetBasicAuth(c.Username, c.Password)
		req.Header.Set("X-Peer", "cmd")
	} else if c.APIKey != "" && c.APISecret != "" {
		// API Key + Secret authentication
		req.Header.Set("X-API-Key", c.APIKey)
		req.Header.Set("X-API-Secret", c.APISecret)
	}

	// Add account ID if available
	if c.AccountID != "" {
		req.Header.Set("X-Account-ID", c.AccountID)
	}

	return req, nil
}

func (c *Client) newRequestWithQuery(method, endpoint string, body interface{}, queryParams map[string]string) (*http.Request, error) {
	versionPrefix := fmt.Sprintf("/api/%s/", c.Version)

	rel := &url.URL{Path: path.Join(fmt.Sprintf("%s%s", c.BaseURL.Path, versionPrefix), endpoint)}
	u := c.BaseURL.ResolveReference(rel)

	// Add query parameters
	if len(queryParams) > 0 {
		q := u.Query()
		for key, value := range queryParams {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
	}

	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Set authentication based on client type
	if c.Username != "" && c.Password != "" {
		// Basic Auth for peer communication
		req.SetBasicAuth(c.Username, c.Password)
		req.Header.Set("X-Peer", "cmd")
	} else if c.APIKey != "" && c.APISecret != "" {
		// API Key + Secret authentication
		req.Header.Set("X-API-Key", c.APIKey)
		req.Header.Set("X-API-Secret", c.APISecret)
	}

	// Add account ID if available
	if c.AccountID != "" {
		req.Header.Set("X-Account-ID", c.AccountID)
	}

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}

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

// Define the response structs

type PaginatedCredentialsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total       int          `json:"total"`
		Offset      int          `json:"offset"`
		Limit       int          `json:"limit"`
		Credentials []Credential `json:"credentials"`
	} `json:"data"`
}

type Credential struct {
	ID        int    `json:"id"`
	Archived  bool   `json:"archived"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
	CreatedAt string `json:"date_created"`
}

// CredentialResponse represents the response for a single credential
type CredentialResponse struct {
	Success bool       `json:"success"`
	Data    Credential `json:"data"`
}

// CreateCredential creates a new credential
func (c *Client) CreateCredential() (*CredentialResponse, error) {
	req, err := c.newRequest("POST", "/credentials", nil)
	if err != nil {
		return nil, err
	}

	var result CredentialResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCredential retrieves a single credential by ID
func (c *Client) GetCredential(id string) (*CredentialResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/credentials/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var result CredentialResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateCredential updates an existing credential
func (c *Client) UpdateCredential(id string) (*CredentialResponse, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("/credentials/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var result CredentialResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteCredential deletes a credential by ID
func (c *Client) DeleteCredential(id string) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/credentials/%s", id), nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// Execution represents a job execution log
type Execution struct {
	ID                int    `json:"id"`
	UniqueID          string `json:"uniqueId"`
	State             string `json:"state"`
	NodeID            string `json:"nodeId"`
	JobID             string `json:"jobId"`
	LastExecutionTime string `json:"lastExecutionDatetime"`
	NextExecutionTime string `json:"nextExecutionDatetime"`
	JobQueueVersion   int    `json:"jobQueueVersion"`
	ExecutionVersion  int    `json:"executionVersion"`
	Logs              string `json:"logs"`
	CreatedAt         string `json:"date_created"`
}

// ExecutionResponse represents the response for a single execution
type ExecutionResponse struct {
	Success bool      `json:"success"`
	Data    Execution `json:"data"`
}

// PaginatedExecutionsResponse represents a paginated list of executions
type PaginatedExecutionsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total      int         `json:"total"`
		Offset     int         `json:"offset"`
		Limit      int         `json:"limit"`
		Executions []Execution `json:"executions"`
	} `json:"data"`
}

// ListExecutions retrieves job executions with query parameters
func (c *Client) ListExecutions(startDate, endDate string, projectID, jobID int64, limit, offset int) (*PaginatedExecutionsResponse, error) {
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

	req, err := c.newRequestWithQuery("GET", "/executions", nil, queryParams)
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

// Executor represents a job executor
type Executor struct {
	ID               int64   `json:"id"`
	AccountID        int64   `json:"accountId"`
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	Region           string  `json:"region"`
	CloudProvider    string  `json:"cloudProvider"`
	CloudResourceURL string  `json:"cloudResourceUrl"`
	CloudAPIKey      string  `json:"cloudApiKey"`
	CloudAPISecret   string  `json:"cloudApiSecret"`
	WebhookURL       string  `json:"webhookUrl"`
	WebhookSecret    string  `json:"webhookSecret"`
	WebhookMethod    string  `json:"webhookMethod"`
	DateCreated      string  `json:"dateCreated"`
	DateModified     *string `json:"dateModified"`
	DateDeleted      *string `json:"dateDeleted"`
	CreatedBy        string  `json:"createdBy"`
	ModifiedBy       *string `json:"modifiedBy"`
	DeletedBy        *string `json:"deletedBy"`
}

// ExecutorResponse represents the response for a single executor
type ExecutorResponse struct {
	Success bool     `json:"success"`
	Data    Executor `json:"data"`
}

// ExecutorRequestBody represents the request body for creating/updating an executor
type ExecutorRequestBody struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	Region           string `json:"region"`
	CloudProvider    string `json:"cloudProvider"`
	CloudResourceURL string `json:"cloudResourceUrl"`
	CloudAPIKey      string `json:"cloudApiKey,omitempty"`
	CloudAPISecret   string `json:"cloudApiSecret,omitempty"`
	WebhookURL       string `json:"webhookUrl,omitempty"`
	WebhookSecret    string `json:"webhookSecret,omitempty"`
	WebhookMethod    string `json:"webhookMethod,omitempty"`
}

// PaginatedExecutorsResponse represents a paginated list of executors
type PaginatedExecutorsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total     int        `json:"total"`
		Offset    int        `json:"offset"`
		Limit     int        `json:"limit"`
		Executors []Executor `json:"executors"`
	} `json:"data"`
}

// ListExecutors retrieves all executors with optional query parameters
func (c *Client) ListExecutors(limit, offset int, orderBy, orderByDirection string) (*PaginatedExecutorsResponse, error) {
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

	req, err := c.newRequestWithQuery("GET", "/executors", nil, queryParams)
	if err != nil {
		return nil, err
	}

	var result PaginatedExecutorsResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateExecutor creates a new executor
func (c *Client) CreateExecutor(body *ExecutorRequestBody) (*ExecutorResponse, error) {
	req, err := c.newRequest("POST", "/executors", body)
	if err != nil {
		return nil, err
	}

	var result ExecutorResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetExecutor retrieves a single executor by ID
func (c *Client) GetExecutor(id string) (*ExecutorResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/executors/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var result ExecutorResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateExecutor updates an existing executor
func (c *Client) UpdateExecutor(id string, body *ExecutorRequestBody) (*ExecutorResponse, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("/executors/%s", id), body)
	if err != nil {
		return nil, err
	}

	var result ExecutorResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteExecutor deletes an executor by ID
func (c *Client) DeleteExecutor(id string) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/executors/%s", id), nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

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

// AccountResponse represents the response for a single account
type AccountResponse struct {
	Success bool    `json:"success"`
	Data    Account `json:"data"`
}

// Feature represents a feature
type Feature struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	DateCreated  string  `json:"dateCreated"`
	DateModified *string `json:"dateModified"`
}

// FeatureRequest represents a request to add/remove a feature
type FeatureRequest struct {
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

// AsyncTask represents an async task
type AsyncTask struct {
	ID          int64  `json:"id"`
	RequestID   string `json:"requestId"`
	Input       string `json:"input"`
	Output      string `json:"output"`
	Service     string `json:"service"`
	State       int    `json:"state"`
	DateCreated string `json:"dateCreated"`
}

// AsyncTaskResponse represents the response for a single async task
type AsyncTaskResponse struct {
	Success bool      `json:"success"`
	Data    AsyncTask `json:"data"`
}

// Project represents a project
type Project struct {
	ID          string `json:"id"`
	AccountID   string `json:"accountId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DateCreated string `json:"date_created"`
}

// ProjectResponse represents the response for a single project
type ProjectResponse struct {
	Success bool    `json:"success"`
	Data    Project `json:"data"`
}

// ProjectRequestBody represents the request body for creating a project
type ProjectRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProjectUpdateRequestBody represents the request body for updating a project
type ProjectUpdateRequestBody struct {
	Description string `json:"description"`
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

// ListProjects retrieves all projects with optional query parameters
func (c *Client) ListProjects(limit, offset int) (*PaginatedProjectsResponse, error) {
	queryParams := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
	}

	req, err := c.newRequestWithQuery("GET", "/projects", nil, queryParams)
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

// CreateProject creates a new project
func (c *Client) CreateProject(body *ProjectRequestBody) (*ProjectResponse, error) {
	req, err := c.newRequest("POST", "/projects", body)
	if err != nil {
		return nil, err
	}

	var result ProjectResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetProject retrieves a single project by ID
func (c *Client) GetProject(id string) (*ProjectResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/projects/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var result ProjectResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateProject updates an existing project
func (c *Client) UpdateProject(id string, body *ProjectUpdateRequestBody) (*ProjectResponse, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("/projects/%s", id), body)
	if err != nil {
		return nil, err
	}

	var result ProjectResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteProject deletes a project by ID
func (c *Client) DeleteProject(id string) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/projects/%s", id), nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// Job represents a scheduled job
type Job struct {
	ID                int64   `json:"id,omitempty"`
	AccountID         int64   `json:"accountId,omitempty"`
	ProjectID         int64   `json:"projectId,omitempty"`
	ExecutorID        *int64  `json:"executorId,omitempty"`
	Data              string  `json:"data,omitempty"`
	Spec              string  `json:"spec,omitempty"`
	StartDate         string  `json:"startDate,omitempty"`
	EndDate           string  `json:"endDate,omitempty"`
	LastExecutionDate string  `json:"lastExecutionDate,omitempty"`
	Timezone          string  `json:"timezone,omitempty"`
	TimezoneOffset    int64   `json:"timezoneOffset,omitempty"`
	RetryMax          int     `json:"retryMax,omitempty"`
	ExecutionID       string  `json:"executionId,omitempty"`
	Status            string  `json:"status,omitempty"`
	DateCreated       string  `json:"dateCreated,omitempty"`
	DateModified      *string `json:"dateModified,omitempty"`
	CreatedBy         string  `json:"createdBy,omitempty"`
	ModifiedBy        *string `json:"modifiedBy,omitempty"`
	DeletedBy         *string `json:"deletedBy,omitempty"`
}

// JobResponse represents the response for a single job
type JobResponse struct {
	Success bool `json:"success"`
	Data    Job  `json:"data"`
}

// JobRequestBody represents the request body for creating a job
type JobRequestBody struct {
	ProjectID      int64  `json:"projectId"`
	Timezone       string `json:"timezone"`
	ExecutorID     *int64 `json:"executorId,omitempty"`
	Data           string `json:"data,omitempty"`
	Spec           string `json:"spec,omitempty"`
	StartDate      string `json:"startDate,omitempty"`
	EndDate        string `json:"endDate,omitempty"`
	TimezoneOffset int64  `json:"timezoneOffset,omitempty"`
	RetryMax       int    `json:"retryMax,omitempty"`
	Status         string `json:"status,omitempty"`
}

// JobUpdateRequestBody represents the request body for updating a job
type JobUpdateRequestBody struct {
	ProjectID      int64  `json:"projectId,omitempty"`
	ExecutorID     *int64 `json:"executorId,omitempty"`
	Data           string `json:"data,omitempty"`
	Spec           string `json:"spec,omitempty"`
	StartDate      string `json:"startDate,omitempty"`
	EndDate        string `json:"endDate,omitempty"`
	Timezone       string `json:"timezone,omitempty"`
	TimezoneOffset int64  `json:"timezoneOffset,omitempty"`
	RetryMax       int    `json:"retryMax,omitempty"`
	Status         string `json:"status,omitempty"`
}

// PaginatedJobsResponse represents a paginated list of jobs
type PaginatedJobsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Total  int   `json:"total"`
		Offset int   `json:"offset"`
		Limit  int   `json:"limit"`
		Jobs   []Job `json:"jobs"`
	} `json:"data"`
}

// ListJobs retrieves all jobs with optional query parameters
func (c *Client) ListJobs(projectID string, limit, offset int, orderBy, orderByDirection string) (*PaginatedJobsResponse, error) {
	queryParams := map[string]string{
		"projectId": projectID,
		"limit":     fmt.Sprintf("%d", limit),
		"offset":    fmt.Sprintf("%d", offset),
	}

	if orderBy != "" {
		queryParams["orderBy"] = orderBy
	}
	if orderByDirection != "" {
		queryParams["orderByDirection"] = orderByDirection
	}

	req, err := c.newRequestWithQuery("GET", "/jobs", nil, queryParams)
	if err != nil {
		return nil, err
	}

	var result PaginatedJobsResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateJob creates a new job
func (c *Client) CreateJob(body *JobRequestBody) (*JobResponse, error) {
	req, err := c.newRequest("POST", "/jobs", body)
	if err != nil {
		return nil, err
	}

	var result JobResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetJob retrieves a single job by ID
func (c *Client) GetJob(id string) (*JobResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/jobs/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var result JobResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateJob updates an existing job
func (c *Client) UpdateJob(id string, body *JobUpdateRequestBody) (*JobResponse, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("/jobs/%s", id), body)
	if err != nil {
		return nil, err
	}

	var result JobResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteJob deletes a job by ID
func (c *Client) DeleteJob(id string) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/jobs/%s", id), nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// Account Management Methods

// CreateAccount creates a new account
func (c *Client) CreateAccount(body *AccountCreateRequestBody) (*AccountResponse, error) {
	req, err := c.newRequest("POST", "/accounts", body)
	if err != nil {
		return nil, err
	}

	var result AccountResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetAccount retrieves a single account by ID
func (c *Client) GetAccount(id string) (*AccountResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/accounts/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var result AccountResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AddFeatureToAccount adds a feature to an account
func (c *Client) AddFeatureToAccount(accountID string, body *FeatureRequest) (*FeatureRequestResponse, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("/accounts/%s/feature", accountID), body)
	if err != nil {
		return nil, err
	}

	var result FeatureRequestResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveFeatureFromAccount removes a feature from an account
func (c *Client) RemoveFeatureFromAccount(accountID string, body *FeatureRequest) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/accounts/%s/feature", accountID), body)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// Feature Management Methods

// ListFeatures retrieves all available features
func (c *Client) ListFeatures() (*FeaturesResponse, error) {
	req, err := c.newRequest("GET", "/features", nil)
	if err != nil {
		return nil, err
	}

	var result FeaturesResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// AsyncTask Methods

// GetAsyncTask retrieves an async task by request ID
func (c *Client) GetAsyncTask(requestID string) (*AsyncTaskResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/async-tasks/%s", requestID), nil)
	if err != nil {
		return nil, err
	}

	var result AsyncTaskResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RaftStats represents the Raft cluster statistics
type RaftStats struct {
	AppliedIndex        string `json:"applied_index"`
	CommitIndex         string `json:"commit_index"`
	FSMPending          string `json:"fsm_pending"`
	LastContact         string `json:"last_contact"`
	LastLogIndex        string `json:"last_log_index"`
	LastLogTerm         string `json:"last_log_term"`
	LastSnapshotIndex   string `json:"last_snapshot_index"`
	LastSnapshotTerm    string `json:"last_snapshot_term"`
	LatestConfiguration string `json:"latest_configuration"`
	LatestConfigIndex   string `json:"latest_configuration_index"`
	NumPeers            string `json:"num_peers"`
	ProtocolVersion     string `json:"protocol_version"`
	ProtocolVersionMax  string `json:"protocol_version_max"`
	ProtocolVersionMin  string `json:"protocol_version_min"`
	SnapshotVersionMax  string `json:"snapshot_version_max"`
	SnapshotVersionMin  string `json:"snapshot_version_min"`
	State               string `json:"state"`
	Term                string `json:"term"`
}

// HealthcheckData represents the healthcheck response data
type HealthcheckData struct {
	LeaderAddress string    `json:"leaderAddress"`
	LeaderID      string    `json:"leaderId"`
	RaftStats     RaftStats `json:"raftStats"`
}

// HealthcheckResponse represents the healthcheck response
type HealthcheckResponse struct {
	Success bool            `json:"success"`
	Data    HealthcheckData `json:"data"`
}

// Healthcheck retrieves the current leader and raft stats (no authentication required)
func (c *Client) Healthcheck() (*HealthcheckResponse, error) {
	// Create a request without authentication for healthcheck
	versionPrefix := fmt.Sprintf("/api/%s/", c.Version)
	rel := &url.URL{Path: path.Join(fmt.Sprintf("%s%s", c.BaseURL.Path, versionPrefix), "healthcheck")}
	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	var result HealthcheckResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
