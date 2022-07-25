package tokens

import (
	"fmt"
	"math/big"
)

// SwapType type
type SwapType uint32

// SwapType constants
const (
	NoSwapType SwapType = iota
	SwapinType
	SwapoutType
)

func (s SwapType) String() string {
	switch s {
	case NoSwapType:
		return "noswap"
	case SwapinType:
		return "swapin"
	case SwapoutType:
		return "swapout"
	default:
		return fmt.Sprintf("unknown swap type %d", s)
	}
}

// SwapTxType type
type SwapTxType uint32

// SwapTxType constants
const (
	SwapinTx     SwapTxType = iota // 0
	SwapoutTx                      // 1
	P2shSwapinTx                   // 2
)

func (s SwapTxType) String() string {
	switch s {
	case SwapinTx:
		return "swapintx"
	case SwapoutTx:
		return "swapouttx"
	case P2shSwapinTx:
		return "p2shswapintx"
	default:
		return fmt.Sprintf("unknown swaptx type %d", s)
	}
}

//swapInfo.TxTo = strings.ToLower(receipt.Recipient.String()) // TxTo
//swapInfo.From = strings.ToLower(receipt.From.String())      // From
//
//from, to, value, err := ParseErc20SwapinTxLogs(receipt.Logs, token.ContractAddress, token.DepositAddress)
//if err != nil {
//if !errors.Is(err, tokens.ErrTxWithWrongReceiver) {
//log.Debug(b.ChainConfig.BlockChain+" ParseErc20SwapinTxLogs failed", "tx", swapInfo.Hash, "err", err)
//}
//return err
//}
//swapInfo.To = strings.ToLower(to)     // To
//swapInfo.Value = value                // Value
//swapInfo.Bind = strings.ToLower(from) // Bind

//  SwapOut
// TxSwapInfo struct
type TxSwapInfo struct {
	PairID    string `json:"pairid"`
	Hash      string `json:"hash"`
	Height    uint64 `json:"height"`
	Timestamp uint64 `json:"timestamp"`
	//strings.ToLower(receipt.From.String())
	From string `json:"from"`
	//swapInfo.TxTo = txRecipient                            // TxTo
	//swapInfo.To = txRecipient
	//泰达币（泰达公司发行的）合约地址                            // To
	TxTo string `json:"txto"`
	To   string `json:"to"`
	//topic【2】，MPC提供的收钱地址
	Bind  string   `json:"bind"`
	Value *big.Int `json:"value"`
}

/*
SwapIn Erc20
// TxSwapInfo struct
type SwapTxInfo struct {
	SwapInfo `json:"swapinfo"`
	SwapType SwapType `json:"swaptype"`
	Hash      string `json:"hash"`
	Height    uint64 `json:"height"`
	Timestamp uint64 `json:"timestamp"`
	//MPC提供的地址
	//swapInfo.From = strings.ToLower(receipt.From.String())
	From string `json:"from"`
	//swapInfo.TxTo = strings.ToLower(receipt.Recipient.String())
	//anyswapToken跨链合约地址
	TxTo string `json:"txto"`

	//合约事件
	//to   ： 跨链金额接收地址
	To string `json:"to"`
	//swapInfo.Bind = strings.ToLower(from) // Bind
	//from ： 跨链token合约中的0x000地址

	Bind string `json:"bind"`
	//跨链的金额
	Value       *big.Int `json:"value"`
	FromChainID *big.Int `json:"fromChainID"`
	ToChainID   *big.Int `json:"toChainID"`

	// 前端用户触发的跨链事件
	LogIndex int `json:"logIndex"`
}*/

// TxStatus struct
type TxStatus struct {
	Receipt       interface{} `json:"receipt,omitempty"`
	Confirmations uint64      `json:"confirmations"`
	BlockHeight   uint64      `json:"block_height"`
	BlockHash     string      `json:"block_hash"`
	BlockTime     uint64      `json:"block_time"`
}

