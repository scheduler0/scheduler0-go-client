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
}

func NewClient(baseURL, version, apiKey, apiSecret string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		BaseURL:    u,
		HTTPClient: &http.Client{},
		APIKey:     apiKey,
		APISecret:  apiSecret,
		Version:    version,
	}, nil
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
	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("X-API-Secret", c.APISecret)

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

func (c *Client) ListCredentials() (*PaginatedCredentialsResponse, error) {
	req, err := c.newRequest("GET", "/credentials", nil)
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

// ListExecutions retrieves all job executions
func (c *Client) ListExecutions() (*PaginatedExecutionsResponse, error) {
	req, err := c.newRequest("GET", "/executions", nil)
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
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Region           string `json:"region"`
	CloudProvider    string `json:"cloudProvider"`
	FilePath         string `json:"filePath"`
	CloudResourceURL string `json:"cloudResourceUrl"`
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
	FilePath         string `json:"filePath"`
	CloudResourceURL string `json:"cloudResourceUrl"`
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

// ListExecutors retrieves all executors
func (c *Client) ListExecutors() (*PaginatedExecutorsResponse, error) {
	req, err := c.newRequest("GET", "/executors", nil)
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

// Project represents a project
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"date_created"`
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

// ListProjects retrieves all projects
func (c *Client) ListProjects() (*PaginatedProjectsResponse, error) {
	req, err := c.newRequest("GET", "/projects", nil)
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
	ID          string `json:"id"`
	ProjectID   string `json:"projectId"`
	Description string `json:"description"`
	ExecutorID  string `json:"executorId"`
	Data        string `json:"data"`
	Spec        string `json:"spec"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Timezone    string `json:"timezone"`
}

// JobResponse represents the response for a single job
type JobResponse struct {
	Success bool `json:"success"`
	Data    Job  `json:"data"`
}

// JobRequestBody represents the request body for creating a job
type JobRequestBody struct {
	Description string `json:"description"`
	Timezone    string `json:"timezone"`
	CallbackURL string `json:"callback_url"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

// JobUpdateRequestBody represents the request body for updating a job
type JobUpdateRequestBody struct {
	ProjectID   string `json:"project_id"`
	Spec        string `json:"spec"`
	CallbackURL string `json:"callback_url"`
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

// ListJobs retrieves all jobs
func (c *Client) ListJobs() (*PaginatedJobsResponse, error) {
	req, err := c.newRequest("GET", "/jobs", nil)
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

// Healthcheck retrieves the current leader and raft stats
func (c *Client) Healthcheck() (*HealthcheckResponse, error) {
	req, err := c.newRequest("GET", "/healthcheck", nil)
	if err != nil {
		return nil, err
	}

	var result HealthcheckResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
