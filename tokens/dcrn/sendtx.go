package dcrn

import (
	"bytes"
	"decred.org/dcrwallet/wallet/txauthor"
	"encoding/hex"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/tokens"
)

// SendTransaction send signed tx
func (b *Bridge) SendTransaction(signedTx interface{}) (txHash string, err error) {
	//todo  txauthor.AuthoredTx 里面的参数和BTC不一样 没有PrevInputValues，应该会影响到buildTx
	authoredTx, ok := signedTx.(*txauthor.AuthoredTx)
	if !ok {
		return "", tokens.ErrWrongRawTx
	}

	tx := authoredTx.Tx
	if tx == nil {
		return "", tokens.ErrWrongRawTx
	}

	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	err = tx.Serialize(buf)
	if err != nil {
		return "", err
	}
	txHex := hex.EncodeToString(buf.Bytes())
	log.Info("Bridge send tx", "hash", tx.TxHash())

	return b.PostTransaction(txHex)
}
