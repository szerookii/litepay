package cryptocurrency

import (
	"fmt"
	"github.com/szerookii/litepay/cryptocurrency/ethereum"
	"github.com/szerookii/litepay/cryptocurrency/types"
	"os"
	"strings"
	"sync"

	"github.com/phuslu/log"
	"github.com/szerookii/litepay/cryptocurrency/bitcoincash"
	"github.com/szerookii/litepay/cryptocurrency/litecoin"
)

type Blockchain interface {
	Name() string
	LegacyRPC() bool // if true, it's a legacy RPC, if false, it's a new RPC (like Ethereum that's doesn't need a password)
	Symbol() string
	RequiredConfirmations() int
	Synced() bool
	Price(string) (float64, error)
	EstimateFees() (float64, error)
	GetNewAddress(string) (string, error)
	RecentTransactions(string) ([]*types.Transaction, error)
}

var (
	cryptos = make(map[string]Blockchain)
	mutex   = new(sync.RWMutex)
)

func Init() {
	Add(new(litecoin.Litecoin))
	Add(new(bitcoincash.BitcoinCash))
	Add(new(ethereum.Ethereum))
}

func Add(crypto Blockchain) {
	mutex.Lock()
	defer mutex.Unlock()

	if os.Getenv(fmt.Sprintf("%s_RPC_HOST", crypto.Symbol())) == "" {
		log.Warn().Msgf("Missing %s_RPC_HOST environment variable, ignoring %s", crypto.Symbol(), crypto.Name())
		return
	}

	if crypto.LegacyRPC() {
		if os.Getenv(fmt.Sprintf("%s_RPC_USER", crypto.Symbol())) == "" {
			log.Error().Msgf("Missing %s_RPC_USER environment variable, ignoring %s", crypto.Symbol(), crypto.Name())
			return
		}

		if os.Getenv(fmt.Sprintf("%s_RPC_PASSWORD", crypto.Symbol())) == "" {
			log.Error().Msgf("Missing %s_RPC_PASSWORD environment variable, ignoring %s", crypto.Symbol(), crypto.Name())
			return
		}
	}

	if !crypto.Synced() {
		log.Warn().Msgf("%s is not synced, skipping", crypto.Name())
		return
	}

	cryptos[crypto.Name()] = crypto

	log.Info().Msgf("Enabled support for %s", crypto.Name())
}

func Get(name string) Blockchain {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, crypto := range cryptos {
		if strings.EqualFold(crypto.Name(), name) {
			return crypto
		}
	}

	return nil
}

func GetBySymbol(symbol string) Blockchain {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, crypto := range cryptos {
		if strings.EqualFold(crypto.Symbol(), symbol) {
			return crypto
		}
	}

	return nil
}
