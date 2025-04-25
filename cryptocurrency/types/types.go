package types

type BlockchainInfo struct {
	Synced bool `json:"synced"`
}

type Transaction struct {
	Txid          string  `json:"txid"`
	Amount        float64 `json:"amount"`
	Confirmations int     `json:"confirmations"`
}
