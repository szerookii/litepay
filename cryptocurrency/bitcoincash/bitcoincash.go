package bitcoincash

import (
	"github.com/phuslu/log"
	"github.com/szerookii/litepay/cryptocurrency/jsonrpc"
	"github.com/szerookii/litepay/cryptocurrency/types"
	"github.com/szerookii/litepay/utils/coingecko"
	"os"
	"strings"
)

type BitcoinCash struct{}

func (b *BitcoinCash) Name() string {
	return "Bitcoin-Cash"
}

func (b *BitcoinCash) LegacyRPC() bool {
	return true
}

func (b *BitcoinCash) Symbol() string {
	return "BCH"
}

func (b *BitcoinCash) RequiredConfirmations() int {
	return 1 // 1 for dev and 3 or 6 for mainnet
}

func (b *BitcoinCash) Synced() bool {
	res, err := jsonrpc.CallRPC[*jsonrpc.BlockchainInfo]("1.0", "BCH", "getblockchaininfo", nil)
	if err != nil {
		return false
	}

	return res.Result.VerificationProgress > 0.98
}

func (b *BitcoinCash) Price(s string) (float64, error) {
	res, err := coingecko.GetPrice([]string{strings.ToLower(b.Name())}, s)
	if err != nil {
		return 0, err
	}

	return res[strings.ToLower(b.Name())][s], nil
}

func (b *BitcoinCash) EstimateFees() (float64, error) {
	// TODO: Understand how fees work and implement this function lol

	return 0, nil
}

func (b *BitcoinCash) GetNewAddress(label string) (string, error) {
	wallets, err := b.ListWallets()
	if err != nil {
		return "", err
	}

	var walletLoaded bool
	for _, wallet := range wallets {
		if wallet == os.Getenv("WALLET_NAME") {
			walletLoaded = true
			break
		}
	}

	if !walletLoaded {
		log.Info().Msgf("wallet %s not found, loading...", os.Getenv("WALLET_NAME"))
		if err := b.LoadWallet(os.Getenv("WALLET_NAME")); err != nil {
			return "", err
		}
	}

	res, err := jsonrpc.CallRPC[string]("1.0", "LTC", "getnewaddress", []string{label})
	if err != nil {
		return "", err
	}

	return res.Result, nil
}

func (b *BitcoinCash) RecentTransactions(address string) ([]*types.Transaction, error) {
	res, err := jsonrpc.CallRPC[[]*jsonrpc.Transaction]("1.0", "LTC", "listunspent", []any{0, 9999999, []string{address}})
	if err != nil {
		return nil, err
	}

	var transactions []*types.Transaction
	for _, tx := range res.Result {
		if tx.Confirmations < b.RequiredConfirmations() {
			continue
		}

		transactions = append(transactions, &types.Transaction{
			Txid:          tx.Txid,
			Amount:        tx.Amount,
			Confirmations: tx.Confirmations,
		})
	}

	return transactions, nil
}
