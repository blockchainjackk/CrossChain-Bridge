package dcrn

import (
	"encoding/hex"
	"errors"
	"regexp"
	"strings"

	"decred.org/dcrwallet/wallet/txauthor"
	"github.com/anyswap/CrossChain-Bridge/common"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/btc/electrs"
	"github.com/anyswap/CrossChain-Bridge/tokens/dcrn/ctl"
)

var (
	regexMemo = regexp.MustCompile(`^OP_RETURN OP_PUSHBYTES_\d* `)
)

// GetTransaction impl
func (b *Bridge) GetTransaction(txHash string) (interface{}, error) {
	return b.GetTransactionByHash(txHash)
}

//todo
// GetTransactionStatus impl
func (b *Bridge) GetTransactionStatus(txHash string) (*tokens.TxStatus, error) {
	txStatus := &tokens.TxStatus{}
	electStatus, err := b.GetElectTransactionStatus(txHash)
	if err != nil {
		log.Trace(b.ChainConfig.BlockChain+" Bridge::GetElectTransactionStatus fail", "tx", txHash, "err", err)
		return txStatus, err
	}
	if !*electStatus.Confirmed {
		return txStatus, tokens.ErrTxNotStable
	}
	if electStatus.BlockHash != nil {
		txStatus.BlockHash = *electStatus.BlockHash
	}
	if electStatus.BlockTime != nil {
		txStatus.BlockTime = *electStatus.BlockTime
	}
	if electStatus.BlockHeight != nil {
		txStatus.BlockHeight = *electStatus.BlockHeight
		latest, errt := b.GetLatestBlockNumber()
		if errt == nil {
			//爆块就相当于 Confirmations=1
			if latest >= txStatus.BlockHeight {
				txStatus.Confirmations = latest - txStatus.BlockHeight + 1
			}
		}
	}
	return txStatus, nil
}

// VerifyMsgHash verify msg hash
func (b *Bridge) VerifyMsgHash(rawTx interface{}, msgHash []string) (err error) {
	authoredTx, ok := rawTx.(*txauthor.AuthoredTx)
	if !ok {
		return tokens.ErrWrongRawTx
	}
	for i, preScript := range authoredTx.PrevScripts {
		sigScript := preScript
		if b.IsPayToScriptHash(sigScript) {
			//todo   受影响，
			sigScript, err = b.getRedeemScriptByOutputScrpit(preScript)
			if err != nil {
				return err
			}
		}
		// todo  受影响
		sigHash, err := b.CalcSignatureHash(sigScript, authoredTx.Tx, i)
		if err != nil {
			return err
		}
		if hex.EncodeToString(sigHash) != msgHash[i] {
			log.Trace("message hash mismatch", "index", i, "want", msgHash[i], "have", hex.EncodeToString(sigHash))
			return tokens.ErrMsgHashMismatch
		}
	}
	return nil
}

// VerifyTransaction impl
func (b *Bridge) VerifyFormTransaction(pairID string, params *ctl.SwapInParam, allowUnstable bool) (*tokens.TxSwapInfo, error) {
	if !b.IsSrc {
		return nil, tokens.ErrBridgeDestinationNotSupported
	}
	fromAddress := params.FromAddress
	txHash := params.TxID
	toAddress := params.ToAddress
	signMsg := params.SignMsg
	if fromAddress != "" {
		fromVerify := b.verifyFrom(txHash, fromAddress)
		if !fromVerify {
			return nil, errors.New("fromAddress verify fail")
		}
	}
	if signMsg != "" {
		bindVerify := b.verifyBind(fromAddress, signMsg, toAddress)
		if !bindVerify {
			return nil, errors.New("toAddress verify fail")
		}
	}
	return b.verifySwapinTx(pairID, fromAddress, txHash, toAddress, allowUnstable)
}

// VerifyTransaction impl
func (b *Bridge) VerifyTransaction(pairID, txHash string, allowUnstable bool) (*tokens.TxSwapInfo, error) {
	if !b.IsSrc {
		return nil, tokens.ErrBridgeDestinationNotSupported
	}
	return b.verifySwapinTx(pairID, "", txHash, "", allowUnstable)
}

