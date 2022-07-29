package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/anyswap/CrossChain-Bridge/params"
	"github.com/anyswap/CrossChain-Bridge/tokens/eth/contracts/keymanager"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethclient "github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	httpClientUrls               []string
	websocketClientUrls          []string
	keyManagerContractAddressHex string
	keyManagerOwnerPrivateKey    string
)

func initParam() {
	poaConfig := params.GetPoaConfig()
	httpClientUrls = poaConfig.HttpClientUrl
	websocketClientUrls = poaConfig.WebsocketClientUrl
	keyManagerContractAddressHex = poaConfig.KeyManagerContractAddressHex
	keyManagerOwnerPrivateKey = poaConfig.KeyManagerOwnerPrivateKey
}

//从私链的智能合约中获取N把私钥
//如果num=0,则（从链上）获取所有私钥
func GetPrivateKey(num uint8) (*[]string, error) {
	if uint8(len(privateKeysCache)) >= num && num > 0 {
		return GetPrivateKeyFromCache(num)
	}
	initParam()
	instance, err := getContractInstance()
	if err != nil {
		return nil, err
	}

	ownerAddress := GetAddressByPrivateKey(keyManagerOwnerPrivateKey)
	privateKeys, err := instance.GetPriKeys(&bind.CallOpts{From: ownerAddress}, num)
	if err != nil {
		return nil, err
	}
	privateKeysCache = privateKeys
	return &privateKeys, nil
}

//从私链的智能合约中获取所有私钥
func GetAllPriKeys() (*[]string, error) {
	if privateKeysCache != nil {
		return &privateKeysCache, nil
	}
	initParam()
	instance, err := getContractInstance()
	if err != nil {
		return nil, err
	}

	ownerAddress := GetAddressByPrivateKey(keyManagerOwnerPrivateKey)
	privateKeys, err := instance.GetAll(&bind.CallOpts{From: ownerAddress})
	if err != nil {
		return nil, err
	}
	privateKeysCache = privateKeys
	return &privateKeys, nil
}

func GetPrivateKeyFromCache(num uint8) (*[]string, error) {
	privateKeys := privateKeysCache[:2]
	return &privateKeys, nil
}

func AddPrivateKey(newPrivateKey string) (bool, error) {
	client, err := getClient()
	if err != nil {
		return false, err
	}

	instance, err := getContractInstanceByClient(client)
	if err != nil {
		return false, err
	}
	ownerAddress := GetAddressByPrivateKey(keyManagerOwnerPrivateKey)
	privateKey, _ := crypto.HexToECDSA(newPrivateKey)

	nonce, err := client.PendingNonceAt(context.Background(), ownerAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	_, err = instance.AddPriKey(auth, newPrivateKey)
	if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

func getContractInstance() (*keymanager.Keymanager, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	contractAddress := common.HexToAddress(keyManagerContractAddressHex)
	return keymanager.NewKeymanager(contractAddress, client)
}

func getContractInstanceByClient(client *ethclient.Client) (*keymanager.Keymanager, error) {
	contractAddress := common.HexToAddress(keyManagerContractAddressHex)
	return keymanager.NewKeymanager(contractAddress, client)
}

//以下两个函数属于tool类的函数，后面应放到更合适的包中
//根据私钥获取公钥
func GetPublicKeyByPrivateKey(hexKey string) string {
	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	return hex.EncodeToString(crypto.FromECDSAPub(publicKeyECDSA))
}

//根据私钥获取地址
func GetAddressByPrivateKey(hexKey string) common.Address {
	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

func getClient() (client *ethclient.Client, err error) {

	for _, httpClientUrl := range httpClientUrls {
		client, err = ethclient.Dial(httpClientUrl)
		if err == nil {
			fmt.Println("we hava poa chain connection")
			return client, nil
		}
	}
	return
}
