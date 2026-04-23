package litecoin

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/szerookii/litepay/backend/cryptocurrency/rpc"
)

func (l *Litecoin) DerivePrivKey(masterSeed []byte, accountIndex, paymentIndex uint32) ([]byte, error) {
	params := l.getParams()
	master, err := hdkeychain.NewMaster(masterSeed, params)
	if err != nil {
		return nil, fmt.Errorf("LTC master key: %w", err)
	}
	path := []uint32{
		hdkeychain.HardenedKeyStart + 84,
		hdkeychain.HardenedKeyStart + 2,
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
	priv, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	return priv.Serialize(), nil
}

type rawTxVout struct {
	N            uint32  `json:"n"`
	Value        float64 `json:"value"`
	ScriptPubKey struct {
		Hex     string `json:"hex"`
		Address string `json:"address"`
	} `json:"scriptPubKey"`
}

type rawTx struct {
	Vout []rawTxVout `json:"vout"`
	Vin  []struct {
		TxID string `json:"txid"`
		Vout uint32 `json:"vout"`
	} `json:"vin"`
}

type utxoInfo struct {
	txid     string
	vout     uint32
	amount   float64
	pkScript []byte
}

func (l *Litecoin) collectUTXOs(addr string) ([]utxoInfo, error) {
	url := l.rpcURL()
	desc := fmt.Sprintf("addr(%s)", addr)
	if len(addr) > 4 && (addr[:4] == "ltc1" || addr[:4] == "tltc") {
		desc = fmt.Sprintf("wpkh(%s)", addr)
	}

	scan, err := rpc.Call[scanResult](url, "scantxoutset", []any{
		"start", []any{map[string]any{"desc": desc}},
	})
	if err != nil {
		return nil, err
	}

	var utxos []utxoInfo
	for _, u := range scan.Unspents {
		tx, err := rpc.Call[rawTx](url, "getrawtransaction", []any{u.TxID, true})
		if err != nil {
			return nil, fmt.Errorf("getrawtransaction %s: %w", u.TxID, err)
		}
		for _, vout := range tx.Vout {
			if vout.ScriptPubKey.Address == addr {
				pkScriptBytes, err := hex.DecodeString(vout.ScriptPubKey.Hex)
				if err != nil {
					return nil, err
				}
				utxos = append(utxos, utxoInfo{
					txid:     u.TxID,
					vout:     vout.N,
					amount:   vout.Value,
					pkScript: pkScriptBytes,
				})
			}
		}
	}
	return utxos, nil
}

func (l *Litecoin) GetOnChainBalance(addr string) (float64, error) {
	utxos, err := l.collectUTXOs(addr)
	if err != nil {
		return 0, err
	}
	var total float64
	for _, u := range utxos {
		total += u.amount
	}
	return total, nil
}

const (
	p2wpkhInputVbytes  = 68
	p2wpkhOutputVbytes = 31
	txOverheadVbytes   = 10
	defaultFeeRateSat  = 5
)

func (l *Litecoin) SendFunds(privKeyBytes []byte, fromAddr, toAddr string, amount float64) (string, error) {
	url := l.rpcURL()
	params := l.getParams()

	utxos, err := l.collectUTXOs(fromAddr)
	if err != nil {
		return "", fmt.Errorf("collect UTXOs: %w", err)
	}
	if len(utxos) == 0 {
		return "", fmt.Errorf("no UTXOs at %s", fromAddr)
	}

	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

	estimatedVbytes := txOverheadVbytes + p2wpkhInputVbytes*len(utxos) + p2wpkhOutputVbytes
	feeSat := int64(estimatedVbytes * defaultFeeRateSat)

	var totalSat int64
	for _, u := range utxos {
		totalSat += int64(u.amount * 1e8)
	}

	var sendSat int64
	if amount == 0 {
		sendSat = totalSat - feeSat
	} else {
		sendSat = int64(amount * 1e8)
	}
	if sendSat <= 0 || sendSat > totalSat-feeSat {
		return "", fmt.Errorf("insufficient funds: have %d sat, need %d + %d fee", totalSat, sendSat, feeSat)
	}

	destAddr, err := btcutil.DecodeAddress(toAddr, params)
	if err != nil {
		return "", fmt.Errorf("invalid destination address: %w", err)
	}
	destScript, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		return "", err
	}

	tx := wire.NewMsgTx(wire.TxVersion)
	for _, u := range utxos {
		hash, err := chainhash.NewHashFromStr(u.txid)
		if err != nil {
			return "", err
		}
		tx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(hash, u.vout), nil, nil))
	}
	tx.AddTxOut(wire.NewTxOut(sendSat, destScript))

	fetcher := txscript.NewMultiPrevOutFetcher(nil)
	for i, u := range utxos {
		hash, _ := chainhash.NewHashFromStr(u.txid)
		fetcher.AddPrevOut(
			wire.OutPoint{Hash: *hash, Index: u.vout},
			wire.NewTxOut(int64(u.amount*1e8), u.pkScript),
		)
		_ = i
	}
	sigHashes := txscript.NewTxSigHashes(tx, fetcher)

	for i, u := range utxos {
		witness, err := txscript.WitnessSignature(
			tx, sigHashes, i,
			int64(u.amount*1e8),
			u.pkScript,
			txscript.SigHashAll,
			privKey, true,
		)
		if err != nil {
			return "", fmt.Errorf("sign input %d: %w", i, err)
		}
		tx.TxIn[i].Witness = witness
	}

	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return "", err
	}
	txHex := hex.EncodeToString(buf.Bytes())

	txHash, err := rpc.Call[string](url, "sendrawtransaction", []any{txHex})
	if err != nil {
		return "", fmt.Errorf("sendrawtransaction: %w", err)
	}
	return txHash, nil
}

func (l *Litecoin) GetSender(txHash, receiverAddr string) (string, error) {
	url := l.rpcURL()
	params := l.getParams()

	tx, err := rpc.Call[rawTx](url, "getrawtransaction", []any{txHash, true})
	if err != nil {
		return "", err
	}
	if len(tx.Vin) == 0 {
		return "", fmt.Errorf("no inputs in tx %s", txHash)
	}

	prevTx, err := rpc.Call[rawTx](url, "getrawtransaction", []any{tx.Vin[0].TxID, true})
	if err != nil {
		return "", err
	}
	for _, vout := range prevTx.Vout {
		if vout.N == tx.Vin[0].Vout {
			addr := vout.ScriptPubKey.Address
			if addr == "" {
				pkScript, err := hex.DecodeString(vout.ScriptPubKey.Hex)
				if err != nil {
					return "", err
				}
				_, addrs, _, err := txscript.ExtractPkScriptAddrs(pkScript, params)
				if err != nil || len(addrs) == 0 {
					return "", fmt.Errorf("cannot extract sender address")
				}
				addr = addrs[0].EncodeAddress()
			}
			return addr, nil
		}
	}
	return "", fmt.Errorf("sender output not found in tx %s", tx.Vin[0].TxID)
}
