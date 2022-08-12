package ctl

import (
	"fmt"
	"strconv"

	"github.com/anyswap/CrossChain-Bridge/tokens/btc/electrs"
	"github.com/decred/dcrd/rpc/jsonrpc/types/v2"
)

const (
	p2pkhType    = "p2pkh"
	p2shType     = "p2sh"
	opReturnType = "op_return"
)

// DcrnTxStatus2ElectTxStatus converts dcrn TxStatus to elect TxStatus
func DcrnTxStatus2ElectTxStatus(status *DcrnTxStatus) electrs.ElectTxStatus {
	return electrs.ElectTxStatus{
		Confirmed:   &status.Confirmed,
		BlockHeight: &status.BlockHeight,
		BlockHash:   &status.BlockHash,
		BlockTime:   &status.BlockTime,
	}
}

// TxStatus make elect tx status from dcrn tx raw result
func TxStatus(tx *types.TxRawResult) *electrs.ElectTxStatus {
	status := &electrs.ElectTxStatus{
		Confirmed:   new(bool),
		BlockHeight: new(uint64),
		BlockHash:   new(string),
		BlockTime:   new(uint64),
	}
	*status.Confirmed = tx.Confirmations > 6
	*status.BlockHeight = uint64(tx.BlockHeight)
	*status.BlockHash = tx.BlockHash
	*status.BlockTime = uint64(tx.Blocktime)
	return status
}

// UnspentStatus make elect tx status from dcrn list unspent result
func UnspentStatus(unspent *ListUnspentResult) *electrs.ElectTxStatus {
	status := &electrs.ElectTxStatus{
		Confirmed: new(bool),
		//查看ElectUtxo中的ElectTxStatus只用到了Confirmed属性
	}
	*status.Confirmed = unspent.Confirmations > 6
	return status
}

// TxRawResult2ElectTx converts dcrn TxRawResult to elect Tx
func TxRawResult2ElectTx(tx *types.TxRawResult) *electrs.ElectTx {
	etx := &electrs.ElectTx{
		Txid:     &tx.Txid,
		Version:  new(uint32),
		Locktime: new(uint32),
		Size:     new(uint32),
		Weight:   new(uint32),
		Fee:      new(uint64),
		Vin:      make([]*electrs.ElectTxin, 0),
		Vout:     make([]*electrs.ElectTxOut, 0),
		Status:   TxStatus(tx),
	}
	*etx.Version = uint32(tx.Version)
	*etx.Locktime = uint32(tx.Time)
	// *etx.Size = uint32(tx.Size)
	// *etx.Weight = uint32(tx.Weight)
	for i := 0; i < len(tx.Vin); i++ {
		evin := ConvertVin(&tx.Vin[i])
		etx.Vin = append(etx.Vin, evin)
	}
	for j := 0; j < len(tx.Vout); j++ {
		evout := ConvertVout(&tx.Vout[j])
		etx.Vout = append(etx.Vout, evout)
	}
	return etx
}

// GetBlockVerboseResult2ElectBlock converts dcrn GetBlockVerboseResult to elect Block
func GetBlockVerboseResult2ElectBlock(block *types.GetBlockVerboseResult) *electrs.ElectBlock {
	eblk := &electrs.ElectBlock{
		Hash:         new(string),
		Height:       new(uint32),
		Version:      new(uint32),
		Timestamp:    new(uint32),
		TxCount:      new(uint32),
		Size:         new(uint32),
		Weight:       new(uint32),
		MerkleRoot:   new(string),
		PreviousHash: new(string),
		Nonce:        new(uint32),
		Bits:         new(uint32),
		Difficulty:   new(uint64),
	}
	*eblk.Hash = block.Hash
	*eblk.Height = uint32(block.Height)
	*eblk.Version = uint32(block.Version)
	*eblk.Timestamp = uint32(block.Time)
	*eblk.TxCount = uint32(len(block.Tx) + len(block.STx))
	*eblk.Size = uint32(block.Size)
	// *eblk.Weight = uint32(block.Weight)
	*eblk.MerkleRoot = block.MerkleRoot
	*eblk.PreviousHash = block.PreviousHash
	*eblk.Nonce = block.Nonce
	if bits, err := strconv.ParseUint(block.Bits, 16, 32); err == nil {
		*eblk.Bits = uint32(bits)
	}
	*eblk.Difficulty = uint64(block.Difficulty)
	return eblk
}

