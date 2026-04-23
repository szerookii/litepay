package db

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
	"github.com/szerookii/litepay/backend/ent"
	entuser "github.com/szerookii/litepay/backend/ent/user"
)

// GenerateAPIKey returns a cryptographically random 48-char hex API key.
func GenerateAPIKey() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func CreateUser(email, passwordHash, apiKey string) (*ent.User, error) {
	ctx := context.Background()
	tx, err := Client().Tx(ctx)
	if err != nil {
		return nil, err
	}

	// Assign account_index = current user count (unique, DB constraint enforces no race)
	count, err := tx.User.Query().Count(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	u, err := tx.User.Create().
		SetEmail(strings.ToLower(email)).
		SetPasswordHash(passwordHash).
		SetAPIKey(apiKey).
		SetAccountIndex(count).
		Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return u, tx.Commit()
}

func UserByEmail(email string) (*ent.User, error) {
	return Client().User.Query().
		Where(entuser.Email(strings.ToLower(email))).
		First(context.Background())
}

func UserByID(id uuid.UUID) (*ent.User, error) {
	return Client().User.Get(context.Background(), id)
}

func UserByAPIKey(apiKey string) (*ent.User, error) {
	return Client().User.Query().
		Where(entuser.APIKey(apiKey)).
		First(context.Background())
}

func UpdateUserWebhook(id uuid.UUID, webhookURL *string) (*ent.User, error) {
	q := Client().User.UpdateOneID(id)
	if webhookURL == nil || *webhookURL == "" {
		q.ClearWebhookURL()
	} else {
		q.SetWebhookURL(*webhookURL)
	}
	return q.Save(context.Background())
}

func RegenerateAPIKey(id uuid.UUID) (string, error) {
	key, err := GenerateAPIKey()
	if err != nil {
		return "", err
	}
	if _, err := Client().User.UpdateOneID(id).SetAPIKey(key).Save(context.Background()); err != nil {
		return "", err
	}
	return key, nil
}
