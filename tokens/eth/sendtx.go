package eth

import (
	"errors"
	"strings"

	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/params"
	"github.com/anyswap/CrossChain-Bridge/types"
)

// SendTransaction send signed tx
func (b *Bridge) SendTransaction(signedTx interface{}) (txHash string, err error) {
	tx, ok := signedTx.(*types.Transaction)
	if !ok {
		log.Printf("signed tx is %+v", signedTx)
		return "", errors.New("wrong signed transaction type")
	}
	txHash, err = b.SendSignedTransaction(tx)
	if err != nil {
		if !strings.Contains(err.Error(), "json-rpc error -32000, already known") {
			log.Error("SendTransaction failed", "hash", txHash, "err", err)
		}
	} else {
		log.Info("SendTransaction success", "hash", txHash)
	}
	if params.IsDebugMode() {
		log.Infof("SendTransaction rawtx is %v", tx.RawStr())
	}
	return txHash, err
}
