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
}

func (client Client) Read(path string, key string) ReadResult {
	secret, err := client.api.Logical().Read(path)
	if err != nil {
		return ReadResult{"", nil, err}
	}

	if secret == nil {
		return ReadResult{"", nil, fmt.Errorf("could not find a secret on path %s", path)}
	}

	secrets, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return ReadResult{"", secret.Warnings, fmt.Errorf("could read secret data on path %s", path)}
	}

	value, ok := secrets[key].(string)
	if !ok {
		return ReadResult{"", secret.Warnings, fmt.Errorf("could not cast key %s on path %s to value", key, path)}
	}

	return ReadResult{value, secret.Warnings, nil}
}
