package dcrn

import (
	"fmt"
	"github.com/decred/dcrd/chaincfg/v3"
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/decred/dcrd/txscript/v3"
	"strings"
)

// Inheritable interface
type Inheritable interface {
	GetChainParams() *chaincfg.Params
}

type dcrnAmountType dcrutil.Amount

func isValidValue(value dcrnAmountType) bool {
	return value > 0 && value <= dcrutil.MaxAmount
}

func newAmount(value float64) (dcrnAmountType, error) {
	amount, err := dcrutil.NewAmount(value)
	return dcrnAmountType(amount), err
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

//todo  dcrn 中没有pkScript
//func (b *Bridge) ParsePkScript(pkScript []byte) (txscript.Pk, error) {
//	//return txscript.
//	return
//}

// GetPayToAddrScript get pay to address script
func (b *Bridge) GetPayToAddrScript(address string) ([]byte, error) {
	toAddr, err := b.DecodeAddress(address)
	if err != nil {
		return nil, fmt.Errorf("decode dcrn address '%v' failed. %w", address, err)
	}
	return txscript.PayToAddrScript(toAddr)
}

// GetP2shRedeemScript get p2sh redeem script
func (b *Bridge) GetP2shRedeemScript(memo, pubKeyHash []byte) (redeemScript []byte, err error) {
	return txscript.NewScriptBuilder().
		AddData(memo).AddOp(txscript.OP_DROP).
		AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).AddData(pubKeyHash).
		AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).
		Script()
}