// SwapInfo struct
type SwapInfo struct {
	PairID     string     `json:"pairid,omitempty"`
	SwapID     string     `json:"swapid,omitempty"`
	SwapType   SwapType   `json:"swaptype,omitempty"`
	TxType     SwapTxType `json:"txtype,omitempty"`
	Bind       string     `json:"bind,omitempty"`
	Identifier string     `json:"identifier,omitempty"`
	Reswapping bool       `json:"reswapping,omitempty"`
}

// IsSwapin is swapin type
func (s *SwapInfo) IsSwapin() bool {
	return s.SwapType == SwapinType
}

// BuildTxArgs struct
type BuildTxArgs struct {
	SwapInfo    `json:"swapInfo,omitempty"`
	From        string     `json:"from,omitempty"`
	To          string     `json:"to,omitempty"`
	OriginFrom  string     `json:"originFrom,omitempty"`
	OriginTxTo  string     `json:"originTxTo,omitempty"`
	Value       *big.Int   `json:"value,omitempty"`
	OriginValue *big.Int   `json:"originValue,omitempty"`
	SwapValue   *big.Int   `json:"swapvalue,omitempty"`
	Memo        string     `json:"memo,omitempty"`
	Input       *[]byte    `json:"input,omitempty"`
	Extra       *AllExtras `json:"extra,omitempty"`
}

// GetReplaceNum get rplace swap count
func (args *BuildTxArgs) GetReplaceNum() uint64 {
	if args.Extra != nil {
		return args.Extra.ReplaceNum
	}
	return 0
}

// GetExtraArgs get extra args
func (args *BuildTxArgs) GetExtraArgs() *BuildTxArgs {
	return &BuildTxArgs{
		SwapInfo: args.SwapInfo,
		Extra:    args.Extra,
	}
}

// GetTxGasPrice get tx gas price
func (args *BuildTxArgs) GetTxGasPrice() *big.Int {
	if args.Extra != nil && args.Extra.EthExtra != nil && args.Extra.EthExtra.GasPrice != nil {
		return args.Extra.EthExtra.GasPrice
	}
	return nil
}

// GetTxNonce get tx nonce
func (args *BuildTxArgs) GetTxNonce() uint64 {
	if args.Extra != nil && args.Extra.EthExtra != nil && args.Extra.EthExtra.Nonce != nil {
		return *args.Extra.EthExtra.Nonce
	}
	if args.Extra != nil && args.Extra.RippleExtra != nil && args.Extra.RippleExtra.Sequence != nil {
		return uint64(*args.Extra.RippleExtra.Sequence)
	}
	return 0
}

// AllExtras struct
type AllExtras struct {
	ReplaceNum  uint64        `json:"replaceNum,omitempty"`
	BtcExtra    *BtcExtraArgs `json:"btcExtra,omitempty"`
	EthExtra    *EthExtraArgs `json:"ethExtra,omitempty"`
	RippleExtra *RippleExtra  `json:"rippleExtra,omitempty"`
}

// EthExtraArgs struct
type EthExtraArgs struct {
	Gas       *uint64  `json:"gas,omitempty"`
	GasPrice  *big.Int `json:"gasPrice,omitempty"`
	GasTipCap *big.Int `json:"gasTipCap,omitempty"`
	GasFeeCap *big.Int `json:"gasFeeCap,omitempty"`
	Nonce     *uint64  `json:"nonce,omitempty"`
}

// RippleExtra struct
type RippleExtra struct {
	Sequence *uint32 `json:"sequence,omitempty"`
	Fee      *int64  `json:"fee,omitempty"`
}

// BtcOutPoint struct
type BtcOutPoint struct {
	Hash  string `json:"hash"`
	Index uint32 `json:"index"`
}

// BtcExtraArgs struct
type BtcExtraArgs struct {
	RelayFeePerKb     *int64         `json:"relayFeePerKb,omitempty"`
	ChangeAddress     *string        `json:"-"`
	PreviousOutPoints []*BtcOutPoint `json:"previousOutPoints,omitempty"`
}

// P2shAddressInfo struct
type P2shAddressInfo struct {
	BindAddress        string
	P2shAddress        string
	RedeemScript       string
	RedeemScriptDisasm string
}
