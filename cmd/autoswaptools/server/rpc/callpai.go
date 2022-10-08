package rpc

import (
	"fmt"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/dcrn/ctl"
	"github.com/decred/dcrd/rpc/jsonrpc/types/v2"
	"strings"
)

func PostTransaction(b tokens.CrossChainBridge, txHex string) (txHash string, err error) {
	gateway := b.GetGatewayConfig()
	var success bool
	for _, apiAddress := range gateway.APIAddress {
		hash0, err0 := ctl.CallPost(apiAddress, "sendrawtransaction", txHex)
		if err0 == nil && !success {
			success = true
			txHash = strings.Trim(hash0, "\"")
			return
		} else if err0 != nil {
			err = err0
		}
	}
	return txHash, err
}

func SendFromAccount(b tokens.CrossChainBridge, account string, to string, value float64) (txHash string, err error) {
	gateway := b.GetGatewayConfig()
	var success bool
	for _, apiAddress := range gateway.APIAddress {
		hash0, err0 := ctl.CallPost(apiAddress, "sendfrom", account, to, value)
		if err0 == nil && !success {
			success = true
			txHash = strings.Trim(hash0, "\"")
			return
		} else if err0 != nil {
			err = err0
		}
	}
	return txHash, err
}

func GetAccountBalance(b tokens.CrossChainBridge, account string) (balance float64, err error) {
	gateway := b.GetGatewayConfig()
	var balanceRet GetBalanceResult
	for _, apiAddress := range gateway.APIAddress {
		err = ctl.CallGet(&balanceRet, apiAddress, "getbalance", account)
		if err == nil {
			for _, value := range balanceRet.Balances {
				if value.AccountName == account {
					balance = value.Spendable
					return
				}
			}
			return 0, fmt.Errorf("%v not exist", account)
		}
	}
	return 0, err
}

func GetDcrnTransactionByHash(b tokens.CrossChainBridge, txHash string) (*types.TxRawResult, error) {
	gateway := b.GetGatewayConfig()
	var result types.TxRawResult
	var err error
	for _, apiAddress := range gateway.APIAddress {
		//注意：要用getrawtransaction（节点启动时需要增加--txindex参数），不要用gettransaction
		err = ctl.CallGet(&result, apiAddress, "getrawtransaction", txHash, 1)
		if err == nil {
			return &result, nil
		}
	}
	return nil, err
}

func GetSignMsg(b tokens.CrossChainBridge, from string, msg string) (string, error) {
	gateway := b.GetGatewayConfig()
	var result string
	var err error
	for _, apiAddress := range gateway.APIAddress {
		err = ctl.CallGet(&result, apiAddress, "signmessage", from, msg)
		if err == nil {
			return result, nil
		}
	}
	return "", err
}
