package bitcoincash

import "github.com/szerookii/litepay/cryptocurrency/jsonrpc"

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
