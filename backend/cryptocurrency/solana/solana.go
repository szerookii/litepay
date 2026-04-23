package solana

import (
	"bytes"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mr-tron/base58/base58"
	"github.com/phuslu/log"
	"github.com/szerookii/litepay/backend/cryptocurrency"
	"github.com/szerookii/litepay/backend/utils/coingecko"
)

const defaultRPCURL = "https://api.mainnet-beta.solana.com"

type Solana struct{}

func (s *Solana) Name() string               { return "Solana" }
func (s *Solana) Symbol() string             { return "SOL" }
func (s *Solana) CoinGeckoID() string        { return "solana" }
func (s *Solana) RequiredConfirmations() int { return 1 }

func (s *Solana) rpcURL() string {
	if u := os.Getenv("SOL_RPC_URL"); u != "" {
		return u
	}
	return defaultRPCURL
}

// slip10DeriveED25519 implements SLIP-0010 private key derivation for ED25519.
// All indices MUST be hardened (caller must OR with 0x80000000).
func slip10DeriveED25519(seed []byte, path []uint32) ([]byte, error) {
	mac := hmac.New(sha512.New, []byte("ed25519 seed"))
	mac.Write(seed)
	I := mac.Sum(nil)
	kL, kR := I[:32], I[32:]

	for _, idx := range path {
		mac = hmac.New(sha512.New, kR)
		mac.Write([]byte{0x00})
		mac.Write(kL)
		idxBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(idxBytes, idx)
		mac.Write(idxBytes)
		I = mac.Sum(nil)
		kL, kR = I[:32], I[32:]
	}
	return kL, nil
}

const hardenedOffset = uint32(0x80000000)

// NewPaymentAddress derives a unique SOL address at m/44'/501'/accountIndex'/0'/paymentIndex'.
func (s *Solana) NewPaymentAddress(masterSeed []byte, accountIndex uint32, paymentIndex uint32) (*cryptocurrency.PaymentAddress, error) {
	path := []uint32{
		hardenedOffset + 44,
		hardenedOffset + 501,
		hardenedOffset + accountIndex,
		hardenedOffset + 0,
		hardenedOffset + paymentIndex,
	}
	privKey, err := slip10DeriveED25519(masterSeed, path)
	if err != nil {
		return nil, fmt.Errorf("SOL derive: %w", err)
	}

	pub := ed25519.NewKeyFromSeed(privKey).Public().(ed25519.PublicKey)
	addr := base58.Encode(pub)

	return &cryptocurrency.PaymentAddress{
		Address: addr,
		Index:   paymentIndex,
	}, nil
}

// ── Solana JSON-RPC helpers ───────────────────────────────────────────────────

type rpcReq struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func (s *Solana) call(method string, params []interface{}, result interface{}) error {
	body, _ := json.Marshal(rpcReq{JSONRPC: "2.0", ID: 1, Method: method, Params: params})
	resp, err := http.Post(s.rpcURL(), "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(result)
}

type sigsResult struct {
	Result []struct {
		Signature          string `json:"signature"`
		ConfirmationStatus string `json:"confirmationStatus"`
	} `json:"result"`
}

type txResult struct {
	Result *struct {
		Meta *struct {
			Err          interface{} `json:"err"`
			PreBalances  []int64     `json:"preBalances"`
			PostBalances []int64     `json:"postBalances"`
		} `json:"meta"`
		Transaction struct {
			Message struct {
				AccountKeys []string `json:"accountKeys"`
			} `json:"message"`
		} `json:"transaction"`
	} `json:"result"`
}

func (s *Solana) CheckPayment(addr *cryptocurrency.PaymentAddress, requiredSOL float64) (*cryptocurrency.PaymentStatus, error) {
	var sigs sigsResult
	if err := s.call("getSignaturesForAddress", []interface{}{
		addr.Address,
		map[string]interface{}{"limit": 10},
	}, &sigs); err != nil {
		return nil, err
	}
	return s.verifySignatures(sigs, addr.Address, requiredSOL)
}

func (s *Solana) verifySignatures(sigs sigsResult, targetAddr string, requiredSOL float64) (*cryptocurrency.PaymentStatus, error) {
	requiredLamports := int64(requiredSOL * 1e9)
	log.Debug().Str("addr", targetAddr).Int("sig_count", len(sigs.Result)).Msg("Cron: Solana verifying signatures")

	for _, sig := range sigs.Result {
		var txr txResult
		if err := s.call("getTransaction", []interface{}{
			sig.Signature,
			map[string]interface{}{"encoding": "json", "commitment": "confirmed"},
		}, &txr); err != nil || txr.Result == nil || txr.Result.Meta == nil {
			log.Debug().Str("sig", sig.Signature).Msg("Cron: Solana skipped tx (not found or error)")
			continue
		}

		if txr.Result.Meta.Err != nil {
			log.Debug().Str("sig", sig.Signature).Interface("err", txr.Result.Meta.Err).Msg("Cron: Solana tx has error")
			continue
		}

		keys := txr.Result.Transaction.Message.AccountKeys
		idx := -1
		for i, k := range keys {
			if k == targetAddr {
				idx = i
				break
			}
		}

		if idx < 0 {
			continue
		}

		received := txr.Result.Meta.PostBalances[idx] - txr.Result.Meta.PreBalances[idx]
		log.Debug().
			Str("sig", sig.Signature).
			Int64("received_lamports", received).
			Int64("required_lamports", requiredLamports).
			Msg("Cron: Solana balance change calculated")

		if received >= requiredLamports-100 {
			confs := 0
			if sig.ConfirmationStatus == "confirmed" {
				confs = 1
			} else if sig.ConfirmationStatus == "finalized" {
				confs = 32
			}

			return &cryptocurrency.PaymentStatus{
				ReceivedAmount: float64(received) / 1e9,
				Confirmations:  confs,
				TxHash:         sig.Signature,
				IsPending:      confs == 0,
			}, nil
		}
	}

	return &cryptocurrency.PaymentStatus{}, nil
}

func (s *Solana) Price(currency string) (float64, error) {
	prices, err := coingecko.GetPrice([]string{s.CoinGeckoID()}, strings.ToLower(currency))
	if err != nil {
		return 0, err
	}
	p, ok := prices[s.CoinGeckoID()][strings.ToLower(currency)]
	if !ok {
		return 0, fmt.Errorf("no price for SOL/%s", currency)
	}
	return p, nil
}
