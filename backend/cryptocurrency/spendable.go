package cryptocurrency

// Spendable is implemented by coins that support signing and broadcasting transactions.
type Spendable interface {
	// DerivePrivKey returns the raw private key for a given HD path.
	DerivePrivKey(masterSeed []byte, accountIndex, paymentIndex uint32) ([]byte, error)
	// GetOnChainBalance returns the current spendable balance at addr.
	GetOnChainBalance(addr string) (float64, error)
	// SendFunds builds, signs, and broadcasts a transaction.
	// privKey is the raw private key bytes. fromAddr is used to locate UTXOs/balance.
	// Returns the broadcast tx hash.
	SendFunds(privKey []byte, fromAddr, toAddr string, amount float64) (string, error)
	// GetSender returns the originating address from a received transaction.
	// receiverAddr is our address, used to disambiguate sender vs receiver.
	GetSender(txHash, receiverAddr string) (string, error)
}
