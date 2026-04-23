package db

import (
	"os"
	"sync"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/szerookii/litepay/backend/ent"
)

var (
	client *ent.Client
	dbMu   sync.RWMutex
)

func mustParseDSN(dsn string) *pgx.ConnConfig {
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}
	return config
}

func Client() *ent.Client {
	dbMu.Lock()
	defer dbMu.Unlock()

	if client == nil {
		sqlDB := stdlib.OpenDB(*mustParseDSN(os.Getenv("DATABASE_URL")))
		drv := entsql.OpenDB(dialect.Postgres, sqlDB)
		client = ent.NewClient(ent.Driver(drv))
	}

	return client
}

func Close() {
	dbMu.RLock()
	defer dbMu.RUnlock()

	if client != nil {
		client.Close()
	}
}

