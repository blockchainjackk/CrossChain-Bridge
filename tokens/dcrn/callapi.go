package dcrn

import (
	"github.com/anyswap/CrossChain-Bridge/tokens/dcrn/ctl"
	"github.com/decred/dcrwallet/rpc/jsonrpc/types"
)

func (b *Bridge) GetTransactionByHash(txHash string) (*types.GetTransactionResult, error) {
	return ctl.GetTransactionByHash(b, txHash)
}

// PostTransaction impl
func (b *Bridge) PostTransaction(txHex string) (txHash string, err error) {
	return ctl.PostTransaction(b, txHex)
}
