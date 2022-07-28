package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"github.com/anyswap/CrossChain-Bridge/cmd/utils"
	"github.com/anyswap/CrossChain-Bridge/common"
	"github.com/anyswap/CrossChain-Bridge/common/hexutil"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/etc"
	"github.com/anyswap/CrossChain-Bridge/tokens/eth/abicoder"
	"github.com/anyswap/CrossChain-Bridge/tools/crypto"
	"github.com/anyswap/CrossChain-Bridge/tools/keystore"
	"github.com/anyswap/CrossChain-Bridge/types"
	"github.com/urfave/cli/v2"
	"math/big"
)

// ./swaptools sendethtx --gateway https://rinkeby.infura.io/v3/3d1d1da6d3df414d916dffc285cab3f3 --gasPrice 5000000001  --from 0x9ff2a4f6a478b435ADD643ee5e02b74e06D4f315 --to 0x98662A1e391697C844E6405f3A32A897413BAb8f --value 0 --privateKey ee260de5beeff23260bc2d52b7ae34c8d0cc0468ffea80aed5f7c8a083016ed3 --input 0xa9059cbb000000000000000000000000bd444ad8fe2b8a0626b1311774543023d50592b60000000000000000000000000000000000000000000000000000000005f5e100 --dryrun
//  multiSwapIn Input : 0x65663532000000000000000000000000000000000000000000000000000000000011111100000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000098968000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000003200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002a3078436232443932313445623732623739316661366531383546353763383334434431354243463230320000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000011c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000011c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000002070ff2adb07f09c49c5d3a0ed2754457557add221e02eba0ad7cda42e1f5213b60000000000000000000000000000000000000000000000000000000000000020a18bea6154d7f201a2e061a373f54cfdf0a4028c44838ca9388ad2becae7e34d000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000002067037967345b93a1fbda0a4c255193e4aee26a6d3433efed7af4f717e65fcbbb00000000000000000000000000000000000000000000000000000000000000202fbbfdbe48d6cf21d35273f47acc2df7ce39237ab5a4b5af8076c9476956cf69

//  ./swaptools sendethtx --gateway https://www.ethercluster.com/mordor --gasPrice 5000000001  --from 0x30542709eaEA5Db50bFD00B00FCcA150192c5a60 --to 0xCb2D9214Eb72b791fa6e185F57c834CD15BCF202 --value 0 --privateKey ee260de5beeff23260bc2d52b7ae34c8d0cc0468ffea80aed5f7c8a083016ed3 --input 0x65663532000000000000000000000000000000000000000000000000000000000011111100000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000098968000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000003200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002a3078436232443932313445623732623739316661366531383546353763383334434431354243463230320000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000011c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000011c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000002070ff2adb07f09c49c5d3a0ed2754457557add221e02eba0ad7cda42e1f5213b60000000000000000000000000000000000000000000000000000000000000020a18bea6154d7f201a2e061a373f54cfdf0a4028c44838ca9388ad2becae7e34d000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000002067037967345b93a1fbda0a4c255193e4aee26a6d3433efed7af4f717e65fcbbb00000000000000000000000000000000000000000000000000000000000000202fbbfdbe48d6cf21d35273f47acc2df7ce39237ab5a4b5af8076c9476956cf69
var (
	// nolint:lll // allow long line of example
	sendEthTxCommand = &cli.Command{
		Action:    sendEthTx,
		Name:      "sendethtx",
		Usage:     "send eth transaction",
		ArgsUsage: " ",
		Description: `
send eth tx command, sign tx with keystore and password file.

Example:

./swaptools sendethtx --gateway http://1.2.3.4:5555 --keystore ./UTC.json --password ./password.txt --from 0x1111111111111111111111111111111111111111 --to 0x2222222222222222222222222222222222222222 --value 1000000000000000000 --input 0x0123456789 --dryrun
`,
		Flags: []cli.Flag{
			utils.GatewayFlag,
			utils.KeystoreFileFlag,
			utils.PasswordFileFlag,
			utils.PrivateKeyFlag,
			senderFlag,
			receiverFlag,
			valueFlag,
			inputDataFlag,
			gasLimitFlag,
			gasPriceFlag,
			accountNonceFlag,
			dryRunFlag,
		},
	}
)