// ConvertVin converts dcrn vin to elect vin
func ConvertVin(vin *types.Vin) *electrs.ElectTxin {
	evin := &electrs.ElectTxin{
		Txid:         &vin.Txid,
		Vout:         &vin.Vout,
		Scriptsig:    new(string),
		ScriptsigAsm: new(string),
		IsCoinbase:   new(bool),
		Sequence:     &vin.Sequence,
		Prevout:      new(electrs.ElectTxOut),
	}
	if vin.ScriptSig != nil {
		*evin.Scriptsig = vin.ScriptSig.Hex
		*evin.ScriptsigAsm = vin.ScriptSig.Asm
	}
	*evin.IsCoinbase = (vin.Coinbase != "")
	//Prevout参数的填充后续由AddPrevout函数实现
	return evin
}

// ConvertVout converts dcrn vout to elect vout
func ConvertVout(vout *types.Vout) *electrs.ElectTxOut {
	evout := &electrs.ElectTxOut{
		Scriptpubkey:        &vout.ScriptPubKey.Hex,
		ScriptpubkeyAsm:     &vout.ScriptPubKey.Asm,
		ScriptpubkeyType:    new(string),
		ScriptpubkeyAddress: new(string),
		Value:               new(uint64),
	}
	switch vout.ScriptPubKey.Type {
	case "pubkeyhash":
		*evout.ScriptpubkeyType = p2pkhType
	case "scripthash":
		*evout.ScriptpubkeyType = p2shType
	default:
		*evout.ScriptpubkeyType = opReturnType
	}
	if len(vout.ScriptPubKey.Addresses) == 1 {
		*evout.ScriptpubkeyAddress = vout.ScriptPubKey.Addresses[0]
	}
	if len(vout.ScriptPubKey.Addresses) > 1 {
		*evout.ScriptpubkeyAddress = fmt.Sprintf("%+v", vout.ScriptPubKey.Addresses)
	}
	*evout.Value = uint64(vout.Value * 1e8)
	return evout
}

// SliceUnspentResult2ElectUtxo converts dcrn ListUnspentResult(slice) to elect Utxo(slice)
func SliceUnspentResult2ElectUtxo(unspentArray []*ListUnspentResult) []*electrs.ElectUtxo {
	var utxoArray []*electrs.ElectUtxo
	for _, unspent := range unspentArray {
		utxo := UnspentResult2ElectUtxo(unspent)
		utxoArray = append(utxoArray, utxo)
	}
	return utxoArray
}

// UnspentResult2ElectUtxo converts dcrn ListUnspentResult to elect Utxo
func UnspentResult2ElectUtxo(unspent *ListUnspentResult) *electrs.ElectUtxo {
	utxo := &electrs.ElectUtxo{
		Txid:   &unspent.TxID,
		Vout:   &unspent.Vout,
		Value:  new(uint64),
		Status: UnspentStatus(unspent),
	}
	*utxo.Value = uint64(unspent.Amount * 1e8) //*1e8不确定
	return utxo
}

// TxOutspend make elect outspend from dcrn GetTxOutResult raw result
func TxOutResult2ElectOutspend(txout *types.GetTxOutResult) *electrs.ElectOutspend {
	outspend := &electrs.ElectOutspend{
		Spent: new(bool),
	}
	if txout == nil {
		*outspend.Spent = true
	}
	return outspend
}
