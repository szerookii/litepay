package cryptocurrency

import "sync"

// PaymentAddress is a generated deposit destination.
type PaymentAddress struct {
	Address   string // derived deposit address
	Reference string // SOL: unique reference pubkey; BTC/LTC: empty
	Index     uint32 // payment index used for derivation
}

// PaymentStatus holds the current on-chain state for a payment address.
type PaymentStatus struct {
	ReceivedAmount float64
	Confirmations  int
	TxHash         string
	IsPending      bool // tx seen in mempool but not yet confirmed
}

// Blockchain is implemented by each supported coin.
type Blockchain interface {
	Name() string
	Symbol() string
	CoinGeckoID() string
	RequiredConfirmations() int
	// NewPaymentAddress derives a unique deposit address from the master seed.
	// accountIndex is the user's unique HD account slot; paymentIndex is per-user per-coin.
	NewPaymentAddress(masterSeed []byte, accountIndex uint32, paymentIndex uint32) (*PaymentAddress, error)
	// CheckPayment returns the current on-chain status for a payment address.
	CheckPayment(addr *PaymentAddress, requiredAmount float64) (*PaymentStatus, error)
	Price(currency string) (float64, error)
}

var (
	mu       sync.RWMutex
	registry = map[string]Blockchain{}
)

func Register(b Blockchain) {
	mu.Lock()
	defer mu.Unlock()
	registry[b.Symbol()] = b
}

func GetBySymbol(symbol string) (Blockchain, bool) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := registry[symbol]
	return b, ok
}

func All() []Blockchain {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]Blockchain, 0, len(registry))
	for _, b := range registry {
		out = append(out, b)
	}
	return out
}