type ethTxSender struct {
	gateway      string
	keystoreFile string
	passwordFile string
	sender       string
	receiver     string
	dryRun       bool

	value      *big.Int
	input      []byte
	keyWrapper *keystore.Key
}

var (
	ethBridge *etc.Bridge
	ethSender = &ethTxSender{}
	ethExtra  = &tokens.EthExtraArgs{}
)

func decodePriKeyString(priKey string) (*ecdsa.PrivateKey, error) {

	//addr, err := hex.DecodeString(keyJSON.Address)
	//if err != nil {
	//	return err
	//}
	//Key := priKey[2:]
	privkey, err := crypto.HexToECDSA(priKey)
	if err != nil {
		return nil, err
	}
	return privkey, nil

	//k.Address = common.BytesToAddress(addr)
	//k.PrivateKey = privkey
}
func (ets *ethTxSender) initArgs(ctx *cli.Context) {
	ets.gateway = ctx.String(utils.GatewayFlag.Name)
	ets.keystoreFile = ctx.String(utils.KeystoreFileFlag.Name)
	ets.passwordFile = ctx.String(utils.PasswordFileFlag.Name)
	ets.sender = ctx.String(senderFlag.Name)
	ets.receiver = ctx.String(receiverFlag.Name)
	ets.dryRun = ctx.Bool(dryRunFlag.Name)

	//if ets.keystoreFile == "" || ets.passwordFile == "" {
	//	log.Fatal("must specify '-keystore' and '-password' flag")
	//}
	if ctx.IsSet(utils.PrivateKeyFlag.Name) {
		name := ctx.String(utils.PrivateKeyFlag.Name)
		key, err := decodePriKeyString(name)
		if err != nil {
			log.Fatalf("decode private key string err. %v", err)
		}
		ets.keyWrapper = &keystore.Key{}
		ets.keyWrapper.PrivateKey = key
	}
	if ets.sender == "" {
		log.Fatal("must specify '-from' flag")
	}

	if ctx.IsSet(valueFlag.Name) {
		value, err := common.GetBigIntFromStr(ctx.String(valueFlag.Name))
		if err != nil {
			log.Fatalf("wrong value. %v", err)
		}
		ets.value = value
	}
	if ctx.IsSet(inputDataFlag.Name) {
		ets.input = common.FromHex(ctx.String(inputDataFlag.Name))
	}

	if ctx.IsSet(gasLimitFlag.Name) {
		gasLimitValue := ctx.Uint64(gasLimitFlag.Name)
		ethExtra.Gas = &gasLimitValue
		log.Printf("gas limit is set to %v", gasLimitValue)
	}
	if ctx.IsSet(gasPriceFlag.Name) {
		gasPriceValue, err := common.GetBigIntFromStr(ctx.String(gasPriceFlag.Name))
		if err != nil {
			log.Fatalf("wrong gas price. %v", err)
		}
		ethExtra.GasPrice = gasPriceValue
		log.Printf("gas price is set to %v", gasPriceValue)
	}
	if ctx.IsSet(accountNonceFlag.Name) {
		nonceValue := ctx.Uint64(accountNonceFlag.Name)
		ethExtra.Nonce = &nonceValue
		log.Printf("account nonce is set to %v", nonceValue)
	}

	log.Info("initArgs finished", "gateway", ets.gateway,
		"from", ets.sender, "to", ets.receiver, "value", ets.value,
		"input", common.ToHex(ets.input), "dryRun", ets.dryRun)
}

func (ets *ethTxSender) doInit() {
	//var err error
	//ets.keyWrapper, err = tools.LoadKeyStore(ets.keystoreFile, ets.passwordFile)
	//if err != nil {
	//	log.Fatal("load keystore failed", "err", err)
	//}
	//
	//keyAddr := ets.keyWrapper.Address.String()
	//if !strings.EqualFold(keyAddr, ets.sender) {
	//	log.Fatal("sender mismatch", "sender", ets.sender, "keyAddr", keyAddr)
	//}
	//log.Info("load keystore success", "address", keyAddr)

	ets.initBridge()
}

