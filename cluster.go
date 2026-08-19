package scheduler0_go_client

import "fmt"

// RemoveSelfFromCluster removes this node from Raft membership and unregisters it
// from etcd. Intended to be called by a node against itself during shutdown.
func (c *Client) RemoveSelfFromCluster() (*ClusterStatusResponse, error) {
	req, err := c.newRequest("POST", "/cluster/remove-self", nil, "")
	if err != nil {
		return nil, err
	}

	var result ClusterStatusResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddSelfToCluster ensures this node is registered in etcd and part of the Raft
// cluster. Intended to be called by a node against itself during startup.
func (c *Client) AddSelfToCluster() (*ClusterStatusResponse, error) {
	req, err := c.newRequest("POST", "/cluster/add-self", nil, "")
	if err != nil {
		return nil, err
	}

	var result ClusterStatusResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ForceRebuildCluster forces a rebuild of the Raft cluster from seedNodeID. This
// should only be called on the seed node.
func (c *Client) ForceRebuildCluster(seedNodeID uint64) (*ClusterStatusResponse, error) {
	queryParams := map[string]string{
		"seedNodeId": fmt.Sprintf("%d", seedNodeID),
	}
	req, err := c.newRequestWithQuery("POST", "/cluster/force-rebuild", nil, queryParams, "")
	if err != nil {
		return nil, err
	}

	var result ClusterStatusResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ResetRaftState clears local Raft state on the target node and causes it to exit
// the process. The server sends its response before exiting, but the connection may
// still be reset if the process exits before the response is fully flushed to the
// client; callers should treat a connection error here as possibly-successful.
func (c *Client) ResetRaftState() (*ClusterStatusResponse, error) {
	req, err := c.newRequest("POST", "/cluster/reset-raft", nil, "")
	if err != nil {
		return nil, err
	}

	var result ClusterStatusResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveNode removes nodeID from the Raft cluster. Only the leader can perform this
// operation; the server responds 403 if called against a non-leader node.
func (c *Client) RemoveNode(nodeID uint64) (*ClusterStatusResponse, error) {
	queryParams := map[string]string{
		"nodeId": fmt.Sprintf("%d", nodeID),
	}
	req, err := c.newRequestWithQuery("POST", "/cluster/remove-node", nil, queryParams, "")
	if err != nil {
		return nil, err
	}

	var result ClusterStatusResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddNode adds a node to the Raft cluster. Only the leader can perform this
// operation; the server responds 403 if called against a non-leader node.
func (c *Client) AddNode(nodeID uint64, nodeAddress string, clientAddress string) (*ClusterStatusResponse, error) {
	queryParams := map[string]string{
		"nodeId":        fmt.Sprintf("%d", nodeID),
		"nodeAddress":   nodeAddress,
		"clientAddress": clientAddress,
	}
	req, err := c.newRequestWithQuery("POST", "/cluster/add-node", nil, queryParams, "")
	if err != nil {
		return nil, err
	}

	var result ClusterStatusResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PromoteNode promotes a non-voter node to a voter in the Raft cluster. Only the
// leader can perform this operation; the server responds 403 if called against a
// non-leader node.
func (c *Client) PromoteNode(nodeID uint64) (*ClusterStatusResponse, error) {
	queryParams := map[string]string{
		"nodeId": fmt.Sprintf("%d", nodeID),
	}
	req, err := c.newRequestWithQuery("POST", "/cluster/promote-node", nil, queryParams, "")
	if err != nil {
		return nil, err
	}

	var result ClusterStatusResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DemoteNode demotes a voter node to a non-voter in the Raft cluster. Only the
// leader can perform this operation; the server responds 403 if called against a
// non-leader node.
func (c *Client) DemoteNode(nodeID uint64) (*ClusterStatusResponse, error) {
	queryParams := map[string]string{
		"nodeId": fmt.Sprintf("%d", nodeID),
	}
	req, err := c.newRequestWithQuery("POST", "/cluster/demote-node", nil, queryParams, "")
	if err != nil {
		return nil, err
	}

	var result ClusterStatusResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TransferLeadership transfers Raft leadership to another voter node. Only the
// leader can perform this operation; the server responds 403 if called against a
// non-leader node.
func (c *Client) TransferLeadership() (*ClusterStatusResponse, error) {
	req, err := c.newRequest("POST", "/cluster/transfer-leadership", nil, "")
	if err != nil {
		return nil, err
	}

	var result ClusterStatusResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListNodes returns every node currently in the Raft cluster.
func (c *Client) ListNodes() (*ListNodesResponse, error) {
	req, err := c.newRequest("GET", "/cluster/list-nodes", nil, "")
	if err != nil {
		return nil, err
	}

	var result ListNodesResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DumpScheduleQueue returns the in-memory schedule queue for debugging. See
// ClusterDumpResponse for why Data is left as raw JSON.
func (c *Client) DumpScheduleQueue() (*ClusterDumpResponse, error) {
	req, err := c.newRequest("GET", "/cluster/dump/schedule-queue", nil, "")
	if err != nil {
		return nil, err
	}

	var result ClusterDumpResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DumpJobExecutionsCache returns the in-memory job executions cache for debugging.
// See ClusterDumpResponse for why Data is left as raw JSON.
func (c *Client) DumpJobExecutionsCache() (*ClusterDumpResponse, error) {
	req, err := c.newRequest("GET", "/cluster/dump/job-executions-cache", nil, "")
	if err != nil {
		return nil, err
	}

	var result ClusterDumpResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DumpJobQueues returns every persisted job queue for debugging. See
// ClusterDumpResponse for why Data is left as raw JSON.
func (c *Client) DumpJobQueues() (*ClusterDumpResponse, error) {
	req, err := c.newRequest("GET", "/cluster/dump/job-queues", nil, "")
	if err != nil {
		return nil, err
	}

	var result ClusterDumpResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DumpJobQueueVersions returns every persisted job queue version for debugging. See
// ClusterDumpResponse for why Data is left as raw JSON.
func (c *Client) DumpJobQueueVersions() (*ClusterDumpResponse, error) {
	req, err := c.newRequest("GET", "/cluster/dump/job-queue-versions", nil, "")
	if err != nil {
		return nil, err
	}

	var result ClusterDumpResponse
	if err = c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
