package ctl

// SortableUnspentSlice sortable
type SortableUnspentSlice []*ListUnspentResult

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

//以下三个struct（ListUnspentResult、SignRawTransactionError、SignRawTransactionResult）来自于github.com/decred/dcrwallet/rpc/jsonrpc/types
//由于与github.com/decred/dcrd/rpc/jsonrpc/types/v2包可能存在冲突问题，所以直接复制过来进行使用
// ListUnspentResult models a successful response from the listunspent request.
// Contains Decred additions.
type ListUnspentResult struct {
	TxID          string  `json:"txid"`
	Vout          uint32  `json:"vout"`
	Tree          int8    `json:"tree"`
	TxType        int     `json:"txtype"`
	Address       string  `json:"address"`
	Account       string  `json:"account"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	RedeemScript  string  `json:"redeemScript,omitempty"`
	Amount        float64 `json:"amount"`
	Confirmations int64   `json:"confirmations"`
	Spendable     bool    `json:"spendable"`
}

// SignRawTransactionError models the data that contains script verification
// errors from the signrawtransaction request.
type SignRawTransactionError struct {
	TxID      string `json:"txid"`
	Vout      uint32 `json:"vout"`
	ScriptSig string `json:"scriptSig"`
	Sequence  uint32 `json:"sequence"`
	Error     string `json:"error"`
}

// SignRawTransactionResult models the data from the signrawtransaction
// command.
type SignRawTransactionResult struct {
	Hex      string                    `json:"hex"`
	Complete bool                      `json:"complete"`
	Errors   []SignRawTransactionError `json:"errors,omitempty"`
}

type SwapInParam struct {
	FromAddress string
	TxHash      string
	BindAddress string
	SignMessage string
}
