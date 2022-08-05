package ctl

import (
	"fmt"
	"sort"

	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/btc/electrs"
	dtypes "github.com/decred/dcrd/rpc/jsonrpc/types/v2"
	"github.com/decred/dcrwallet/rpc/jsonrpc/types"
)

// GetLatestBlockNumberOf
//查询最新块高
func GetLatestBlockNumberOf(apiAddress string) (uint64, error) {
	var result uint64
	err := CallGet(&result, apiAddress, "getblockcount")
	if err == nil {
		return result, nil
	}
	return 0, err
}

// GetLatestBlockNumber
//查询最新块高
func GetLatestBlockNumber(b tokens.CrossChainBridge) (result uint64, err error) {
	gateway := b.GetGatewayConfig()
	for _, apiAddress := range gateway.APIAddress {
		err := CallGet(&result, apiAddress, "getblockcount")
		if err == nil {
			return result, nil
		}
	}
	return 0, err
}

// GetTransactionByHash
//获取交易信息
func GetTransactionByHash(b tokens.CrossChainBridge, txHash string) (*electrs.ElectTx, error) {
	gateway := b.GetGatewayConfig()
	var txRawResult dtypes.TxRawResult
	var err error
	for _, apiAddress := range gateway.APIAddress {
		//注意：要用getrawtransaction（节点启动时需要增加--txindex参数），不要用gettransaction
		err := CallGet(&txRawResult, apiAddress, "getrawtransaction", txHash, 1)
		if err == nil {
			result := TxRawResult2ElectTx(&txRawResult)
			return result, nil
		}
	}
	return nil, err
}

//GetTransactionStatus
//获取交易状态
func GetTransactionStatus(b tokens.CrossChainBridge, txHash string) (*electrs.ElectTxStatus, error) {
	transaction, err := GetTransactionByHash(b, txHash)
	if err != nil {
		return nil, err
	}
	status := transaction.Status
	return status, nil
}

// FindUtxos
//查找utxo
func FindUtxos(b tokens.CrossChainBridge, addrs []string) ([]*electrs.ElectUtxo, error) {
	gateway := b.GetGatewayConfig()
	minconf := 1
	maxconf := 999
	var unspentResult []*types.ListUnspentResult
	var err error
	for _, apiAddress := range gateway.APIAddress {
		err := CallGet(&unspentResult, apiAddress, "listunspent", minconf, maxconf, addrs)
		if err == nil {
			result := SliceUnspentResult2ElectUtxo(unspentResult)
			// sort.Sort(SortableUnspentSlice(unspentResult))
			sort.Sort(electrs.SortableElectUtxoSlice(result))
			return result, nil
		}
	}
	return nil, err
}

// GetPoolTxidList
func GetPoolTxidList(b tokens.CrossChainBridge) (result []string, err error) {
	gateway := b.GetGatewayConfig()
	for _, apiAddress := range gateway.APIAddress {
		err := CallGet(&result, apiAddress, "getrawmempool")
		if err == nil {
			return result, nil
		}
	}
	return nil, err
}

// GetPoolTransactions
func GetPoolTransactions(b tokens.CrossChainBridge, addr string) (result []*electrs.ElectTx, err error) {
	//btc的该方法的上层函数没有被调用，应该不用实现
	return
}

// GetTransactionHistory
func GetTransactionHistory(b tokens.CrossChainBridge, addr, lastSeenTxid string) (result []*electrs.ElectTx, err error) {
	//根据地址获取交易历史，应该要通过浏览器，尤其是有lastSeenTxid参数
	return
}

// GetOutspend
func GetOutspend(b tokens.CrossChainBridge, txHash string, vout uint32) (result *electrs.ElectOutspend, err error) {
	gateway := b.GetGatewayConfig()
	var txOutResult dtypes.GetTxOutResult
	for _, apiAddress := range gateway.APIAddress {
		err := CallGet(&txOutResult, apiAddress, "gettxout", txHash, vout, true)
		if err == nil {
			result = TxOutResult2ElectOutspend(&txOutResult)
			return result, nil
		}
	}
	return
}

