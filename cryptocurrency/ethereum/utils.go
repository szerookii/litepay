package ethereum

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
)

func Client() (*ethclient.Client, error) {
	client, err := ethclient.Dial(os.Getenv("ETH_RPC_HOST"))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func WeiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
}

func GweiToEther(gwei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(gwei), big.NewFloat(1e9))
}
