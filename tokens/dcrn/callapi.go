package dcrn

import (
	"fmt"
	"math/big"

	"github.com/anyswap/CrossChain-Bridge/tokens/btc/electrs"
	"github.com/anyswap/CrossChain-Bridge/tokens/dcrn/ctl"
)

// GetLatestBlockNumberOf impl
func (b *Bridge) GetLatestBlockNumberOf(apiAddress string) (uint64, error) {
	return ctl.GetLatestBlockNumberOf(apiAddress)
}

// GetLatestBlockNumber impl
func (b *Bridge) GetLatestBlockNumber() (uint64, error) {
	return ctl.GetLatestBlockNumber(b)
}

func (b *Bridge) GetTransactionByHash(txHash string) (*electrs.ElectTx, error) {
	return ctl.GetTransactionByHash(b, txHash)
}

// GetElectTransactionStatus impl
func (b *Bridge) GetElectTransactionStatus(txHash string) (*electrs.ElectTxStatus, error) {
	return ctl.GetTransactionStatus(b, txHash)
}

// FindUtxos impl
func (b *Bridge) FindUtxos(addr string) ([]*electrs.ElectUtxo, error) {
	addrs := []string{addr}
	return ctl.FindUtxos(b, addrs)
}

// GetPoolTxidList impl
func (b *Bridge) GetPoolTxidList() ([]string, error) {
	return ctl.GetPoolTxidList(b)
}

// GetPoolTransactions impl
func (b *Bridge) GetPoolTransactions(addr string) ([]*electrs.ElectTx, error) {
	return ctl.GetPoolTransactions(b, addr)
}

// GetTransactionHistory impl
func (b *Bridge) GetTransactionHistory(addr, lastSeenTxid string) ([]*electrs.ElectTx, error) {
	return ctl.GetTransactionHistory(b, addr, lastSeenTxid)
}

// GetOutspend impl
// func (b *Bridge) GetOutspend(txHash string, vout uint32) (*electrs.ElectOutspend, error) {
// 	return ctl.GetOutspend(b, txHash, vout)
// }

// PostTransaction impl
func (b *Bridge) PostTransaction(txHex string) (txHash string, err error) {
	return ctl.PostTransaction(b, txHex)
}

// GetBlockHash impl
func (b *Bridge) GetBlockHash(height uint64) (string, error) {
	return ctl.GetBlockHash(b, height)
}

// GetBlockTxids impl
func (b *Bridge) GetBlockTxids(blockHash string) ([]string, error) {
	return ctl.GetBlockTxids(b, blockHash)
}

// GetBlock impl
func (b *Bridge) GetBlock(blockHash string) (*electrs.ElectBlock, error) {
	return ctl.GetBlock(b, blockHash)
}

// GetBlockTransactions impl
func (b *Bridge) GetBlockTransactions(blockHash string, startIndex uint32) ([]*electrs.ElectTx, error) {
	return ctl.GetBlockTransactions(b, blockHash, startIndex)
}

// EstimateFeePerKb impl
func (b *Bridge) EstimateFeePerKb(blocks int) (int64, error) {
	return ctl.EstimateFeePerKb(b, blocks)
}

// GetBalance impl
func (b *Bridge) GetBalance(account string) (*big.Int, error) {
	utxos, err := b.FindUtxos(account)
	if err != nil {
		return nil, err
	}
	var balance uint64
	for _, utxo := range utxos {
		balance += *utxo.Value
	}
	return new(big.Int).SetUint64(balance), nil
}

// GetTokenBalance impl
func (b *Bridge) GetTokenBalance(tokenType, tokenAddress, accountAddress string) (*big.Int, error) {
	return nil, fmt.Errorf("[%v] can not get token balance of token with type '%v'", b.ChainConfig.BlockChain, tokenType)
}

// GetTokenSupply impl
func (b *Bridge) GetTokenSupply(tokenType, tokenAddress string) (*big.Int, error) {
	return nil, fmt.Errorf("[%v] can not get token supply of token with type '%v'", b.ChainConfig.BlockChain, tokenType)
}

//
func (b *Bridge) SignRawtransaction(unsignedHex string) (signedHex string, err error) {
	return ctl.SignRawtransaction(b, unsignedHex)
}
