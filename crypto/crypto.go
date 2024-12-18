package crypto

import (
	"github.com/szerookii/litepay/crypto/jsonrpc"
	"github.com/szerookii/litepay/crypto/litecoin"
	"strings"
	"sync"
)

type Blockchain interface {
	RequiredConfirmations() int
	Name() string
	Symbol() string
	Info() (*jsonrpc.BlockchainInfo, error)
	Price(string) (float64, error)
	ListWallets() ([]string, error)
	CreateWallet(string) error
	LoadWallet(string) error
	GetNewAddress(string) (string, error)
	GetAddressLabel(string) (string, error)
	ListUnspent(string) ([]*jsonrpc.Transaction, error)
}

var (
	cryptos = make(map[string]Blockchain)
	mutex   = new(sync.RWMutex)
)

func init() {
	cryptos["litecoin"] = new(litecoin.Litecoin)
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
