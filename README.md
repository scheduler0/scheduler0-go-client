# Scheduler0 Go Client

A Go client library for interacting with the [Scheduler0 API](https://scheduler0.com/api). This client provides a convenient way to manage credentials, executions, executors, projects, jobs, and monitor the health of your Scheduler0 cluster.

## Features

- **Credentials Management**
  - List credentials
  - Create new credentials
  - Get credential details
  - Update credentials
  - Delete credentials

- **Executions Management**
  - List job executions
  - View execution details and logs

- **Executors Management**
  - List executors
  - Create new executors
  - Get executor details
  - Update executors
  - Delete executors

- **Projects Management**
  - List projects
  - Create new projects
  - Get project details
  - Update projects
  - Delete projects

- **Jobs Management**
  - List jobs
  - Create new jobs
  - Get job details
  - Update jobs
  - Delete jobs

- **Health Monitoring**
  - Check cluster health
  - View raft statistics
  - Monitor leader status

## Installation

```bash
go get github.com/scheduler0/scheduler0-go-client
```

## Usage

### Creating a Client

```go
import "github.com/scheduler0/scheduler0-go-client"

client, err := scheduler0_go_client.NewClient(
    "http://localhost:7070",  // Base URL
    "v1",                    // API Version
    "your-api-key",          // API Key
    "your-api-secret",       // API Secret
)
if err != nil {
    log.Fatal(err)
}
```

### Managing Credentials

```go
// List credentials
credentials, err := client.ListCredentials()
if err != nil {
    log.Fatal(err)
}

// Create a new credential
credential, err := client.CreateCredential()
if err != nil {
    log.Fatal(err)
}

// Get a specific credential
credential, err := client.GetCredential("credential-id")
if err != nil {
    log.Fatal(err)
}

// Update a credential
credential, err := client.UpdateCredential("credential-id")
if err != nil {
    log.Fatal(err)
}

// Delete a credential
err := client.DeleteCredential("credential-id")
if err != nil {
    log.Fatal(err)
}
```

### Managing Executions

```go
// List executions
executions, err := client.ListExecutions()
if err != nil {
    log.Fatal(err)
}
```

### Managing Executors

```go
// List executors
executors, err := client.ListExecutors()
if err != nil {
    log.Fatal(err)
}

// Create a new executor
executor := &scheduler0_go_client.ExecutorRequestBody{
    Name:            "my-executor",
    Type:            "cloud_function",
    Region:          "us-west-1",
    CloudProvider:   "aws",
    FilePath:        "/path/to/function",
    CloudResourceURL: "https://example.com/function",
}
result, err := client.CreateExecutor(executor)
if err != nil {
    log.Fatal(err)
}

// Get a specific executor
executor, err := client.GetExecutor("executor-id")
if err != nil {
    log.Fatal(err)
}

// Update an executor
executor := &scheduler0_go_client.ExecutorRequestBody{
    Name: "updated-executor",
    // ... other fields
}
result, err := client.UpdateExecutor("executor-id", executor)
if err != nil {
    log.Fatal(err)
}

// Delete an executor
err := client.DeleteExecutor("executor-id")
if err != nil {
    log.Fatal(err)
}
```

### Managing Projects

```go
// List projects
projects, err := client.ListProjects()
if err != nil {
    log.Fatal(err)
}

// Create a new project
project := &scheduler0_go_client.ProjectRequestBody{
    Name:        "My Project",
    Description: "Project description",
}
result, err := client.CreateProject(project)
if err != nil {
    log.Fatal(err)
}

// Get a specific project
project, err := client.GetProject("project-id")
if err != nil {
    log.Fatal(err)
}

// Update a project
update := &scheduler0_go_client.ProjectUpdateRequestBody{
    Description: "Updated description",
}
result, err := client.UpdateProject("project-id", update)
if err != nil {
    log.Fatal(err)
}

// Delete a project
err := client.DeleteProject("project-id")
if err != nil {
    log.Fatal(err)
}
```

### Managing Jobs

```go
// List jobs
jobs, err := client.ListJobs()
if err != nil {
    log.Fatal(err)
}

// Create a new job
job := &scheduler0_go_client.JobRequestBody{
    Description: "My Job",
    Timezone:    "UTC",
    CallbackURL: "https://example.com/callback",
    StartDate:   "2024-01-01T00:00:00Z",
    EndDate:     "2024-12-31T23:59:59Z",
}
result, err := client.CreateJob(job)
if err != nil {
    log.Fatal(err)
}

// Get a specific job
job, err := client.GetJob("job-id")
if err != nil {
    log.Fatal(err)
}

// Update a job
update := &scheduler0_go_client.JobUpdateRequestBody{
    ProjectID:   "project-id",
    Spec:        "0 30 * * * *",
    CallbackURL: "https://example.com/updated-callback",
}
result, err := client.UpdateJob("job-id", update)
if err != nil {
    log.Fatal(err)
}

// Delete a job
err := client.DeleteJob("job-id")
if err != nil {
    log.Fatal(err)
}
```

### Health Monitoring

```go
// Check cluster health
health, err := client.Healthcheck()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Leader: %s\n", health.Data.LeaderAddress)
fmt.Printf("Raft State: %s\n", health.Data.RaftStats.State)
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. 