package bitcoin

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

type Bitcoin struct{}

func (b *Bitcoin) Name() string               { return "Bitcoin" }
func (b *Bitcoin) Symbol() string             { return "BTC" }
func (b *Bitcoin) CoinGeckoID() string        { return "bitcoin" }
func (b *Bitcoin) RequiredConfirmations() int { return 2 }

func (b *Bitcoin) rpcURL() string {
	return strings.TrimRight(os.Getenv("BTC_RPC_URL"), "/")
}

func (b *Bitcoin) getParams() *chaincfg.Params {
	url := strings.ToLower(b.rpcURL())
	if strings.Contains(url, "testnet") {
		return &chaincfg.TestNet3Params
	}
	if strings.Contains(url, "signet") {
		return &chaincfg.SigNetParams
	}
	if strings.Contains(url, "regtest") {
		return &chaincfg.RegressionNetParams
	}
	return &chaincfg.MainNetParams
}

// NewPaymentAddress derives a BIP84 P2WPKH address at m/84'/0'/accountIndex'/0/paymentIndex.
func (b *Bitcoin) NewPaymentAddress(masterSeed []byte, accountIndex uint32, paymentIndex uint32) (*cryptocurrency.PaymentAddress, error) {
	params := b.getParams()

	master, err := hdkeychain.NewMaster(masterSeed, params)
	if err != nil {
		return nil, fmt.Errorf("BTC master key: %w", err)
	}

	// m/84'/0'/accountIndex'/0/paymentIndex
	path := []uint32{
		hdkeychain.HardenedKeyStart + 84,
		hdkeychain.HardenedKeyStart + 0,
		hdkeychain.HardenedKeyStart + accountIndex,
		0,
		paymentIndex,
	}
	key := master
	for _, idx := range path {
		key, err = key.Derive(idx)
		if err != nil {
			return nil, fmt.Errorf("BTC derive: %w", err)
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

func (b *Bitcoin) CheckPayment(addr *cryptocurrency.PaymentAddress, requiredBTC float64) (*cryptocurrency.PaymentStatus, error) {
	url := b.rpcURL()
	if url == "" {
		return nil, fmt.Errorf("BTC_RPC_URL not set")
	}

	blockCount, err := rpc.Call[int](url, "getblockcount", nil)
	if err != nil {
		return nil, fmt.Errorf("getblockcount: %w", err)
	}

	desc := fmt.Sprintf("addr(%s)", addr.Address)
	if strings.HasPrefix(addr.Address, "bc1") || strings.HasPrefix(addr.Address, "tb1") {
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

func (b *Bitcoin) Price(currency string) (float64, error) {
	prices, err := coingecko.GetPrice([]string{b.CoinGeckoID()}, strings.ToLower(currency))
	if err != nil {
		return 0, err
	}
	p, ok := prices[b.CoinGeckoID()][strings.ToLower(currency)]
	if !ok {
		return 0, fmt.Errorf("no price for BTC/%s", currency)
	}
	return p, nil
}
