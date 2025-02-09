package vaultclient

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

// VaultClient struct
type VaultClient struct {
	client *api.Client
}

// NewVaultClient initializes a new Vault client with JWT authentication
func NewVaultClient(vaultAddr, mountpoint, role string, cloudflareHeaders map[string]string) (*VaultClient, error) {
	config := api.DefaultConfig()
	config.Address = vaultAddr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	// Add Cloudflare headers
	for key, value := range cloudflareHeaders {
		client.AddHeader(key, value)
	}

	// Authenticate using JWT
	authData := map[string]interface{}{
		"jwt":  os.Getenv("TERRAFORM_VAULT_AUTH_JWT"),
		"role": role,
	}

	authPath := fmt.Sprintf("auth/%s/login", mountpoint)
	secret, err := client.Logical().Write(authPath, authData)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Vault: %w", err)
	}

	// Set Vault token
	client.SetToken(secret.Auth.ClientToken)

	return &VaultClient{client: client}, nil
}

// ReadSecret retrieves a secret from Vault
func (v *VaultClient) ReadSecret(path string) (map[string]interface{}, error) {
	secret, err := v.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path: %s", path)
	}
	return secret.Data, nil
}

// WriteSecret stores a secret in Vault
func (v *VaultClient) WriteSecret(path string, data map[string]interface{}) error {
	_, err := v.client.Logical().Write(path, data)
	return err
}
