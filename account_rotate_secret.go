package scheduler0_go_client

// RotateSecretRequest is the body for POST /account/rotate-secret. OldSecretKey is the
// previous SecretKey, needed to decrypt existing rows before re-encrypting them with the
// server's currently-loaded (new) SecretKey.
type RotateSecretRequest struct {
	OldSecretKey string `json:"oldSecretKey"`
}

// RotateSecretResponse is the response body returned by POST /account/rotate-secret.
// The counts are the number of rows re-encrypted in each subsystem.
type RotateSecretResponse struct {
	Success bool `json:"success"`
	Data    struct {
		CredentialsRotated uint64 `json:"credentialsRotated"`
		ExecutorsRotated   uint64 `json:"executorsRotated"`
		AISettingsRotated  uint64 `json:"aiSettingsRotated"`
	} `json:"data"`
}

// RotateSecret triggers a server-side re-encryption of all secrets stored under the
// server's SecretKey — credential api secrets, executor cloud provider credentials, and
// per-account AI provider keys — from oldSecretKey to the server's currently-loaded
// SecretKey.
//
// A credential's api_secret is stored encrypted and verified by decrypt-then-compare, so
// it is re-encrypted too; the api_key is a stable opaque identifier that is left unchanged,
// so rotation does not invalidate any client's credential.
//
// This endpoint is self-hosting only and requires a BasicAuth client
// (NewBasicAuthClient). The operator must update SecretKey in the secrets source (and
// reload/restart the server) before calling this method, then pass the previous key as
// oldSecretKey.
func (c *Client) RotateSecret(oldSecretKey string) (*RotateSecretResponse, error) {
	req, err := c.newRequest("POST", "/account/rotate-secret", RotateSecretRequest{OldSecretKey: oldSecretKey})
	if err != nil {
		return nil, err
	}

	var result RotateSecretResponse
	if err := c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
