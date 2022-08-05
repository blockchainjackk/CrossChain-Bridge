package ctl

import "github.com/decred/dcrwallet/rpc/jsonrpc/types"

// SortableUnspentSlice sortable
type SortableUnspentSlice []*types.ListUnspentResult

// Len impl Sortable
func (s SortableUnspentSlice) Len() int {
	return len(s)
}

// Swap impl Sortable
func (s SortableUnspentSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less impl Sortable
// sort utxos
// 1. confirmed fisrt
// 2. value first
func (s SortableUnspentSlice) Less(i, j int) bool {
	spendable1 := s[i].Spendable
	spendable2 := s[j].Spendable
	if spendable1 != spendable2 {
		return spendable1
	}
	//为什么要逆序排序？
	return s[i].Amount > s[j].Amount
}

// DcrnTxStatus struct
type DcrnTxStatus struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight uint64 `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   uint64 `json:"block_time"`
}
