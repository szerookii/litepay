package db

import (
	"github.com/szerookii/litepay/prisma/db"
	"os"
	"sync"
)

var (
	client *db.PrismaClient
	dbMu   = sync.RWMutex{}
)

func Client() *db.PrismaClient {
	dbMu.Lock()
	defer dbMu.Unlock()

	if client == nil {
		client = db.NewClient(db.WithDatasourceURL(os.Getenv("DATABASE_URL")))

		if client == nil {
			panic("Failed to connect to the database.")
		}

		err := client.Connect()
		if err != nil {
			panic(err)
		}
	}

	return client
}

func Close() {
	dbMu.RLock()
	defer dbMu.RUnlock()

	if client != nil {
		client.Disconnect()
	}
}