// PostTransaction
//发送交易
func PostTransaction(b tokens.CrossChainBridge, txHex string) (txHash string, err error) {
	gateway := b.GetGatewayConfig()
	var success bool
	for _, apiAddress := range gateway.APIAddress {
		hash0, err0 := CallPost(apiAddress, "sendrawtransaction", txHex)
		if err0 == nil && !success {
			success = true
			txHash = hash0
		} else if err0 != nil {
			err = err0
		}
	}
	return txHash, err
}

// GetBlockHash
// 根据块高查询区块的hash
func GetBlockHash(b tokens.CrossChainBridge, height uint64) (blockHash string, err error) {
	gateway := b.GetGatewayConfig()
	for _, apiAddress := range gateway.APIAddress {
		err := CallGet(&blockHash, apiAddress, "getblockhash", height)
		if err == nil {
			return blockHash, nil
		}
	}
	return "", err
}

// GetBlockTxids
// 根据区块hash查询该区块中包含的（普通交易）Txids
func GetBlockTxids(b tokens.CrossChainBridge, blockHash string) (result []string, err error) {
	resultGetBlockVerbose, err := GetDcrnBlock(b, blockHash)
	if err != nil {
		return nil, err
	}
	transactions := resultGetBlockVerbose.Tx
	if len(transactions) <= 1 {
		return
	}
	result = transactions[1:] //除第一个交易其余交易才是普通交易
	return
}

// GetBlock
// 根据区块hash查询区块信息
func GetBlock(b tokens.CrossChainBridge, blockHash string) (*electrs.ElectBlock, error) {
	gateway := b.GetGatewayConfig()
	var blockVerboseResult dtypes.GetBlockVerboseResult
	var err error
	for _, apiAddress := range gateway.APIAddress {
		err = CallGet(&blockVerboseResult, apiAddress, "getblock", blockHash)
		if err == nil {
			result := GetBlockVerboseResult2ElectBlock(&blockVerboseResult)
			return result, nil
		}
	}
	return nil, err
}

// GetDcrnBlock
func GetDcrnBlock(b tokens.CrossChainBridge, blockHash string) (*dtypes.GetBlockVerboseResult, error) {
	gateway := b.GetGatewayConfig()
	var result dtypes.GetBlockVerboseResult
	var err error
	for _, apiAddress := range gateway.APIAddress {
		err = CallGet(&result, apiAddress, "getblock", blockHash)
		if err == nil {
			return &result, nil
		}
	}
	return nil, err
}

// GetBlockTransactions
//注：无视startIndex
func GetBlockTransactions(b tokens.CrossChainBridge, blockHash string, startIndex uint32) (result []*electrs.ElectTx, err error) {
	resultGetBlockVerbose, err := GetDcrnBlock(b, blockHash)
	if err != nil {
		return nil, err
	}
	transactions := resultGetBlockVerbose.Tx
	if len(transactions) <= 1 {
		return
	}
	trans := transactions[1:] //除第一个交易其余交易才是普通交易
	for _, tran := range trans {
		tx, err := GetTransactionByHash(b, tran)
		if err == nil {
			// result = append(result, tx)
			fmt.Println("tx:", tx)
		}
	}
	return
}

// EstimateFeePerKb call /fee-estimates and multiply 1000
func EstimateFeePerKb(b tokens.CrossChainBridge, blocks int) (fee int64, err error) {
	gateway := b.GetGatewayConfig()
	var result float64
	for _, apiAddress := range gateway.APIAddress {
		err = CallGet(&result, apiAddress, "estimatefee", blocks)
		if err == nil {
			fee = int64(result * 1e8)
			return
		}
	}
	return 0, err
}

//创建交易
func CreateRawTransaction() {

}

//签名交易
func SignRawtransaction(b tokens.CrossChainBridge, hex string) (signedHex string, txHash string, err error) {
	gateway := b.GetGatewayConfig()
	var resultSignRawTransaction types.SignRawTransactionResult
	for _, apiAddress := range gateway.APIAddress {
		err := CallGet(&resultSignRawTransaction, apiAddress, "signrawtransaction", hex)
		if err == nil {
			return resultSignRawTransaction.Hex, "", nil
		}
	}
	return signedHex, "", err
}
