package client

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Client is the core of vault-to-env
// responsible for communication with Vault and other app logic
type Client struct {
	api *api.Client
}

// Config is used to parameterize Client
type Config struct {
	VaultAddress string
	VaultToken   string
}

// NewClient is used to create a new Client
func NewClient(config *Config) (*Client, error) {
	api, err := api.NewClient(&api.Config{
		Address: config.VaultAddress,
	})
	if err != nil {
		return nil, err
	}

	api.SetToken(config.VaultToken)

	return &Client{
		api: api,
	}, nil
}

// ReadResult represents a result of reading a secret
type ReadResult struct {
	Value    string
	Warnings []string
	Error    error
	LeaseID  string
}

// getSecretDataKey extracts a key as a string from a given secret
func (client Client) getSecretDataKey(secret *api.Secret, key []string) (string, error) {
	var intermediate = secret.Data

	for i, k := range key {
		if i == len(key)-1 {
			value, ok := intermediate[k].(string)

			if !ok {
				return "", fmt.Errorf(`could not cast secret key "%s" to string`, key)
			}

			return value, nil
		}

		m, ok := intermediate[k].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("could not cast secret data to map[string]interface{}")
		}

		intermediate = m
	}

	return "", nil
}

// SetToken is used to set auth token
func (client Client) SetToken(token string) {
	client.api.SetToken(token)
}

// Auth is used to perform generic authention against a given path
func (client Client) Auth(path string, data map[string]interface{}) (string, error) {
	secret, err := client.api.Logical().Write(path, data)
	if err != nil {
		return "", err
	}

	return secret.Auth.ClientToken, nil
}

// Read is used to read secrets
func (client Client) Read(path string, key []string) ReadResult {
	secret, err := client.api.Logical().Read(path)
	if err != nil {
		return ReadResult{"", nil, err, ""}
	}

	if secret == nil {
		return ReadResult{"", nil, fmt.Errorf("could not find a secret on path %s", path), ""}
	}

	value, err := client.getSecretDataKey(secret, key)

	return ReadResult{value, secret.Warnings, err, secret.LeaseID}
}

func (client Client) RenewLease(LeaseID string, increment int) error {
	var payload = make(map[string]interface{})

	payload["lease_id"] = LeaseID
	payload["increment"] = increment

	_, err := client.api.Logical().Write("/sys/leases/renew", payload)
	if err != nil {
		return err
	}

	return nil
}
