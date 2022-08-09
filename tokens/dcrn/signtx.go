package dcrn

import (
	"bytes"
	"decred.org/dcrwallet/wallet/txauthor"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"strings"
)

func (b *Bridge) verifyTransactionWithArgs(tx *txauthor.AuthoredTx, args *tokens.BuildTxArgs) error {
	checkReceiver := args.Bind
	if args.Identifier == tokens.AggregateIdentifier {
		checkReceiver = cfgUtxoAggregateToAddress
	}
	payToReceiverScript, _, err := b.GetPayToAddrScript(checkReceiver)
	if err != nil {
		return err
	}
	isRightReceiver := false
	for _, out := range tx.Tx.TxOut {
		if bytes.Equal(out.PkScript, payToReceiverScript) {
			isRightReceiver = true
			break
		}
	}
	if !isRightReceiver {
		return fmt.Errorf("[sign] verify tx receiver failed")
	}
	return nil
}

// DcrmSignTransaction dcrm sign raw tx
func (b *Bridge) DcrmSignTransaction(rawTx interface{}, args *tokens.BuildTxArgs) (signedTx interface{}, txHash string, err error) {
	//todo
	return nil, "", err
}

func checkEqualLength(authoredTx *txauthor.AuthoredTx, msgHash, rsv []string, sigScripts [][]byte) error {
	txIn := authoredTx.Tx.TxIn
	if len(txIn) != len(msgHash) {
		return errors.New("mismatch number of msghashes and tx inputs")
	}
	if len(txIn) != len(rsv) {
		return errors.New("mismatch number of signatures and tx inputs")
	}
	if sigScripts != nil && len(sigScripts) != len(txIn) {
		return errors.New("mismatch number of signatures scripts and tx inputs")
	}
	return nil
}

// VerifyRedeemScript verify redeem script
func (b *Bridge) VerifyRedeemScript(prevScript, redeemScript []byte) error {
	p2shScript, _, err := b.GetP2shSigScript(redeemScript)
	if err != nil {
		return err
	}
	if !bytes.Equal(p2shScript, prevScript) {
		return fmt.Errorf("redeem script %x mismatch", redeemScript)
	}
	return nil
}

func (b *Bridge) verifyPublickeyData(pkData []byte) error {
	tokenCfg := b.GetTokenConfig(PairID)
	if tokenCfg == nil {
		return nil
	}
	dcrmAddress := tokenCfg.DcrmAddress
	if dcrmAddress == "" {
		return nil
	}
	address, err := b.NewAddressPubKeyHash(pkData)
	if err != nil {
		return err
	}
	if address.Address() != dcrmAddress {
		return fmt.Errorf("public key address %v is not the configed dcrm address %v", address, dcrmAddress)
	}
	return nil
}

// SignTransaction sign tx with pairID
func (b *Bridge) SignTransaction(rawTx interface{}, pairID string) (signedTx interface{}, txHash string, err error) {
	authoredTx, ok := rawTx.(*txauthor.AuthoredTx)

	if !ok {
		return nil, "", tokens.ErrWrongRawTx
	}
	var builder strings.Builder
	builder.Grow(2 * authoredTx.Tx.SerializeSize())
	err = authoredTx.Tx.Serialize(hex.NewEncoder(&builder))
	if err != nil {
		return nil, "", err
	}

	return b.SignRawtransaction(builder.String())
}

//todo 这个可能用不到了
// GetCompressedPublicKey get compressed public key
func (b *Bridge) GetCompressedPublicKey(fromPublicKey string, needVerify bool) (cPkData []byte, err error) {

	return nil, nil
}
