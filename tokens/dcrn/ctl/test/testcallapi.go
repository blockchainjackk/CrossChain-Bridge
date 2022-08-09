package main

import (
	"fmt"

	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/btc"
	"github.com/anyswap/CrossChain-Bridge/tokens/dcrn/ctl"
)

var (
	err        error
	apiAddress string
	bridge     *btc.Bridge
)

func init() {
	apiAddress = "http://127.0.0.1:19557"
	gatewayConfig := &tokens.GatewayConfig{APIAddress: []string{apiAddress}}
	crossChainBridgeBase := &tokens.CrossChainBridgeBase{GatewayConfig: gatewayConfig}
	bridge = &btc.Bridge{CrossChainBridgeBase: crossChainBridgeBase}
	// bridge2 := &btc.Bridge{APIAddress: []string{"http://127.0.0.1:19557"}}
}

func main() {

	// 测试GetLatestBlockNumber
	testGetLatestBlockNumber()
	// 测试GetTransactionByHash
	testGetTransactionByHash()
	// 测试FindUtxos
	testFindUtxos()
	// 测试SignRawtransaction
	testSignRawtransaction()
	// 测试GetPoolTxidList
	testGetPoolTxidList()
	// 测试GetOutspend
	testGetOutspend()
	// 测试EstimateFeePerKb
	testEstimateFeePerKb()
	// 测试GetBlockHash
	testGetBlockHash()
	// 测试GetBlockTxids
	testGetBlockTxids()
	// 测试GetBlock
	testGetBlock()
	// 测试GetBlockTransactions
	testGetBlockTransactions()
	// 测试EstimateFeePerKb
	testEstimateFeePerKb()
}

func testGetLatestBlockNumber() {
	result, err := ctl.GetLatestBlockNumber(bridge)
	if err == nil {
		fmt.Printf("result:%v;type:%T\n", result, result)
	} else {
		fmt.Println("err:", err)
	}
}

func testGetTransactionByHash() {
	txHash := "5168b9251b4df0c64dfad5f7a60ed4e9b88d261d2e9afbfa9cfafc524dce1059"
	result, err := ctl.GetTransactionByHash(bridge, txHash)
	if err == nil {
		fmt.Printf("result:%v;type:%T\n", *result, *result)
	} else {
		fmt.Println("err:", err)
	}
}

func testFindUtxos() {
	addr := "SsbwBYjTLrF4TGMXGAdaVKRGXyFrK3fmvMk"
	addrs := []string{addr}
	result, err := ctl.FindUtxos(bridge, addrs)
	if err == nil {
		fmt.Printf("result:%v;type:%T\n", &result, &result)
	} else {
		fmt.Println("err:", err)
	}
	for _, utxo := range result {
		fmt.Println("utxo", utxo)
	}

}

func testSignRawtransaction() {
	hex := "010000000112a9c3b8ebff5a3810123dfc465daac637d7efd0155c9273719e377153bcdce20000000000ffffffff0200e1f5050000000000001976a914000000000000000000000000000000000000000088ac80ce341d0000000000001976a9145332bb44a27de855edcf2005cdd5b554577e824588ac0000000000000000018011c3230000000000000000ffffffff00"
	//first
	signedHex, err := ctl.CallPost(apiAddress, "signrawtransaction", hex)
	if err == nil {
		fmt.Printf("signedHex:%v;type:%T\n", signedHex, signedHex)
	} else {
		fmt.Println("err:", err)
	}
	//second
	signedHex2, _, err := ctl.SignRawtransaction(bridge, hex)
	if err == nil {
		fmt.Printf("signedHex2:%v;type:%T\n", signedHex2, signedHex2)
	} else {
		fmt.Println("err:", err)
	}
}

func testGetPoolTxidList() {
	var result []string
	result, err := ctl.GetPoolTxidList(bridge)
	if err == nil {
		fmt.Printf("result:%v;type:%T\n", result, result)
	} else {
		fmt.Println("err:", err)
	}
}

func testGetOutspend() {
	txHash := "be5701c1a4a63edc77290c26f288f17e64ec89033a4dbad9b0c0012e19bf144b"
	vout := 0
	result, err := ctl.GetOutspend(bridge, txHash, uint32(vout))
	if err == nil {
		fmt.Printf("result:%v;type:%T\n", result, result)
	} else {
		fmt.Println("err:", err)
	}
}

func testGetBlockHash() {
	var height uint64 = 32
	result, err := ctl.GetBlockHash(bridge, height)
	if err == nil {
		fmt.Printf("result:%v;type:%T\n", result, result)
	} else {
		fmt.Println("err:", err)
	}
}

func testGetBlockTxids() {
	blockHash := "4087bb58bdb7dfe77fb527592b73c212d59423c337bde9d6da35fe55507e84a0"
	result, err := ctl.GetBlockTxids(bridge, blockHash)
	if err == nil {
		fmt.Printf("result:%v;type:%T\n", result, result)
	} else {
		fmt.Println("err:", err)
	}
}

func testGetBlock() {
	blockHash := "3781d2870f5a613ddca33e7bcac0f44e6fcdb139805a23bfcad3c74f7b7b0b68"
	result, err := ctl.GetBlock(bridge, blockHash)
	if err == nil {
		fmt.Printf("result:%v;type:%T\n", result, result)
	} else {
		fmt.Println("err:", err)
	}
}

func testGetBlockTransactions() {
	//first 该块中无普通交易
	blockHash1 := "3781d2870f5a613ddca33e7bcac0f44e6fcdb139805a23bfcad3c74f7b7b0b68"
	result, err := ctl.GetBlockTransactions(bridge, blockHash1, 0)
	if err == nil {
		fmt.Printf("result:%v;type:%T\n", result, result)
	} else {
		fmt.Println("err:", err)
	}
	//second 该块中有普通交易
	blockHash2 := "4087bb58bdb7dfe77fb527592b73c212d59423c337bde9d6da35fe55507e84a0"
	result2, err := ctl.GetBlockTransactions(bridge, blockHash2, 0)
	if err == nil {
		fmt.Printf("result2:%v;type:%T\n", result2, result2)
	} else {
		fmt.Println("err:", err)
	}
	for _, one := range result2 {
		fmt.Printf("one:%v;type:%T\n", one, one)
	}
}

func testEstimateFeePerKb() {
	// var result float64
	result, err := ctl.EstimateFeePerKb(bridge, 6)
	if err == nil {
		fmt.Printf("result:%v;type:%T\n", result, result)
	} else {
		fmt.Println("err:", err)
	}
}
