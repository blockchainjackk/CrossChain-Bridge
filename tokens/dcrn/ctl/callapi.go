package ctl

import (
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/decred/dcrwallet/rpc/jsonrpc/types"
)

// GetTransactionByHash
//获取交易信息
func GetTransactionByHash(b tokens.CrossChainBridge, txHash string) (*types.GetTransactionResult, error) {
	gateway := b.GetGatewayConfig()
	var result types.GetTransactionResult
	var err error
	for _, apiAddress := range gateway.APIAddress {
		err := CallGet(&result, apiAddress, "gettransaction", txHash)
		if err == nil {
			return &result, nil
		}
	}
	return nil, err
}

// PostTransaction
//发送交易
func PostTransaction(b tokens.CrossChainBridge, txHex string) (txHash string, err error) {
	gateway := b.GetGatewayConfig()
	// var success bool
	for _, apiAddress := range gateway.APIAddress {
		// url := apiAddress + "/tx"
		// hash0, err0 := client.RPCRawPost(url, txHex)
		// if err0 == nil && !success {
		// 	success = true
		// 	txHash = hash0
		// } else if err0 != nil {
		// 	err = err0
		// }
		CallPost(apiAddress, "sendrawtransaction", txHex)

	}
	return txHash, err
}
