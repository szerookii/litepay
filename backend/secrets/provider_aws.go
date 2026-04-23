package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type awsProvider struct{}

func (p *awsProvider) Load(ctx context.Context) (string, error) {
	region := os.Getenv("AWS_REGION")
	secretID := os.Getenv("AWS_SECRET_ID")
	jsonKey := os.Getenv("AWS_SECRET_KEY")

	if region == "" || secretID == "" {
		return "", fmt.Errorf("AWS_REGION and AWS_SECRET_ID required")
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("aws config: %w", err)
	}

	client := secretsmanager.NewFromConfig(cfg)
	result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	})
	if err != nil {
		return "", fmt.Errorf("aws get secret: %w", err)
	}

	if result.SecretString == nil {
		return "", fmt.Errorf("aws secret %q has no SecretString (binary secrets not supported)", secretID)
	}

	value := *result.SecretString

	if jsonKey != "" {
		var m map[string]string
		if err := json.Unmarshal([]byte(value), &m); err != nil {
			return "", fmt.Errorf("AWS_SECRET_KEY set but SecretString is not valid JSON: %w", err)
		}
		v, ok := m[jsonKey]
		if !ok {
			return "", fmt.Errorf("key %q not found in AWS secret JSON", jsonKey)
		}
		value = v
	}

	return value, nil
}
