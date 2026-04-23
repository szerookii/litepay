package litecoin

import (
	"fmt"
	"os"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/szerookii/litepay/backend/cryptocurrency"
	"github.com/szerookii/litepay/backend/cryptocurrency/rpc"
	"github.com/szerookii/litepay/backend/utils/coingecko"
)

type Litecoin struct{}

func (l *Litecoin) Name() string               { return "Litecoin" }
func (l *Litecoin) Symbol() string             { return "LTC" }
func (l *Litecoin) CoinGeckoID() string        { return "litecoin" }
func (l *Litecoin) RequiredConfirmations() int { return 6 }

func (l *Litecoin) rpcURL() string {
	return strings.TrimRight(os.Getenv("LTC_RPC_URL"), "/")
}

func (l *Litecoin) getParams() *chaincfg.Params {
	url := strings.ToLower(l.rpcURL())
	if strings.Contains(url, "testnet") {
		return &TestNet4Params
	}
	return &MainNetParams
}

// NewPaymentAddress derives a BIP84 P2WPKH address at m/84'/2'/accountIndex'/0/paymentIndex.
func (l *Litecoin) NewPaymentAddress(masterSeed []byte, accountIndex uint32, paymentIndex uint32) (*cryptocurrency.PaymentAddress, error) {
	params := l.getParams()

	master, err := hdkeychain.NewMaster(masterSeed, params)
	if err != nil {
		return nil, fmt.Errorf("LTC master key: %w", err)
	}

	// m/84'/2'/accountIndex'/0/paymentIndex
	path := []uint32{
		hdkeychain.HardenedKeyStart + 84,
		hdkeychain.HardenedKeyStart + 2, // coin type 2 = Litecoin
		hdkeychain.HardenedKeyStart + accountIndex,
		0,
		paymentIndex,
	}
	key := master
	for _, idx := range path {
		key, err = key.Derive(idx)
		if err != nil {
			return nil, fmt.Errorf("LTC derive: %w", err)
		}
	}

	pk, err := key.ECPubKey()
	if err != nil {
		return nil, err
	}
	addr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pk.SerializeCompressed()), params)
	if err != nil {
		return nil, err
	}
	return &cryptocurrency.PaymentAddress{Address: addr.EncodeAddress(), Index: paymentIndex}, nil
}

type scanUnspent struct {
	TxID   string  `json:"txid"`
	Height int     `json:"height"`
	Amount float64 `json:"amount"`
}

type scanResult struct {
	Success     bool          `json:"success"`
	TotalAmount float64       `json:"total_amount"`
	Unspents    []scanUnspent `json:"unspents"`
}

func (l *Litecoin) CheckPayment(addr *cryptocurrency.PaymentAddress, requiredLTC float64) (*cryptocurrency.PaymentStatus, error) {
	url := l.rpcURL()
	if url == "" {
		return nil, fmt.Errorf("LTC_RPC_URL not set")
	}

	blockCount, err := rpc.Call[int](url, "getblockcount", nil)
	if err != nil {
		return nil, fmt.Errorf("getblockcount: %w", err)
	}

	desc := fmt.Sprintf("addr(%s)", addr.Address)
	if strings.HasPrefix(addr.Address, "ltc1") || strings.HasPrefix(addr.Address, "tltc1") {
		desc = fmt.Sprintf("wpkh(%s)", addr.Address)
	}

	scan, err := rpc.Call[scanResult](url, "scantxoutset", []any{
		"start", []any{map[string]any{"desc": desc}},
	})
	if err != nil {
		return nil, fmt.Errorf("scantxoutset: %w", err)
	}

	if len(scan.Unspents) == 0 {
		return &cryptocurrency.PaymentStatus{}, nil
	}

	minConfs := -1
	var lastTxID string
	for _, u := range scan.Unspents {
		confs := blockCount - u.Height + 1
		if confs < 0 {
			confs = 0
		}
		if minConfs < 0 || confs < minConfs {
			minConfs = confs
		}
		lastTxID = u.TxID
	}
	if minConfs < 0 {
		minConfs = 0
	}

	return &cryptocurrency.PaymentStatus{
		ReceivedAmount: scan.TotalAmount,
		Confirmations:  minConfs,
		TxHash:         lastTxID,
	}, nil
}

func (l *Litecoin) Price(currency string) (float64, error) {
	prices, err := coingecko.GetPrice([]string{l.CoinGeckoID()}, strings.ToLower(currency))
	if err != nil {
		return 0, err
	}
	p, ok := prices[l.CoinGeckoID()][strings.ToLower(currency)]
	if !ok {
		return 0, fmt.Errorf("no price for LTC/%s", currency)
	}
	return p, nil
}
