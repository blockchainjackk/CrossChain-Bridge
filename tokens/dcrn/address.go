package dcrn

import (
	"github.com/decred/dcrd/dcrec"
	"github.com/decred/dcrd/dcrutil/v3"
)

// DecodeAddress decode address
func (b *Bridge) DecodeAddress(addr string) (address dcrutil.Address, err error) {
	chainConfig := b.Inherit.GetChainParams()
	address, err = dcrutil.DecodeAddress(addr, chainConfig)
	if err != nil {
		return
	}
	//todo
	// dcrn has no IsForNet 判断地址是否符合此网络
	//可以通过反推出netID，和参数中的netID做对比，参见dcrutil.DecodeAddress(addr, chainConfig)实现
	//decoded, netID, err := base58.CheckDecode(addr)
	//if !address.IsForNet(chainConfig) {
	//	err = fmt.Errorf("invalid address for net")
	//	return
	//}
	return
}

// NewAddressPubKeyHash encap
func (b *Bridge) NewAddressPubKeyHash(pkData []byte) (*dcrutil.AddressPubKeyHash, error) {
	//todo 比 btc多一个签名算法类型参数，每一个如何处理？
	return dcrutil.NewAddressPubKeyHash(dcrutil.Hash160(pkData), b.Inherit.GetChainParams(), dcrec.STEcdsaSecp256k1)
}

// NewAddressScriptHash encap
func (b *Bridge) NewAddressScriptHash(redeemScript []byte) (*dcrutil.AddressScriptHash, error) {
	return dcrutil.NewAddressScriptHash(redeemScript, b.Inherit.GetChainParams())
}

// IsValidAddress check address
func (b *Bridge) IsValidAddress(addr string) bool {
	_, err := b.DecodeAddress(addr)
	return err == nil
}

// IsP2pkhAddress check p2pkh addrss
func (b *Bridge) IsP2pkhAddress(addr string) bool {
	address, err := b.DecodeAddress(addr)
	if err != nil {
		return false
	}
	_, ok := address.(*dcrutil.AddressPubKeyHash)
	return ok
}

// DecodeWIF decode wif
//todo 比比特币多一个net参数
func DecodeWIF(wif string, net [2]byte) (*dcrutil.WIF, error) {
	return dcrutil.DecodeWIF(wif, net)
}
