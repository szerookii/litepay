package ethereum

import (
	"context"
	"github.com/szerookii/litepay/cryptocurrency/jsonrpc"
	"github.com/szerookii/litepay/cryptocurrency/types"
	"github.com/szerookii/litepay/utils/coingecko"
	"strings"
)

type Ethereum struct{}

func (b *Ethereum) Name() string {
	return "Ethereum"
}

func (b *Ethereum) LegacyRPC() bool {
	return false
}

func (b *Ethereum) Symbol() string {
	return "ETH"
}

func (b *Ethereum) RequiredConfirmations() int {
	return 1 // 1 for dev and 3 or 6 for mainnet
}

func (b *Ethereum) Synced() bool {
	client, err := Client()
	if err != nil {
		return false
	}

	res, err := client.SyncProgress(context.Background())
	if err != nil {
		return false
	}

	if res == nil {
		return true // there is no sync progress, so it's synced, right?
	}

	return res.CurrentBlock >= res.HighestBlock
}

func (b *Ethereum) Price(s string) (float64, error) {
	res, err := coingecko.GetPrice([]string{strings.ToLower(b.Name())}, s)
	if err != nil {
		return 0, err
	}

	return res[strings.ToLower(b.Name())][s], nil
}

func (b *Ethereum) EstimateFees() (float64, error) {
	// TODO: Understand how fees work and implement this function lol

	return 0, nil
}

func (b *Ethereum) ListWallets() ([]string, error) {
	res, err := jsonrpc.CallRPC[[]string]("1.0", "LTC", "listwallets", nil)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (b *Ethereum) CreateWallet(name string) error {
	_, err := jsonrpc.CallRPC[*any]("1.0", "LTC", "createwallet", []string{name})
	return err
}

func (b *Ethereum) LoadWallet(name string) error {
	_, err := jsonrpc.CallRPC[*any]("1.0", "LTC", "loadwallet", []string{name})
	return err
}

func (b *Ethereum) GetNewAddress(label string) (string, error) {

	res, err := jsonrpc.CallRPC[string]("1.0", "LTC", "getnewaddress", []string{label})
	if err != nil {
		return "", err
	}

	return res.Result, nil
}

func (b *Ethereum) RecentTransactions(s string) ([]*types.Transaction, error) {
	//TODO implement me
	panic("implement me")
}