func (ets *ethTxSender) initBridge() {
	ethBridge = etc.NewCrossChainBridge(true)
	ethBridge.ChainConfig = &tokens.ChainConfig{
		BlockChain: "ETHCLASSIC",
		NetID:      "Mordor",
	}
	ethBridge.GatewayConfig = &tokens.GatewayConfig{
		APIAddress: []string{ets.gateway},
	}
	ethBridge.VerifyChainID()
}

//    "ef52ef1d": "MultiSwapin(bytes32,address,uint256,address,uint8[],bytes32[],bytes32[])"
func (ets *ethTxSender) buildMultiSwapInTxInput(txHash string, receiver string, swapValue string,
	vs [][]uint8, rs [][]byte, ss [][]byte) []byte {

	//receiver := common.HexToAddress(args.Bind)

	funcHash := "ef52ef1d"
	input := abicoder.PackDataWithFuncHash([]byte(funcHash), txHash, receiver, swapValue, vs, rs, ss)

	return input
}

func (ets *ethTxSender) buildTx() (rawTx interface{}, err error) {
	args := &tokens.BuildTxArgs{
		From:  ets.sender,
		To:    ets.receiver,
		Value: ets.value,
		Input: &ets.input,
		Extra: &tokens.AllExtras{
			EthExtra: ethExtra,
		},
	}
	return ethBridge.BuildRawTransaction(args)
}

func sendEthTx(ctx *cli.Context) error {
	utils.SetLogger(ctx)
	ethSender.initArgs(ctx)

	ethSender.doInit()

	//todo
	// input :  0xa9059cbb000000000000000000000000bd444ad8fe2b8a0626b1311774543023d50592b60000000000000000000000000000000000000000000000000000000005f5e100

	rawTx, err := ethSender.buildTx()
	if err != nil {
		log.Fatal("BuildRawTransaction error", "err", err)
	}

	signedTx, txHash, err := ethBridge.SignTransactionWithPrivateKey(rawTx, ethSender.keyWrapper.PrivateKey)
	if err != nil {
		log.Fatal("SignTransaction failed", "err", err)
	}
	log.Info("SignTransaction success", "txHash", txHash)

	tx, _ := signedTx.(*types.Transaction)
	tx.PrintPretty()

	if !ethSender.dryRun {
		txHash, err = ethBridge.SendTransaction(signedTx)
		if err != nil {
			log.Error("SendTransaction failed", "err", err)
		}
		log.Infof("SendTransaction success , txHash %s\n", txHash)
	} else {
		log.Info("------------ dry run, does not sendtx -------------")
	}
	return nil
}

func GetSignResult(erc20Address, destAddress string, value int64, privateKeyStr string) (vs []uint8, rs []byte, ss []byte) {
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatal("GetSignResult err :", err)
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
		log.Fatal("GetSignResult err :", err)
	}
	fmt.Println(hexutil.Encode(signature))
	fmt.Println("r:", []uint8(signature[:32]))
	fmt.Println("s:", signature[32:64])
	fmt.Println("v:", signature[64:])
	rs = []uint8(signature[:32])
	ss = signature[32:64]
	vs = signature[64:]
	return
}

func AnyToByte(any interface{}) []byte {
	buff := new(bytes.Buffer)
	//数据写入buff
	binary.Write(buff, binary.BigEndian, any)
	return buff.Bytes()
}

func MultiSing(erc20Address, destAddress string, value int64) (vs [][]uint8, rs [][]byte, ss [][]byte) {

	const POANODEPRIVATE1 = "7f79bf8c1f83b7462a51bcb2a9a539f9c9fe54e87bfc53ae8d3d5cd6db4e24a3"
	const POANODEPRIVATE2 = "ea6564e665355498f2a0f3511707b5ee2ae84c7b8a3a26aa0aa54cfb08f8fdd4"

	v, r, s := GetSignResult(erc20Address, destAddress, value, POANODEPRIVATE1)
	vs = append(vs, v)
	rs = append(rs, r)
	ss = append(ss, s)

	v1, r1, s1 := GetSignResult(erc20Address, destAddress, value, POANODEPRIVATE2)
	vs = append(vs, v1)
	rs = append(rs, r1)
	ss = append(ss, s1)

	return

}