func (b *Bridge) verifySwapinTx(pairID, from, txHash, bindAddr string, allowUnstable bool) (*tokens.TxSwapInfo, error) {
	tokenCfg := b.GetTokenConfig(pairID)
	if tokenCfg == nil {
		return nil, tokens.ErrUnknownPairID
	}
	if tokenCfg.DisableSwap {
		return nil, tokens.ErrSwapIsClosed
	}
	swapInfo := &tokens.TxSwapInfo{}
	swapInfo.PairID = pairID // PairID
	swapInfo.Hash = txHash   // Hash
	if !allowUnstable && !b.checkStable(txHash) {
		return swapInfo, tokens.ErrTxNotStable
	}
	tx, err := b.GetTransactionByHash(txHash)
	if err != nil {
		log.Debug("[verifySwapin] "+b.ChainConfig.BlockChain+" Bridge::GetTransaction fail", "tx", txHash, "err", err)
		return swapInfo, tokens.ErrTxNotFound
	}
	txStatus := tx.Status
	if txStatus.BlockHeight != nil {
		swapInfo.Height = *txStatus.BlockHeight // Height
	} else if *tx.Locktime != 0 {
		// tx with locktime should be on chain, prvent DDOS attack
		return swapInfo, tokens.ErrTxNotStable
	}
	if txStatus.BlockTime != nil {
		swapInfo.Timestamp = *txStatus.BlockTime // Timestamp
	}
	depositAddress := tokenCfg.DepositAddress
	value, memoScript, rightReceiver := b.GetReceivedValue(tx.Vout, depositAddress, p2shType)
	if !rightReceiver {
		return swapInfo, tokens.ErrTxWithWrongReceiver
	}
	var bindAddress string
	var bindOk bool
	if bindAddr == "" {
		bindAddress, bindOk = GetBindAddressFromMemoScipt(memoScript)
		if !bindOk {
			log.Debug("wrong memo", "memo", memoScript)
			return swapInfo, tokens.ErrTxWithWrongMemo
		}
	} else {
		bindAddress = bindAddr
	}

	swapInfo.To = depositAddress                 // To
	swapInfo.Value = common.BigFromUint64(value) // Value
	swapInfo.Bind = bindAddress
	// swapInfo.From = getTxFrom(tx.Vin, depositAddress) // From
	swapInfo.From = from //直接使用前端入参

	err = b.checkSwapinInfo(swapInfo)
	if err != nil {
		return swapInfo, err
	}

	if !allowUnstable {
		log.Info("verify swapin pass", "pairID", swapInfo.PairID, "from", swapInfo.From, "to", swapInfo.To, "bind", swapInfo.Bind, "value", swapInfo.Value, "txid", swapInfo.Hash, "height", swapInfo.Height, "timestamp", swapInfo.Timestamp)
	}
	return swapInfo, nil
}

func (b *Bridge) checkSwapinInfo(swapInfo *tokens.TxSwapInfo) error {
	if swapInfo.From == swapInfo.To {
		log.Error("wrong swap sender in swapin, from==to!")
		return tokens.ErrTxWithWrongSender
	}
	if !tokens.CheckSwapValue(swapInfo, b.IsSrc) {
		log.Error("wrong swap value in swapin")
		return tokens.ErrTxWithWrongValue
	}
	if !tokens.DstBridge.IsValidAddress(swapInfo.Bind) {
		log.Debug("wrong bind address in swapin", "bind", swapInfo.Bind)
		return tokens.ErrTxWithWrongMemo
	}
	return nil
}

func (b *Bridge) checkStable(txHash string) bool {
	txStatus, err := b.GetTransactionStatus(txHash)
	if err != nil {
		return false
	}
	confirmations := *b.GetChainConfig().Confirmations
	return txStatus.BlockHeight > 0 && txStatus.Confirmations >= confirmations
}

