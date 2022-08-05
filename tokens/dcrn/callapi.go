package dcrn

import (
	"github.com/anyswap/CrossChain-Bridge/tokens/btc/electrs"
	"github.com/anyswap/CrossChain-Bridge/tokens/dcrn/ctl"
)

// GetLatestBlockNumberOf impl
func (b *Bridge) GetLatestBlockNumberOf(apiAddress string) (uint64, error) {
	return 0, nil
}

// GetLatestBlockNumber impl
func (b *Bridge) GetLatestBlockNumber() (uint64, error) {

	return 0, nil
}

func (b *Bridge) GetTransactionByHash(txHash string) (*electrs.ElectTx, error) {
	return nil, nil
}

// PostTransaction impl
func (b *Bridge) PostTransaction(txHex string) (txHash string, err error) {
	return ctl.PostTransaction(b, txHex)
}

// GetElectTransactionStatus impl
func (b *Bridge) GetElectTransactionStatus(txHash string) (*electrs.ElectTxStatus, error) {

	return nil, nil
}

// EstimateFeePerKb impl
func (b *Bridge) EstimateFeePerKb(blocks int) (int64, error) {
	return 0, nil
}

// FindUtxos impl
func (b *Bridge) FindUtxos(addr string) ([]*electrs.ElectUtxo, error) {
	return electrs.FindUtxos(b, addr)
}

// GetOutspend impl
func (b *Bridge) GetOutspend(txHash string, vout uint32) (*electrs.ElectOutspend, error) {
	return electrs.GetOutspend(b, txHash, vout)
}
