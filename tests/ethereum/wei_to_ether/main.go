package main

import (
	"fmt"
	"github.com/szerookii/litepay/cryptocurrency/ethereum"
	"math/big"
)

func main() {
	wei := int64(10000000000000000)
	fmt.Printf("wei to ether: %f\n", ethereum.WeiToEther(big.NewInt(wei)))

	gwei := int64(10000000)
	fmt.Printf("gwei to ether: %f\n", ethereum.GweiToEther(big.NewInt(gwei)))

}
