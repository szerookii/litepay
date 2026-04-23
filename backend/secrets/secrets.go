package secrets

import (
	"context"
	"crypto/sha512"
	"fmt"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

type Provider interface {
	Load(ctx context.Context) (string, error)
}

var masterMnemonic string

func Load(ctx context.Context) error {
	name := os.Getenv("SECRET_PROVIDER")
	if name == "" {
		name = "env"
	}

	p, err := newProvider(name)
	if err != nil {
		return err
	}

	mnemonic, err := p.Load(ctx)
	if err != nil {
		return fmt.Errorf("secrets(%s): %w", name, err)
	}
	if len(mnemonic) < 12 {
		return fmt.Errorf("secrets(%s): mnemonic missing or too short", name)
	}

	masterMnemonic = mnemonic
	return nil
}

func newProvider(name string) (Provider, error) {
	switch name {
	case "env":
		return &envProvider{}, nil
	case "vault":
		return &vaultProvider{}, nil
	case "bitwarden":
		return &bitwardenProvider{}, nil
	case "aws":
		return &awsProvider{}, nil
	case "gcp":
		return &gcpProvider{}, nil
	default:
		return nil, fmt.Errorf("unknown SECRET_PROVIDER %q (valid: env, vault, bitwarden, aws, gcp)", name)
	}
}

func DeriveMasterSeed() ([]byte, error) {
	if masterMnemonic == "" {
		return nil, fmt.Errorf("master seed not loaded — call secrets.Load first")
	}
	return pbkdf2.Key([]byte(masterMnemonic), []byte("mnemonic"), 2048, 64, sha512.New), nil
}
