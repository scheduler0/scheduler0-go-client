package scheduler0_go_client

import "fmt"

// TestInvokeExecutor fires a synthetic ("test") job through an existing executor
// immediately and synchronously, so a scheduled job can be exercised without
// waiting for its cron spec or start date to elapse. The invocation has no side
// effects: no job is created, no execution log is written, and nothing is
// scheduled or rescheduled.
//
// The body is optional (pass nil to invoke with a default synthetic job). The
// endpoint returns HTTP 200 once the invocation completes; whether the target
// accepted the job is reported in the response's Data.Success field. Local
// (pull-based) executors cannot be test-invoked and return a 400 error.
func (c *Client) TestInvokeExecutor(id string, body *TestInvocationRequestBody) (*TestInvocationResponse, error) {
	req, err := c.newRequest("POST", fmt.Sprintf("/executors/%s/test-invoke", id), body)
	if err != nil {
		return nil, err
	}

	var result TestInvocationResponse
	err = c.do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
