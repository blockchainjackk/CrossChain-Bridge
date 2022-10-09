package server

import (
	"crypto/ecdsa"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/db"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/db/ddl"
	"github.com/anyswap/CrossChain-Bridge/common/hexutil"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/tools/crypto"
)

func BscAccountInit(chainDB *db.CrossChainDB) {

	err := chainDB.CreateTable(ddl.BscAddressesTableName, ddl.CreateAddressTable)
	if err != nil {
		log.Errorf("BscAccountInit CreateTable err %v\n", err)
		return
	}

	count, err := chainDB.RetrieveAddressCount()
	if err != nil {
		log.Errorf("BscAccountInit RetrieveBalanceIndexCount err %v\n", err)
		return
	}
	for i := count; i < AutoSwapConf.MaxBscAccount; i++ {

		key, addr, err := CreateKey()
		if err != nil {
			log.Errorf("BscAccountInit err %v\n", err)
			return
		}

		err = chainDB.InsertAddress(key, addr, 0)
		if err != nil {
			log.Errorf("BscAccountInit InsertBalanceIndex err %v\n", err)
			return
		}
	}
	log.Info("BscAccountInit success .")
}

// CreateKey
func CreateKey() (privs, addr string, err error) {
	//创建私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Error("CreateKey fail: ", "err", err)
		return "", "", err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	priv := hexutil.Encode(privateKeyBytes)[2:]
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return priv, address, nil
}

func HexToPrivateKey(key string) (*ecdsa.PrivateKey, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, err
	}
	return privateKey, err
}
