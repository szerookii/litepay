package crypto

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/phuslu/log"
	"github.com/szerookii/litepay/crypto/bitcoincash"
	"github.com/szerookii/litepay/crypto/jsonrpc"
	"github.com/szerookii/litepay/crypto/litecoin"
)

type Blockchain interface {
	RequiredConfirmations() int
	Name() string
	Symbol() string
	Info() (*jsonrpc.BlockchainInfo, error)
	Price(string) (float64, error)
	EstimateFees() (float64, error)
	ListWallets() ([]string, error)
	CreateWallet(string) error
	LoadWallet(string) error
	GetNewAddress(string) (string, error)
	ListUnspent(string) ([]*jsonrpc.Transaction, error)
}

var (
	cryptos = make(map[string]Blockchain)
	mutex   = new(sync.RWMutex)
)

func Init() {
	Add(new(litecoin.Litecoin))
	Add(new(bitcoincash.BitcoinCash))
}

func Add(crypto Blockchain) {
	mutex.Lock()
	defer mutex.Unlock()

	if os.Getenv(fmt.Sprintf("%s_RPC_HOST", crypto.Symbol())) == "" {
		log.Info().Msgf("missing %s_RPC_HOST environment variable, ignoring %s", crypto.Symbol(), crypto.Name())
		return
	}

	if os.Getenv(fmt.Sprintf("%s_RPC_USER", crypto.Symbol())) == "" {
		log.Info().Msgf("missing %s_RPC_USER environment variable, ignoring %s", crypto.Symbol(), crypto.Name())
		return
	}

	if os.Getenv(fmt.Sprintf("%s_RPC_PASSWORD", crypto.Symbol())) == "" {
		log.Info().Msgf("missing %s_RPC_PASSWORD environment variable, ignoring %s", crypto.Symbol(), crypto.Name())
		return
	}

	_, err := crypto.Info()
	if err != nil {
		log.Error().Err(err).Msgf("failed to get blockchain info for %s", crypto.Name())
		return
	}

	cryptos[crypto.Name()] = crypto

	log.Info().Msgf("Enabled support for %s.", crypto.Name())
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
