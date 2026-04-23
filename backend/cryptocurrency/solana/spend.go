package solana

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"net/http"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
)

func (s *Solana) rpcClient() *solanarpc.Client {
	return solanarpc.New(s.rpcURL())
}

func (s *Solana) DerivePrivKey(masterSeed []byte, accountIndex, paymentIndex uint32) ([]byte, error) {
	path := []uint32{
		hardenedOffset + 44,
		hardenedOffset + 501,
		hardenedOffset + accountIndex,
		hardenedOffset + 0,
		hardenedOffset + paymentIndex,
	}
	return slip10DeriveED25519(masterSeed, path)
}

func (s *Solana) GetOnChainBalance(addr string) (float64, error) {
	lamports, err := s.getLamports(addr)
	if err != nil {
		return 0, err
	}
	return float64(lamports) / 1e9, nil
}

func (s *Solana) getLamports(addr string) (uint64, error) {
	pubkey, err := solanago.PublicKeyFromBase58(addr)
	if err != nil {
		return 0, fmt.Errorf("invalid SOL address: %w", err)
	}
	result, err := s.rpcClient().GetBalance(
		context.Background(),
		pubkey,
		solanarpc.CommitmentConfirmed,
	)
	if err != nil {
		return 0, err
	}
	return result.Value, nil
}

func (s *Solana) SendFunds(privKeyBytes []byte, fromAddr, toAddr string, amount float64) (string, error) {
	// Build full 64-byte ed25519 private key from 32-byte SLIP-0010 seed
	fullKey := ed25519.NewKeyFromSeed(privKeyBytes)
	account := solanago.PrivateKey(fullKey)

	toPubkey, err := solanago.PublicKeyFromBase58(toAddr)
	if err != nil {
		return "", fmt.Errorf("invalid destination: %w", err)
	}

	client := s.rpcClient()

	const feeLamports = uint64(5000)

	balanceLamports, err := s.getLamports(fromAddr)
	if err != nil {
		return "", fmt.Errorf("get balance: %w", err)
	}

	var lamports uint64
	if amount == 0 {
		// Sweep: send everything minus fee. Result is exactly 0 lamports → account closes cleanly.
		if balanceLamports <= feeLamports {
			return "", fmt.Errorf("insufficient SOL: have %d lamports, need > %d for fee", balanceLamports, feeLamports)
		}
		lamports = balanceLamports - feeLamports
	} else {
		lamports = uint64(amount * 1e9)
		if lamports+feeLamports > balanceLamports {
			return "", fmt.Errorf("insufficient SOL: have %d lamports, need %d + %d fee", balanceLamports, lamports, feeLamports)
		}
	}

	recent, err := client.GetLatestBlockhash(context.Background(), solanarpc.CommitmentConfirmed)
	if err != nil {
		return "", fmt.Errorf("get blockhash: %w", err)
	}

	tx, err := solanago.NewTransaction(
		[]solanago.Instruction{
			system.NewTransferInstruction(lamports, account.PublicKey(), toPubkey).Build(),
		},
		recent.Value.Blockhash,
		solanago.TransactionPayer(account.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("build tx: %w", err)
	}

	_, err = tx.Sign(func(key solanago.PublicKey) *solanago.PrivateKey {
		if key.Equals(account.PublicKey()) {
			return &account
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("sign tx: %w", err)
	}

	opts := solanarpc.TransactionOpts{
		SkipPreflight:       false,
		PreflightCommitment: solanarpc.CommitmentProcessed,
		MaxRetries:          &[]uint{3}[0],
	}
	sig, err := client.SendTransactionWithOpts(context.Background(), tx, opts)
	if err != nil {
		return "", fmt.Errorf("send tx: %w", err)
	}
	return sig.String(), nil
}

type solTxResult struct {
	Result *struct {
		Meta *struct {
			PreBalances  []int64 `json:"preBalances"`
			PostBalances []int64 `json:"postBalances"`
		} `json:"meta"`
		Transaction struct {
			Message struct {
				AccountKeys []string `json:"accountKeys"`
			} `json:"message"`
		} `json:"transaction"`
	} `json:"result"`
}

func (s *Solana) GetSender(txHash, receiverAddr string) (string, error) {
	body, _ := json.Marshal(rpcReq{
		JSONRPC: "2.0", ID: 1, Method: "getTransaction",
		Params: []interface{}{
			txHash,
			map[string]interface{}{"encoding": "json", "commitment": "confirmed"},
		},
	})
	resp, err := http.Post(s.rpcURL(), "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result solTxResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Result == nil || result.Result.Meta == nil {
		return "", fmt.Errorf("transaction %s not found", txHash)
	}

	keys := result.Result.Transaction.Message.AccountKeys
	pre := result.Result.Meta.PreBalances
	post := result.Result.Meta.PostBalances

	for i, key := range keys {
		if key == receiverAddr {
			continue
		}
		if i < len(pre) && i < len(post) && post[i] < pre[i] {
			return key, nil
		}
	}
	return "", fmt.Errorf("could not determine sender for tx %s", txHash)
}
