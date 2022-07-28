package eth

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/anyswap/CrossChain-Bridge/common"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/eth/abicoder"
	"github.com/anyswap/CrossChain-Bridge/tools/crypto"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"io/ioutil"
	"math/big"
	"strings"
)

const anySwapV6TokenPath = "./tokens/eth/abicoder/abi/anySwapV6token.abi"

// build input for calling `Swapin(bytes32 txhash, address account, uint256 amount)`
func (b *Bridge) buildSwapinTxInput(args *tokens.BuildTxArgs) (err error) {
	token := b.GetTokenConfig(args.PairID)
	if token == nil {
		return tokens.ErrUnknownPairID
	}

	receiver := common.HexToAddress(args.Bind)
	if receiver == (common.Address{}) || !common.IsHexAddress(args.Bind) {
		log.Warn("swapin to wrong address", "receiver", args.Bind)
		return errInvalidReceiverAddress
	}

	swapValue := tokens.CalcSwappedValue(args.PairID, args.OriginValue, true, args.OriginFrom, args.OriginTxTo)
	swapValue, err = b.adjustSwapValue(args, swapValue)
	if err != nil {
		return err
	}
	args.SwapValue = swapValue // swap value

	funcHash := getSwapinFuncHash()
	txHash := common.HexToHash(args.SwapID)
	input := abicoder.PackDataWithFuncHash(funcHash, txHash, receiver, swapValue)
	args.Input = &input             // input
	args.To = token.ContractAddress // to

	if token.IsDelegateContract && !token.IsAnyswapAdapter {
		return b.checkBalance(token.DelegateToken, token.ContractAddress, swapValue)
	}
	return nil
}

// build input for calling `Swapin(bytes32 txhash, address account, uint256 amount)`
func (b *Bridge) buildMultiSwapinTxInput(args *tokens.BuildTxArgs) (err error) {
	token := b.GetTokenConfig(args.PairID)
	if token == nil {
		return tokens.ErrUnknownPairID
	}

	receiver := common.HexToAddress(args.Bind)
	if receiver == (common.Address{}) || !common.IsHexAddress(args.Bind) {
		log.Warn("swapin to wrong address", "receiver", args.Bind)
		return errInvalidReceiverAddress
	}

	swapValue := tokens.CalcSwappedValue(args.PairID, args.OriginValue, true, args.OriginFrom, args.OriginTxTo)
	swapValue, err = b.adjustSwapValue(args, swapValue)
	if err != nil {
		return err
	}
	args.SwapValue = swapValue // swap value

	txHash := common.HexToHash(args.SwapID)

	tokenContract := common.HexToAddress(token.ContractAddress)
	vs, rs, ss, err := b.MultiSign(tokenContract, receiver, swapValue, nil)

	abiJson, err := ioutil.ReadFile(anySwapV6TokenPath)
	if err != nil {
		log.Errorf("read abi file err", err)
		return err
	}

	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	if err != nil {
		log.Errorf("abi json err", err)
		return err
	}
	input, err := myAbi.Pack("MultiSwapin", txHash, receiver, swapValue, tokenContract,
		vs, rs, ss)
	if err != nil {
		log.Errorf("abi pack err", err)
		return err
	}
	args.Input = &input             // input
	args.To = token.ContractAddress // to

	if token.IsDelegateContract && !token.IsAnyswapAdapter {
		return b.checkBalance(token.DelegateToken, token.ContractAddress, swapValue)
	}
	return nil
}

func (b *Bridge) adjustSwapValue(args *tokens.BuildTxArgs, swapValue *big.Int) (*big.Int, error) {
	isDynamicFeeTx := b.ChainConfig.EnableDynamicFeeTx
	if isDynamicFeeTx {
		return swapValue, nil
	}

	if baseGasPrice == nil {
		return swapValue, nil
	}

	gasPrice := args.GetTxGasPrice()
	if gasPrice.Cmp(baseGasPrice) <= 0 {
		return swapValue, nil
	}

	fee := new(big.Int).Sub(args.OriginValue, swapValue)
	if fee.Sign() == 0 {
		return swapValue, nil
	}
	if fee.Sign() < 0 {
		return nil, tokens.ErrWrongSwapValue
	}

	extraGasPrice := new(big.Int).Sub(gasPrice, baseGasPrice)
	extraFee := new(big.Int).Mul(fee, extraGasPrice)
	extraFee.Div(extraFee, baseGasPrice)

	newSwapValue := new(big.Int).Sub(swapValue, extraFee)
	log.Info("adjust swap value", "isSrc", b.IsSrc, "chainID", b.SignerChainID,
		"pairID", args.PairID, "txid", args.SwapID, "bind", args.Bind, "swapType", args.SwapType.String(),
		"originValue", args.OriginValue, "oldSwapValue", swapValue, "newSwapValue", newSwapValue,
		"oldFee", fee, "extraFee", extraFee, "baseGasPrice", baseGasPrice, "gasPrice", gasPrice, "extraGasPrice", extraGasPrice)
	if newSwapValue.Sign() <= 0 {
		return nil, tokens.ErrWrongSwapValue
	}
	return newSwapValue, nil
}

func MultiSing(erc20Address, destAddress common.Address, value *big.Int) (vs []uint8, rs [][]byte, ss [][]byte) {
	const POANODEPRIVATE1 = "7f79bf8c1f83b7462a51bcb2a9a539f9c9fe54e87bfc53ae8d3d5cd6db4e24a3"
	const POANODEPRIVATE2 = "ea6564e665355498f2a0f3511707b5ee2ae84c7b8a3a26aa0aa54cfb08f8fdd4"

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

func GetSignResult(erc20Address, destAddress common.Address, value *big.Int, privateKeyStr string) (vs []uint8, rs []byte, ss []byte) {
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, nil, nil
	}

	erc20Byte := AnyToByte(erc20Address)
	destByte := AnyToByte(destAddress)
	value256 := common.LeftPadBytes(value.Bytes(), 32)

	firstHash := crypto.Keccak256Hash(erc20Byte, destByte, value256)
	fmt.Println("firstHash:", firstHash.Hex())

	hash := crypto.Keccak256Hash([]byte("\x19Ethereum Signed Message:\n32"), firstHash[:])
	fmt.Println("hash:", hash)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	signature[64] += 27
	if err != nil {
		return nil, nil, nil
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