// GetReceivedValue get received value
func (b *Bridge) GetReceivedValue(vout []*electrs.ElectTxOut, receiver, pubkeyType string) (value uint64, memoScript string, rightReceiver bool) {
	for _, output := range vout {
		switch *output.ScriptpubkeyType {
		case opReturnType:
			memoScript = *output.ScriptpubkeyAsm
			continue
		case pubkeyType:
			if output.ScriptpubkeyAddress == nil || *output.ScriptpubkeyAddress != receiver {
				continue
			}
			rightReceiver = true
			value += *output.Value
		}
	}
	return value, memoScript, rightReceiver
}

// return priorityAddress if has it in Vin
// return the first address in Vin if has no priorityAddress
func getTxFrom(vin []*electrs.ElectTxin, priorityAddress string) string {
	from := ""
	for _, input := range vin {
		if input != nil &&
			input.Prevout != nil &&
			input.Prevout.ScriptpubkeyAddress != nil {
			if *input.Prevout.ScriptpubkeyAddress == priorityAddress {
				return priorityAddress
			}
			if from == "" {
				from = *input.Prevout.ScriptpubkeyAddress
			}
		}
	}
	return from
}

// GetBindAddressFromMemoScipt get bind address
func GetBindAddressFromMemoScipt(memoScript string) (bind string, ok bool) {
	parts := regexMemo.Split(memoScript, -1)
	if len(parts) != 2 {
		return "", false
	}
	memoHex := strings.TrimSpace(parts[1])
	memo := common.FromHex(memoHex)
	memoStr := string(memo)
	if memoStr == tokens.AggregateMemo {
		return "", false
	}
	if len(memo) <= len(tokens.LockMemoPrefix) {
		return "", false
	}
	if !strings.HasPrefix(memoStr, tokens.LockMemoPrefix) {
		return "", false
	}
	bind = string(memo[len(tokens.LockMemoPrefix):])
	return bind, true
}

// 根据txhash查出给质押地址打钱的来源地址，是否与fromAddress一致
func (b *Bridge) verifyFrom(txHash, fromAddress string) bool {
	txRaw, err := ctl.GetDcrnTransactionByHash(b, txHash)
	if err != nil {
		log.Warnf("txHash:%v GetDcrnTransactionByHash fail\n", txHash)
		return false
	}
	vinSlice := txRaw.Vin
	//循环txRaw交易信息的Vin数组内容
	for _, oneVin := range vinSlice {
		vinTxid := oneVin.Txid
		vinVout := oneVin.Vout
		//根据Vin中的txid查询vin的交易信息
		vinTxRaw, err := ctl.GetDcrnTransactionByHash(b, vinTxid)
		if err != nil {
			log.Warnf("vinTxid:%v GetDcrnTransactionByHash fail\n", vinTxid)
			return false
		}
		vinVouSlice := vinTxRaw.Vout
		//根据原交易中的Vin[]中的Vout(索引)直接在上一笔交易中的Vout[]中进行定位
		if int(vinVout) >= len(vinVouSlice) {
			log.Warnln("vinVout >= len(vinVouSlice)")
			return false
		}
		oneVinVout := vinVouSlice[vinVout]
		if oneVinVout.N != vinVout {
			//根据Vout中的参数N进行再次确认
			log.Warnln("oneVinVout.N != vinVout")
			return false
		}
		//在找到的vinVou中拿到Addresses与fromAddress对比
		addresses := oneVinVout.ScriptPubKey.Addresses
		for _, address := range addresses {
			if fromAddress == address {
				//有一个地址即可
				return true
			}
		}
	}
	return false
}

// 根据address（加密地址）与message（内容）,验证加密后的signature（加密后的内容）是否正确
func (b *Bridge) verifyBind(address, signature, message string) bool {
	bindVerify, err := ctl.Verifymessage(b, address, signature, message)
	if err != nil || !bindVerify {
		log.Warnln("verifyBind fail")
		return false
	} else {
		return true
	}
}
