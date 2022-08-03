package dcrn

import (
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/decred/dcrd/chaincfg/v3"

	"strings"
)

const (
	netMainnet  = "mainnet"
	netTestnet3 = "testnet3"
	netSimnet   = "simnet"
)

// PairID unique dcrn pair ID
var PairID = "dcrn"

// Bridge btc bridge
type Bridge struct {
	*tokens.CrossChainBridgeBase
	Inherit Inheritable
}

// GetChainParams get chain config (net params)
func (b *Bridge) GetChainParams() *chaincfg.Params {
	networkID := strings.ToLower(b.ChainConfig.NetID)
	switch networkID {
	case netMainnet:
		return chaincfg.MainNetParams()
	default:
		return chaincfg.TestNet3Params()
	}
}
