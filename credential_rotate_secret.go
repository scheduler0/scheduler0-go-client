package scheduler0_go_client

// RotateSecretResponse is the response body returned by POST /credentials/rotate-secret.
// Rotated is the number of credentials that were re-encrypted with the new secret key.
type RotateSecretResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Rotated uint64 `json:"rotated"`
	} `json:"data"`
}

// RotateCredentialSecret triggers a server-side re-encryption of all active,
// non-expired credentials from the old SecretKey to the new one.
//
// This endpoint is self-hosting only and requires a BasicAuth client
// (NewBasicAuthClient). The operator must update SecretKey in the secrets
// source before calling this method. The operation is resumable — if it is
// interrupted, re-calling RotateCredentialSecret will pick up from where it
// left off.
func (c *Client) RotateCredentialSecret() (*RotateSecretResponse, error) {
	req, err := c.newRequest("POST", "/credentials/rotate-secret", nil)
	if err != nil {
		return nil, err
	}

	var result RotateSecretResponse
	if err := c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
