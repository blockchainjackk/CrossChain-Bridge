package dcrn

import (
	"github.com/anyswap/CrossChain-Bridge/tokens/dcrn/ctl"
	"github.com/decred/dcrwallet/rpc/jsonrpc/types"
)

// GetLatestBlockNumberOf impl
func (b *Bridge) GetLatestBlockNumberOf(apiAddress string) (uint64, error) {
	return 0, nil
}

// GetLatestBlockNumber impl
func (b *Bridge) GetLatestBlockNumber() (uint64, error) {

	return 0, nil
}

func (b *Bridge) GetTransactionByHash(txHash string) (*types.GetTransactionResult, error) {
	return ctl.GetTransactionByHash(b, txHash)
}

// PostTransaction impl
func (b *Bridge) PostTransaction(txHex string) (txHash string, err error) {
	return ctl.PostTransaction(b, txHex)
}
