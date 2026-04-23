package secrets

import (
	"context"
	"os"
)

type envProvider struct{}

func (p *envProvider) Load(_ context.Context) (string, error) {
	val := os.Getenv("MASTER_SEED")
	os.Unsetenv("MASTER_SEED")
	return val, nil
}
