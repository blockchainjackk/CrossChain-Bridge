package dcrn

import (
	"fmt"
	"strings"

	"decred.org/dcrwallet/wallet"
	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/chaincfg/v3"
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/decred/dcrd/txscript/v3"
	"github.com/decred/dcrd/wire"
)

// Inheritable interface
type Inheritable interface {
	GetChainParams() *chaincfg.Params
}

type dcrnAmountType = dcrutil.Amount
type wireTxInType = wire.TxIn
type wireTxOutType = wire.TxOut

func isValidValue(value dcrnAmountType) bool {
	return value > 0 && value <= dcrutil.MaxAmount
}

func newAmount(value float64) (dcrnAmountType, error) {
	amount, err := dcrutil.NewAmount(value)
	return amount, err
}

// GetChainParams get chain config (net params)
func (b *Bridge) GetChainParams() *chaincfg.Params {
	networkID := strings.ToLower(b.ChainConfig.NetID)
	switch networkID {
	case netMainnet:
		return chaincfg.MainNetParams()
	case netSimnet:
		return chaincfg.SimNetParams()
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
func (b *Bridge) GetPayToAddrScript(address string) (pkScript []byte, version uint16, err error) {
	toAddr, err := b.DecodeAddress(address)
	if err != nil {
		return nil, 0, fmt.Errorf("decode dcrn address '%v' failed. %w", address, err)
	}
	//todo 调研这个version
	switch addr := toAddr.(type) {
	case wallet.V0Scripter:
		return addr.ScriptV0(), 0, nil
	default:
		pkScript, err = txscript.PayToAddrScript(addr)
		return pkScript, 0, err
	}
}

// NullDataScript encap
func (b *Bridge) NullDataScript(memo string) ([]byte, error) {
	bmemo := []byte(memo)
	if len(bmemo) > txscript.MaxDataCarrierSize {
		str := fmt.Sprintf("data size %d is larger than max "+
			"allowed size %d", len(bmemo), txscript.MaxDataCarrierSize)
		return nil, fmt.Errorf(str)
	}

	script, err := txscript.NewScriptBuilder().AddOp(txscript.OP_RETURN).AddData(bmemo).Script()
	return script, err

}

// IsPayToScriptHash is p2sh
func (b *Bridge) IsPayToScriptHash(sigScript []byte) bool {
	return txscript.IsPayToScriptHash(sigScript)
}

// CalcSignatureHash calc sig hash
func (b *Bridge) CalcSignatureHash(sigScript []byte, tx *wire.MsgTx, i int) (sigHash []byte, err error) {
	//todo  有问题,dcrn多一个cachedPrefix 参数
	return txscript.CalcSignatureHash(sigScript, txscript.SigHashAll, tx, i, nil)
}

// GetP2shRedeemScript get p2sh redeem script
func (b *Bridge) GetP2shRedeemScript(memo, pubKeyHash []byte) (redeemScript []byte, err error) {
	return txscript.NewScriptBuilder().
		AddData(memo).AddOp(txscript.OP_DROP).
		AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).AddData(pubKeyHash).
		AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).
		Script()
}

// NewTxIn new txin
func (b *Bridge) NewTxIn(txid string, vout uint32, value int64, pkScript []byte) (*wire.TxIn, error) {
	txHash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		return nil, err
	}
	//todo dcrn多一个tree参数
	prevOutPoint := wire.NewOutPoint(txHash, vout, 0)
	//todo 参数不一样 多一个value
	return wire.NewTxIn(prevOutPoint, value, pkScript), nil
}

// NewTxOut new txout
func (b *Bridge) NewTxOut(amount int64, pkScript []byte) *wire.TxOut {
	return wire.NewTxOut(amount, pkScript)
}
