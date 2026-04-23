package secrets

import (
	"context"
	"fmt"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type gcpProvider struct{}

func (p *gcpProvider) Load(ctx context.Context) (string, error) {
	project := os.Getenv("GCP_PROJECT_ID")
	name := os.Getenv("GCP_SECRET_NAME")
	version := os.Getenv("GCP_SECRET_VERSION")

	if project == "" || name == "" {
		return "", fmt.Errorf("GCP_PROJECT_ID and GCP_SECRET_NAME required")
	}
	if version == "" {
		version = "latest"
	}

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("gcp secret manager client: %w", err)
	}
	defer client.Close()

	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", project, name, version),
	})
	if err != nil {
		return "", fmt.Errorf("gcp access secret version: %w", err)
	}

	return string(result.Payload.Data), nil
}
