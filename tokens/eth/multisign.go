package eth

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/anyswap/CrossChain-Bridge/common"
	"github.com/anyswap/CrossChain-Bridge/common/hexutil"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/tools/crypto"
)

type MultiSignParam struct {
	MultiSignContract common.Address //暂不使用
	Erc20Contract     common.Address
	Receiver          common.Address
	Value             big.Int
	Nonce             big.Int //暂不使用
}

const minSignNum = 2

func multiSigns(param MultiSignParam) ([][]byte, error) {

	privateKeyStrs, err := getPrivateKeys()
	if err != nil {
		return nil, err
	}
	length := len(privateKeyStrs)
	var signatures [][]byte = make([][]byte, length)
	for i := 0; i < length; i++ {
		signature, err := multiSignOne(param, privateKeyStrs[i])
		if err != nil {
			return nil, err
		}
		signatures[i] = signature
	}
	return signatures, nil
}

func getPrivateKeys() ([]string, error) {

	privateKeys, err := GetAllPriKeys()
	if err != nil {
		return nil, err
	}
	return *privateKeys, nil
}

func multiSignOne(param MultiSignParam, privateKeyStr string) ([]byte, error) {
	erc20contractByte := param.Erc20Contract.Bytes()
	receiverByte := param.Receiver.Bytes()
	value := common.LeftPadBytes(param.Value.Bytes(), 32)

	firstHash := crypto.Keccak256Hash(erc20contractByte, receiverByte, value)
	hash := crypto.Keccak256Hash([]byte("\x19Ethereum Signed Message:\n32"), firstHash[:])

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatalf("privateKey parse fail")
		return nil, err
	}

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatalf("multiSign fail")
		return nil, err
	}
	signature[64] += 27
	fmt.Println("signature:", hexutil.Encode(signature))
	return signature, nil
}

func sigs2rsv(sigs [][]byte) (rs, ss [][32]byte, vs []uint8) {
	length := len(sigs)
	rs = make([][32]byte, length)
	ss = make([][32]byte, length)
	vs = make([]uint8, length)

	for i, sig := range sigs {
		r, s, v := sig2rsv(sig)
		rs[i] = r
		ss[i] = s
		vs[i] = v
	}
	return
}

//签名后的结果拆分成R S V
func sig2rsv(sig []byte) (r, s [32]byte, v uint8) {
	copy(r[:], sig[:32])
	copy(s[:], sig[32:64])
	bytesBuffer := bytes.NewBuffer(sig[64:])
	var x uint8
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	v = x
	return
}
