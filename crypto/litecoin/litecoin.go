package litecoin

import (
	"github.com/szerookii/litepay/crypto/jsonrpc"
	"github.com/szerookii/litepay/utils/coingecko"
	"strings"
)

type Litecoin struct{}

func (l *Litecoin) RequiredConfirmations() int {
	return 6
}

func (l *Litecoin) Name() string {
	return "Litecoin"
}

func (l *Litecoin) Symbol() string {
	return "LTC"
}

func (l *Litecoin) Info() (*jsonrpc.BlockchainInfo, error) {
	res, err := jsonrpc.CallRPC[*BlockchainInfo]("1.0", "getblockchaininfo", nil)
	if err != nil {
		return nil, err
	}

	return &jsonrpc.BlockchainInfo{
		Blocks:               int64(res.Result.Blocks),
		Headers:              int64(res.Result.Headers),
		Difficulty:           res.Result.Difficulty,
		BestBlockHash:        res.Result.Bestblockhash,
		VerificationProgress: res.Result.Verificationprogress,
		InitialBlockDownload: res.Result.Initialblockdownload,
	}, nil
}

func (l *Litecoin) Price(s string) (float64, error) {
	res, err := coingecko.GetPrice([]string{strings.ToLower(l.Name())}, s)
	if err != nil {
		return 0, err
	}

	return res[strings.ToLower(l.Name())][s], nil
}

func (l *Litecoin) ListWallets() ([]string, error) {
	res, err := jsonrpc.CallRPC[[]string]("1.0", "listwallets", nil)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (l *Litecoin) CreateWallet(name string) error {
	_, err := jsonrpc.CallRPC[*any]("1.0", "createwallet", []string{name})
	return err
}

func (l *Litecoin) LoadWallet(name string) error {
	_, err := jsonrpc.CallRPC[*any]("1.0", "loadwallet", []string{name})
	return err
}

func (l *Litecoin) GetNewAddress(label string) (string, error) {
	res, err := jsonrpc.CallRPC[string]("1.0", "getnewaddress", []string{label})
	if err != nil {
		return "", err
	}

	return res.Result, nil
}

func (l *Litecoin) GetAddressLabel(address string) (string, error) {
	res, err := jsonrpc.CallRPC[*AddressInfo]("1.0", "getaddressinfo", []string{address})
	if err != nil {
		return "", err
	}

	return res.Result.Label, nil
}

func (l *Litecoin) ListUnspent(address string) ([]*jsonrpc.Transaction, error) {
	res, err := jsonrpc.CallRPC[[]*jsonrpc.Transaction]("1.0", "listunspent", []any{0, 9999999, []string{address}})
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}
