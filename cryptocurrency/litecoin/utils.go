package litecoin

import "github.com/szerookii/litepay/cryptocurrency/jsonrpc"

func (l *Litecoin) ListWallets() ([]string, error) {
	res, err := jsonrpc.CallRPC[[]string]("1.0", "LTC", "listwallets", nil)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (l *Litecoin) CreateWallet(name string) error {
	_, err := jsonrpc.CallRPC[*any]("1.0", "LTC", "createwallet", []string{name})
	return err
}

func (l *Litecoin) LoadWallet(name string) error {
	_, err := jsonrpc.CallRPC[*any]("1.0", "LTC", "loadwallet", []string{name})
	return err
}
