package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/anyswap/CrossChain-Bridge/common"
	"github.com/anyswap/CrossChain-Bridge/tokens/eth/abicoder"
	"github.com/anyswap/CrossChain-Bridge/tools/crypto"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"io/ioutil"
	"log"
	"math/big"
	"strings"
)

const POANODEPRIVATE1 = "7f79bf8c1f83b7462a51bcb2a9a539f9c9fe54e87bfc53ae8d3d5cd6db4e24a3"
const POANODEPRIVATE2 = "ea6564e665355498f2a0f3511707b5ee2ae84c7b8a3a26aa0aa54cfb08f8fdd4"
const Erc20Contract = "0xCb2D9214Eb72b791fa6e185F57c834CD15BCF202"

const DESTADDRESS = "0xbD444ad8Fe2B8A0626B1311774543023D50592B6"
const TXHash = "0x412e70c6c6e0b38ec0c29885dc1cbbd2c9e269681659248f185e87e188e4622c"

const ABIPATH = "/Users/boom/bitmain/projects/crosschain/antpool/CrossChain-Bridge/cmd/swaptools/test/token.abi"

var SwapValue int64 = 100000000

func main() {

	vs, rs, ss := MultiSing(Erc20Contract, DESTADDRESS, SwapValue)

	fmt.Println("vs : ", vs)
	var rsHexArry []string
	for _, rs := range rs {
		hex := common.ToHex(rs)
		rsHexArry = append(rsHexArry, hex)
	}
	fmt.Println("rs : ", rsHexArry)

	var ssHesArry []string
	for _, ss := range ss {
		hex := common.ToHex(ss)
		ssHesArry = append(ssHesArry, hex)
	}
	fmt.Println("ss : ", ssHesArry)

	input := BuildMultiSwapInTxInput(TXHash, DESTADDRESS, SwapValue, Erc20Contract, vs, rs, ss)

	fmt.Println("MultiSwapInTxInput: ", common.ToHex(input))

	//Input := BuildSwapInTxInput(TXHash, DESTADDRESS, SwapValue)
	//
	//fmt.Println("SwapInTxInput: ", common.ToHex(Input))

}

//    "ef52ef1d": "MultiSwapin(bytes32,address,uint256,address,uint8[],bytes32[],bytes32[])"
func BuildMultiSwapInTxInput(txHash string, receiver string, swapValue int64, erc20Contract string,
	vs []uint8, rs [][]byte, ss [][]byte) []byte {

	funcHash := common.FromHex("0xef52ef1d")
	txHash1 := common.HexToHash(txHash)
	receiverAddress := common.HexToAddress(receiver)
	//value := big.NewInt(swapValue)
	value256 := common.LeftPadBytes(big.NewInt(swapValue).Bytes(), 32)
	value := new(big.Int).SetBytes(value256)

	var rsHexArry [][32]byte
	for _, rs := range rs {
		value := [32]byte{}
		copy(value[:], rs)
		rsHexArry = append(rsHexArry, value)
	}
	fmt.Println("rs : ", rsHexArry)

	var ssHesArry [][32]byte
	for _, ss := range ss {
		value := [32]byte{}
		copy(value[:], ss)
		ssHesArry = append(ssHesArry, value)
	}
	fmt.Println("ss : ", ssHesArry)

	erc20ContractAddress := common.HexToAddress(erc20Contract)

	abiJson, err := ioutil.ReadFile(ABIPATH)
	if err != nil {
		log.Fatalln(err)
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Fatalln(err)
	}
	Abiinput, err := myAbi.Pack("MultiSwapin", txHash1, receiverAddress, value, erc20ContractAddress,
		vs, rsHexArry, ssHesArry)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("ABI  MultiSwapInTxInput: ", common.ToHex(Abiinput))

	input := abicoder.PackDataWithFuncHash(funcHash, txHash1, receiverAddress, value, erc20ContractAddress,
		vs, rsHexArry, ssHesArry)

	return input
}

//"ec126c77": "Swapin(bytes32,address,uint256)",

func BuildSwapInTxInput(txHash string, receiver string, swapValue int64) []byte {

	funcHash := common.FromHex("0xec126c77")
	txHash1 := common.HexToHash(txHash)
	receiverAddress := common.HexToAddress(receiver)
	value256 := common.LeftPadBytes(big.NewInt(swapValue).Bytes(), 32)
	value := new(big.Int).SetBytes(value256)

	input := abicoder.PackDataWithFuncHash([]byte(funcHash), txHash1, receiverAddress, value)

	return input
}
func MultiSing(erc20Address, destAddress string, value int64) (vs []uint8, rs [][]byte, ss [][]byte) {

	v, r, s := GetSignResult(erc20Address, destAddress, value, POANODEPRIVATE1)
	vs = append(vs, v...)
	rs = append(rs, []byte(r[:32]))
	ss = append(ss, s)

	v1, r1, s1 := GetSignResult(erc20Address, destAddress, value, POANODEPRIVATE2)
	vs = append(vs, v1...)
	rs = append(rs, r1)
	ss = append(ss, s1)

	return

}

func GetSignResult(erc20Address, destAddress string, value int64, privateKeyStr string) (vs []uint8, rs []byte, ss []byte) {
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatal(err)
	}

	erc20 := common.HexToAddress(erc20Address)
	dest := common.HexToAddress(destAddress)

	erc20Byte := AnyToByte(erc20)
	destByte := AnyToByte(dest)
	value256 := common.LeftPadBytes(big.NewInt(value).Bytes(), 32)

	firstHash := crypto.Keccak256Hash(erc20Byte, destByte, value256)
	fmt.Println("firstHash:", firstHash.Hex())

	hash := crypto.Keccak256Hash([]byte("\x19Ethereum Signed Message:\n32"), firstHash[:])
	fmt.Println("hash:", hash)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	signature[64] += 27
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(hexutil.Encode(signature))
	//fmt.Println("r:", []uint8(signature[:32]))
	//fmt.Println("s:", signature[32:64])
	//fmt.Println("v:", signature[64:])
	rs = signature[:32]
	ss = signature[32:64]
	vs = []uint8(signature[64:])
	return
}

func AnyToByte(any interface{}) []byte {
	buff := new(bytes.Buffer)
	//数据写入buff
	binary.Write(buff, binary.BigEndian, any)
	return buff.Bytes()
}
