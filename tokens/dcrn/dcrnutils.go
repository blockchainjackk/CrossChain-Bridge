package dcrn

import (
	"crypto/ecdsa"
	"decred.org/dcrwallet/wallet"
	"fmt"
	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/chaincfg/v3"
	"github.com/decred/dcrd/dcrec"
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/decred/dcrd/txscript/v3"
	"github.com/decred/dcrd/wire"
	"math/big"
	"strings"
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
func (b *Bridge) GetPayToAddrScript(address string) (pkScript []byte, version uint16, err error) {
	//todo 在DecodeAddress中有问题
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
	//todo dcrn 中没有这个函数了
	txscript.NullDataScript([]byte(memo))

	return nil, nil

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
func (b *Bridge) NewTxIn(txid string, vout uint32, pkScript []byte) (*wire.TxIn, error) {
	txHash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		return nil, err
	}
	//todo dcrn多一个tree参数
	prevOutPoint := wire.NewOutPoint(txHash, vout, 0)
	//todo 参数不一样 多一个value
	return wire.NewTxIn(prevOutPoint, pkScript), nil
}

// NewTxOut new txout
func (b *Bridge) NewTxOut(amount int64, pkScript []byte) *wire.TxOut {
	return wire.NewTxOut(amount, pkScript)
}

//todo
// GetSigScript get script
func (b *Bridge) GetSigScript(sigScripts [][]byte, prevScript, signData, cPkData []byte, i int) (sigScript []byte, err error) {
	scriptClass := txscript.GetScriptClass(prevScript)
	switch scriptClass {
	case txscript.PubKeyHashTy:
		sigScript, err = txscript.NewScriptBuilder().AddData(signData).AddData(cPkData).Script()
	case txscript.ScriptHashTy:
		if sigScripts == nil {
			err = fmt.Errorf("call MakeSignedTransaction spend p2sh without redeem scripts")
		} else {
			redeemScript := sigScripts[i]
			err = b.VerifyRedeemScript(prevScript, redeemScript)
			if err == nil {
				sigScript, err = txscript.NewScriptBuilder().AddData(signData).AddData(cPkData).AddData(redeemScript).Script()
			}
		}
	default:
		err = fmt.Errorf("unsupport to spend '%v' output", scriptClass.String())
	}
	return sigScript, err
}

// SerializeSignature serialize signature
func (b *Bridge) SerializeSignature(r, s *big.Int) []byte {
	sign := &dcrec.Signature{R: r, S: s}
	return append(sign.Serialize(), byte(txscript.SigHashAll))
	return nil
}

// SignWithECDSA sign with ecdsa private key
func (b *Bridge) SignWithECDSA(privKey *ecdsa.PrivateKey, msgHash []byte) (rsv string, err error) {
	signature, err := (*dcrec.PrivateKey)(privKey).Sign(msgHash)
	if err != nil {
		return "", err
	}
	rr := fmt.Sprintf("%064X", signature.R)
	ss := fmt.Sprintf("%064X", signature.S)
	rsv = fmt.Sprintf("%s%s00", rr, ss)
	return rsv, nil
}

// GetPublicKeyFromECDSA get public key from ecdsa private key
func (b *Bridge) GetPublicKeyFromECDSA(privKey *ecdsa.PrivateKey, compressed bool) []byte {
	if compressed {
		return (*dcrec.PublicKey)(&privKey.PublicKey).SerializeCompressed()
	}
	return (*dcrec.PublicKey)(&privKey.PublicKey).SerializeUncompressed()
}
