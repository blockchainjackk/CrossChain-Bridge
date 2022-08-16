package dcrn

import "github.com/anyswap/CrossChain-Bridge/tokens"

// SendTransaction send signed tx
func (b *Bridge) SendTransaction(signedTx interface{}) (txHash string, err error) {
	Tx, ok := signedTx.(string)
	if !ok {
		return "", tokens.ErrWrongRawTx
	}
	return b.PostTransaction(Tx)
}
