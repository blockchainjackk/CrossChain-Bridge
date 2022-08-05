package dcrn

import (
	"fmt"
	"strings"
	"time"

	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/tokens"
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

// NewCrossChainBridge new dcrn bridge
func NewCrossChainBridge(isSrc bool) *Bridge {
	tokens.IsSwapoutToStringAddress = true
	if !isSrc {
		log.Fatalf("dcrn::NewCrossChainBridge error %v", tokens.ErrBridgeDestinationNotSupported)
	}
	instance := &Bridge{CrossChainBridgeBase: tokens.NewCrossChainBridgeBase(isSrc)}
	//todo
	BridgeInstance = instance
	instance.SetInherit(instance)
	return instance
}

// SetInherit set inherit
func (b *Bridge) SetInherit(inherit Inheritable) {
	b.Inherit = inherit
}

// SetChainAndGateway set chain and gateway config
func (b *Bridge) SetChainAndGateway(chainCfg *tokens.ChainConfig, gatewayCfg *tokens.GatewayConfig) {
	b.CrossChainBridgeBase.SetChainAndGateway(chainCfg, gatewayCfg)
	b.VerifyChainConfig()
	b.InitLatestBlockNumber()
}

// VerifyChainConfig verify chain config
func (b *Bridge) VerifyChainConfig() {
	chainCfg := b.ChainConfig
	networkID := strings.ToLower(chainCfg.NetID)
	switch networkID {
	case netMainnet, netTestnet3:
	case netSimnet:
		return
	default:
		log.Fatal("unsupported dcrn network", "netID", chainCfg.NetID)
	}
}

// VerifyTokenConfig verify token config
func (b *Bridge) VerifyTokenConfig(tokenCfg *tokens.TokenConfig) error {
	if !b.IsP2pkhAddress(tokenCfg.DcrmAddress) {
		return fmt.Errorf("invalid dcrm address (not p2pkh): %v", tokenCfg.DcrmAddress)
	}
	if !b.IsValidAddress(tokenCfg.DepositAddress) {
		return fmt.Errorf("invalid deposit address: %v", tokenCfg.DepositAddress)
	}
	if strings.EqualFold(tokenCfg.Symbol, "DCRN") && *tokenCfg.Decimals != 8 {
		return fmt.Errorf("invalid decimals for DCRN: want 8 but have %v", *tokenCfg.Decimals)
	}
	return nil
}

// InitLatestBlockNumber init latest block number
func (b *Bridge) InitLatestBlockNumber() {
	chainCfg := b.ChainConfig
	gatewayCfg := b.GatewayConfig
	var latest uint64
	var err error
	for {
		latest, err = b.GetLatestBlockNumber()
		if err == nil {
			tokens.SetLatestBlockHeight(latest, b.IsSrc)
			log.Info("get latst block number succeed.", "number", latest, "BlockChain", chainCfg.BlockChain, "NetID", chainCfg.NetID)
			break
		}
		log.Error("get latst block number failed.", "BlockChain", chainCfg.BlockChain, "NetID", chainCfg.NetID, "err", err)
		log.Println("retry query gateway", gatewayCfg.APIAddress)
		time.Sleep(3 * time.Second)
	}
}
