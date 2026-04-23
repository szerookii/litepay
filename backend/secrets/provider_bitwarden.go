package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type bitwardenProvider struct{}

func (p *bitwardenProvider) Load(ctx context.Context) (string, error) {
	clientID := os.Getenv("BITWARDEN_CLIENT_ID")
	clientSecret := os.Getenv("BITWARDEN_CLIENT_SECRET")
	secretID := os.Getenv("BITWARDEN_SECRET_ID")
	identityURL := os.Getenv("BITWARDEN_IDENTITY_URL")
	apiURL := os.Getenv("BITWARDEN_API_URL")

	if clientID == "" || clientSecret == "" || secretID == "" {
		return "", fmt.Errorf("BITWARDEN_CLIENT_ID, BITWARDEN_CLIENT_SECRET, BITWARDEN_SECRET_ID required")
	}
	if identityURL == "" {
		identityURL = "https://identity.bitwarden.com"
	}
	if apiURL == "" {
		apiURL = "https://api.bitwarden.com"
	}

	tokenURL := strings.TrimRight(identityURL, "/") + "/connect/token"
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"scope":         {"api.secrets"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("bitwarden auth: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bitwarden auth returned HTTP %d", resp.StatusCode)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("bitwarden auth decode: %w", err)
	}

	secretURL := strings.TrimRight(apiURL, "/") + "/secrets/" + secretID
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, secretURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("bitwarden fetch secret: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bitwarden secret returned HTTP %d", resp.StatusCode)
	}

	var secretResp struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&secretResp); err != nil {
		return "", fmt.Errorf("bitwarden secret decode: %w", err)
	}

	return secretResp.Value, nil
}
