package litecoin

import (
	"github.com/btcsuite/btcd/chaincfg"
)

// MainNetParams defines the network parameters for the main Litecoin network.
var MainNetParams = chaincfg.Params{
	Name: "litecoin",
	Net:  0xdbb6c0fb, // Litecoin mainnet magic

	// Address prefixes
	PubKeyHashAddrID: 0x30, // L
	ScriptHashAddrID: 0x32, // M
	PrivateKeyID:     0xb0,

	// BIP32 prefixes (same as Bitcoin)
	HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e},
	HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4},

	// Bech32 prefix
	Bech32HRPSegwit: "ltc",
}

// TestNet4Params defines the network parameters for the Litecoin test network.
var TestNet4Params = chaincfg.Params{
	Name: "litecoin-testnet",
	Net:  0xfdd1c0db, // Litecoin testnet magic

	// Address prefixes
	PubKeyHashAddrID: 0x6f, // n or m
	ScriptHashAddrID: 0x3a, // Q (or 2)
	PrivateKeyID:     0xef,

	// BIP32 prefixes (tpub/tprv)
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf},
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94},

	// Bech32 prefix
	Bech32HRPSegwit: "tltc",
}
