package secrets

import (
	"context"
	"fmt"
	"os"

	vault "github.com/hashicorp/vault-client-go"
)

type vaultProvider struct{}

func (p *vaultProvider) Load(ctx context.Context) (string, error) {
	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	mount := os.Getenv("VAULT_MOUNT")
	path := os.Getenv("VAULT_PATH")
	key := os.Getenv("VAULT_KEY")

	if addr == "" || token == "" || path == "" || key == "" {
		return "", fmt.Errorf("VAULT_ADDR, VAULT_TOKEN, VAULT_PATH, VAULT_KEY all required")
	}
	if mount == "" {
		mount = "secret"
	}

	client, err := vault.New(vault.WithAddress(addr))
	if err != nil {
		return "", fmt.Errorf("vault client init: %w", err)
	}
	if err := client.SetToken(token); err != nil {
		return "", fmt.Errorf("vault set token: %w", err)
	}

	resp, err := client.Secrets.KvV2Read(ctx, path, vault.WithMountPath(mount))
	if err != nil {
		return "", fmt.Errorf("vault kv read: %w", err)
	}

	raw, ok := resp.Data.Data[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in vault secret at %s/%s", key, mount, path)
	}

	val, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("vault key %q is not a string", key)
	}
	return val, nil
}
