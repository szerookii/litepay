package litecoin

import (
	"github.com/phuslu/log"
	"github.com/szerookii/litepay/cryptocurrency/jsonrpc"
	"github.com/szerookii/litepay/cryptocurrency/types"
	"github.com/szerookii/litepay/utils/coingecko"
	"os"
	"strings"
)

type Litecoin struct{}

func (l *Litecoin) Name() string {
	return "Litecoin"
}

func (l *Litecoin) LegacyRPC() bool {
	return true
}

func (l *Litecoin) Symbol() string {
	return "LTC"
}

func (l *Litecoin) RequiredConfirmations() int {
	return 6
}

func (l *Litecoin) Synced() bool {
	res, err := jsonrpc.CallRPC[*BlockchainInfo]("1.0", "LTC", "getblockchaininfo", nil)
	if err != nil {
		return false
	}

	return res.Result.Verificationprogress > 0.98
}

func (l *Litecoin) Price(s string) (float64, error) {
	res, err := coingecko.GetPrice([]string{strings.ToLower(l.Name())}, s)
	if err != nil {
		return 0, err
	}

	return res[strings.ToLower(l.Name())][s], nil
}

func (l *Litecoin) EstimateFees() (float64, error) {
	// TODO: Understand how fees work and implement this function lol

	return 0, nil
}

func (l *Litecoin) GetNewAddress(label string) (string, error) {
	wallets, err := l.ListWallets()
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
		if err := l.LoadWallet(os.Getenv("WALLET_NAME")); err != nil {
			return "", err
		}
	}

	res, err := jsonrpc.CallRPC[string]("1.0", "LTC", "getnewaddress", []string{label})
	if err != nil {
		return "", err
	}

	return res.Result, nil
}

func (l *Litecoin) RecentTransactions(address string) ([]*types.Transaction, error) {
	res, err := jsonrpc.CallRPC[[]*jsonrpc.Transaction]("1.0", "LTC", "listunspent", []any{0, 9999999, []string{address}})
	if err != nil {
		return nil, err
	}

	var transactions []*types.Transaction
	for _, tx := range res.Result {
		if tx.Confirmations < l.RequiredConfirmations() {
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
