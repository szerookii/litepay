package bitcoincash

import (
	"github.com/szerookii/litepay/crypto/jsonrpc"
	"github.com/szerookii/litepay/utils/coingecko"
	"strings"
)

type BitcoinCash struct{}

func (b *BitcoinCash) RequiredConfirmations() int {
	return 1 // 1 for dev and 3 or 6 for mainnet
}

func (b *BitcoinCash) Name() string {
	return "Bitcoin-Cash"
}

func (b *BitcoinCash) Symbol() string {
	return "BCH"
}

func (b *BitcoinCash) Info() (*jsonrpc.BlockchainInfo, error) {
	res, err := jsonrpc.CallRPC[*jsonrpc.BlockchainInfo]("1.0", "BCH", "getblockchaininfo", nil)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
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

func (b *BitcoinCash) ListWallets() ([]string, error) {
	res, err := jsonrpc.CallRPC[[]string]("1.0", "LTC", "listwallets", nil)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (b *BitcoinCash) CreateWallet(name string) error {
	_, err := jsonrpc.CallRPC[*any]("1.0", "LTC", "createwallet", []string{name})
	return err
}

func (b *BitcoinCash) LoadWallet(name string) error {
	_, err := jsonrpc.CallRPC[*any]("1.0", "LTC", "loadwallet", []string{name})
	return err
}

func (b *BitcoinCash) GetNewAddress(label string) (string, error) {
	res, err := jsonrpc.CallRPC[string]("1.0", "LTC", "getnewaddress", []string{label})
	if err != nil {
		return "", err
	}

	return res.Result, nil
}

func (b *BitcoinCash) ListUnspent(address string) ([]*jsonrpc.Transaction, error) {
	res, err := jsonrpc.CallRPC[[]*jsonrpc.Transaction]("1.0", "LTC", "listunspent", []any{0, 9999999, []string{address}})
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}
