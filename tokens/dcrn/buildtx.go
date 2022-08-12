package dcrn

import (
	"decred.org/dcrwallet/wallet/txauthor"
	"decred.org/dcrwallet/wallet/txrules"
	"decred.org/dcrwallet/wallet/txsizes"
	"errors"
	"fmt"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/params"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/btc/electrs"
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/decred/dcrd/txscript/v3"
	"github.com/decred/dcrd/wire"
	"math/big"
	"time"
)

const (
	p2pkhType    = "p2pkh"
	p2shType     = "p2sh"
	opReturnType = "op_return"

	retryCount    = 3
	retryInterval = 3 * time.Second

	// generatedTxVersion is the version of the transaction being generated.
	// It is defined as a constant here rather than using the wire.TxVersion
	// constant since a change in the transaction version will potentially
	// require changes to the generated transaction.  Thus, using the wire
	// constant for the generated transaction version could allow creation
	// of invalid transactions for the updated version.
	generatedTxVersion = 1
)

type scriptChangeSource struct {
	version uint16
	script  []byte
}

func (src *scriptChangeSource) Script() ([]byte, uint16, error) {
	return src.script, src.version, nil
}

func (src *scriptChangeSource) ScriptSize() int {
	return len(src.script)
}

func (b *Bridge) getRelayFeePerKb() (estimateFee int64, err error) {
	for i := 0; i < retryCount; i++ {
		estimateFee, err = b.EstimateFeePerKb(cfgEstimateFeeBlocks)
		if err == nil {
			break
		}
		time.Sleep(retryInterval)
	}
	if err != nil {
		log.Warn("estimate smart fee failed", "err", err)
		return 0, err
	}
	if cfgPlusFeePercentage > 0 {
		estimateFee += estimateFee * int64(cfgPlusFeePercentage) / 100
	}
	if estimateFee > cfgMaxRelayFeePerKb {
		estimateFee = cfgMaxRelayFeePerKb
	} else if estimateFee < cfgMinRelayFeePerKb {
		estimateFee = cfgMinRelayFeePerKb
	}
	return estimateFee, nil
}

func updateExtraInfo(extra *tokens.BtcExtraArgs, txins []*wireTxInType) {
	if len(extra.PreviousOutPoints) > 0 {
		return
	}
	extra.PreviousOutPoints = make([]*tokens.BtcOutPoint, len(txins))
	for i, txin := range txins {
		point := txin.PreviousOutPoint
		extra.PreviousOutPoints[i] = &tokens.BtcOutPoint{
			Hash:  point.Hash.String(),
			Index: point.Index,
		}
	}
}

// BuildRawTransaction build raw tx
func (b *Bridge) BuildRawTransaction(args *tokens.BuildTxArgs) (rawTx interface{}, err error) {
	var (
		pairID        = args.PairID
		token         = b.GetTokenConfig(pairID)
		from          string
		to            string
		changeAddress string
		amount        *big.Int
		memo          string
		relayFeePerKb dcrnAmountType
	)

	if token == nil {
		return nil, fmt.Errorf("swap pair '%v' is not configed", pairID)
	}

	switch args.SwapType {
	case tokens.SwapinType:
		return nil, tokens.ErrSwapTypeNotSupported
	case tokens.SwapoutType:
		from = token.DcrmAddress          // from
		to = args.Bind                    // to
		changeAddress = token.DcrmAddress // change

		amount = tokens.CalcSwappedValue(pairID, args.OriginValue, false, args.OriginFrom, args.OriginTxTo) // amount
		memo = tokens.UnlockMemoPrefix + args.SwapID
	default:
		return nil, tokens.ErrUnknownSwapType
	}

	if from == "" {
		return nil, errors.New("no sender specified")
	}

	var extra *tokens.BtcExtraArgs
	if args.Extra == nil || args.Extra.BtcExtra == nil {
		extra = &tokens.BtcExtraArgs{}
		args.Extra = &tokens.AllExtras{BtcExtra: extra}
	} else {
		extra = args.Extra.BtcExtra
		if extra.ChangeAddress != nil && args.SwapType == tokens.NoSwapType {
			changeAddress = *extra.ChangeAddress
		}
	}

	if extra.RelayFeePerKb != nil {
		relayFeePerKb = dcrnAmountType(*extra.RelayFeePerKb)
	} else {
		relayFee, errf := b.getRelayFeePerKb()
		if errf != nil {
			return nil, errf
		}
		extra.RelayFeePerKb = &relayFee
		relayFeePerKb = dcrnAmountType(relayFee)
	}

	txOuts, err := b.getTxOutputs(to, amount, memo)
	if err != nil {
		return nil, err
	}

	var inputSource txauthor.InputSource
	inputSource = func(target dcrutil.Amount) (detail *txauthor.InputDetail, err error) {
		if len(extra.PreviousOutPoints) != 0 {
			//todo
			return b.getUtxos(from, target, extra.PreviousOutPoints)
		}
		//todo
		return b.selectUtxos(from, target)
	}

	script, version, err := b.GetPayToAddrScript(changeAddress)
	if err != nil {
		return nil, err
	}
	changeSource := &scriptChangeSource{
		version: version,
		script:  script,
	}

	authoredTx, err := b.NewUnsignedTransaction(txOuts, relayFeePerKb, inputSource, changeSource, false)
	if err != nil {
		return nil, err
	}

	updateExtraInfo(extra, authoredTx.Tx.TxIn)

	if args.SwapType != tokens.NoSwapType {
		args.Identifier = params.GetIdentifier()
	}

	return authoredTx, nil
}

//func makeScriptChangeSource(address string, version uint16, params *chaincfg.Params) (*scriptChangeSource, error) {
//	destinationAddress, err := dcrutil.DecodeAddress(address, params)
//	if err != nil {
//		return nil, err
//	}
//
//	var script []byte
//	if addr, ok := destinationAddress.(wallet.V0Scripter); ok && version == 0 {
//		script = addr.ScriptV0()
//	} else {
//		script, err = txscript.PayToAddrScript(destinationAddress)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	source := &scriptChangeSource{
//		version: version,
//		script:  script,
//	}
//
//	return source, nil
//}
func (b *Bridge) getTxOutputs(to string, amount *big.Int, memo string) (txOuts []*wireTxOutType, err error) {
	if amount != nil {
		err = b.addPayToAddrOutput(&txOuts, to, amount.Int64())
		if err != nil {
			return nil, err
		}
	}
	//todo
	if memo != "" {
		err = b.addMemoOutput(&txOuts, memo)
		if err != nil {
			return nil, err
		}
	}

	return txOuts, err
}

func (b *Bridge) addPayToAddrOutput(txOuts *[]*wireTxOutType, to string, amount int64) error {
	if amount <= 0 {
		return nil
	}
	pkscript, _, err := b.GetPayToAddrScript(to)
	if err != nil {
		return err
	}
	*txOuts = append(*txOuts, b.NewTxOut(amount, pkscript))
	return nil
}

func (b *Bridge) addMemoOutput(txOuts *[]*wireTxOutType, memo string) error {
	if memo == "" {
		return nil
	}
	nullScript, err := b.NullDataScript(memo)
	if err != nil {
		return err
	}
	*txOuts = append(*txOuts, b.NewTxOut(0, nullScript))
	return nil
}
func (b *Bridge) findUxtosWithRetry(from string) (utxos []*electrs.ElectUtxo, err error) {
	for i := 0; i < retryCount; i++ {
		utxos, err = b.FindUtxos(from)
		if err == nil {
			break
		}
		time.Sleep(retryInterval)
	}
	return utxos, err
}

func (b *Bridge) getTransactionByHashWithRetry(txid string) (tx *electrs.ElectTx, err error) {
	for i := 0; i < retryCount; i++ {
		tx, err = b.GetTransactionByHash(txid)
		if err == nil {
			break
		}
		time.Sleep(retryInterval)
	}
	return tx, err
}

func (b *Bridge) getOutspendWithRetry(point *tokens.BtcOutPoint) (outspend *electrs.ElectOutspend, err error) {
	for i := 0; i < retryCount; i++ {
		outspend, err = b.GetOutspend(point.Hash, point.Index)
		if err == nil {
			break
		}
		time.Sleep(retryInterval)
	}
	return outspend, err
}

func (b *Bridge) selectUtxos(from string, target dcrnAmountType) (detail *txauthor.InputDetail, err error) {
	p2pkhScript, _, err := b.GetPayToAddrScript(from)
	if err != nil {
		return nil, err
	}

	utxos, err := b.findUxtosWithRetry(from)
	if err != nil {
		return nil, err
	}

	var (
		tx      *electrs.ElectTx
		success bool
		total   dcrnAmountType
	)

	inputs := make([]*wire.TxIn, 0)
	scripts := make([][]byte, 0)
	detail = &txauthor.InputDetail{}
	for _, utxo := range utxos {
		value := dcrnAmountType(*utxo.Value)
		if !isValidValue(value) {
			continue
		}
		tx, err = b.getTransactionByHashWithRetry(*utxo.Txid)
		if err != nil {
			continue
		}
		if *utxo.Vout >= uint32(len(tx.Vout)) {
			continue
		}
		output := tx.Vout[*utxo.Vout]
		if *output.ScriptpubkeyType != p2pkhType {
			continue
		}
		if output.ScriptpubkeyAddress == nil || *output.ScriptpubkeyAddress != from {
			continue
		}

		txIn, errf := b.NewTxIn(*utxo.Txid, *utxo.Vout, int64(value), p2pkhScript)
		if errf != nil {
			continue
		}

		total += value
		inputs = append(inputs, txIn)
		scripts = append(scripts, p2pkhScript)

		if total >= target {
			success = true
			break
		}
	}
	detail.Amount = total
	detail.Inputs = inputs
	detail.Scripts = scripts
	if !success {
		err = fmt.Errorf("not enough balance, total %v < target %v", total, target)
		return nil, err
	}

	return detail, nil
}

func (b *Bridge) getUtxos(from string, target dcrnAmountType, prevOutPoints []*tokens.BtcOutPoint) (detail *txauthor.InputDetail, err error) {
	p2pkhScript, _, err := b.GetPayToAddrScript(from)
	if err != nil {
		return nil, err
	}

	var total dcrnAmountType
	inputs := make([]*wire.TxIn, 0)
	scripts := make([][]byte, 0)
	detail = &txauthor.InputDetail{}
	for _, point := range prevOutPoints {
		outspend, errf := b.getOutspendWithRetry(point)
		if errf != nil {
			return nil, errf
		}
		if *outspend.Spent {
			if outspend.Status != nil && outspend.Status.BlockHeight != nil {
				spentHeight := *outspend.Status.BlockHeight
				err = fmt.Errorf("out point (%v, %v) is spent at %v", point.Hash, point.Index, spentHeight)
			} else {
				err = fmt.Errorf("out point (%v, %v) is spent at txpool", point.Hash, point.Index)
			}
			return nil, nil
		}
		tx, errf := b.getTransactionByHashWithRetry(point.Hash)
		if errf != nil {
			return nil, errf
		}
		if point.Index >= uint32(len(tx.Vout)) {
			err = fmt.Errorf("out point (%v, %v) index overflow", point.Hash, point.Index)
			return nil, err
		}
		output := tx.Vout[point.Index]
		if *output.ScriptpubkeyType != p2pkhType {
			err = fmt.Errorf("out point (%v, %v) script pubkey type %v is not p2pkh", point.Hash, point.Index, *output.ScriptpubkeyType)
			return nil, err
		}
		if output.ScriptpubkeyAddress == nil || *output.ScriptpubkeyAddress != from {
			err = fmt.Errorf("out point (%v, %v) script pubkey address %v is not %v", point.Hash, point.Index, *output.ScriptpubkeyAddress, from)
			return nil, err
		}
		value := dcrnAmountType(*output.Value)
		if value == 0 {
			err = fmt.Errorf("out point (%v, %v) with zero value", point.Hash, point.Index)
			return nil, err
		}

		txIn, errf := b.NewTxIn(point.Hash, point.Index, int64(value), p2pkhScript)
		if errf != nil {
			return nil, errf
		}

		total += value
		inputs = append(inputs, txIn)
		scripts = append(scripts, p2pkhScript)
	}
	detail.Amount = total
	detail.Inputs = inputs
	detail.Scripts = scripts

	if total < target {
		err = fmt.Errorf("not enough balance, total %v < target %v", total, target)
		return nil, err
	}
	return detail, nil
}

type insufficientFundsError struct{}

func (insufficientFundsError) InputSourceError() {}
func (insufficientFundsError) Error() string {
	return "insufficient funds available to construct transaction"
}

func sumOutputValues(outputs []*wire.TxOut) (totalOutput dcrutil.Amount) {
	for _, txOut := range outputs {
		totalOutput += dcrutil.Amount(txOut.Value)
	}
	return totalOutput
}

// todo 换成DCRN的需要
// NewUnsignedTransaction ref btcwallet
// ref. https://github.com/btcsuite/btcwallet/blob/b07494fc2d662fdda2b8a9db2a3eacde3e1ef347/wallet/txauthor/author.go
// we only modify it to support P2PKH change script (the origin only support P2WPKH change script)
// and update estimate size because we are not use P2WKH
func (b *Bridge) NewUnsignedTransaction(outputs []*wireTxOutType, relayFeePerKb dcrnAmountType, fetchInputs txauthor.InputSource, fetchChange txauthor.ChangeSource, isAggregate bool) (*txauthor.AuthoredTx, error) {

	//const op errors.Op = "txauthor.NewUnsignedTransaction"

	targetAmount := sumOutputValues(outputs)
	scriptSizes := []int{txsizes.RedeemP2PKHSigScriptSize}
	changeScript, changeScriptVersion, err := fetchChange.Script()
	if err != nil {
		return nil, err
	}
	changeScriptSize := fetchChange.ScriptSize()

	maxSignedSize := txsizes.EstimateSerializeSize(scriptSizes, outputs, changeScriptSize)
	targetFee := txrules.FeeForSerializeSize(relayFeePerKb, maxSignedSize)

	for {
		inputDetail, err := fetchInputs(targetAmount + targetFee)

		if err != nil {
			return nil, err
		}
		if inputDetail.Amount < targetAmount+targetFee {
			return nil, insufficientFundsError{}
		}

		scriptSizes := make([]int, 0, len(inputDetail.RedeemScriptSizes))
		scriptSizes = append(scriptSizes, inputDetail.RedeemScriptSizes...)

		maxSignedSize = txsizes.EstimateSerializeSize(scriptSizes, outputs, changeScriptSize)
		maxRequiredFee := txrules.FeeForSerializeSize(relayFeePerKb, maxSignedSize)

		if maxRequiredFee < dcrnAmountType(cfgMinRelayFee) {
			maxRequiredFee = dcrnAmountType(cfgMinRelayFee)
		}
		remainingAmount := inputDetail.Amount - targetAmount

		if remainingAmount < maxRequiredFee {
			if isAggregate {
				return nil, insufficientFundsError{}
			}
			targetFee = maxRequiredFee
			continue
		}

		unsignedTransaction := &wire.MsgTx{
			SerType:  wire.TxSerializeFull,
			Version:  generatedTxVersion,
			TxIn:     inputDetail.Inputs,
			TxOut:    outputs,
			LockTime: 0,
			Expiry:   0,
		}
		changeIndex := -1
		changeAmount := inputDetail.Amount - targetAmount - maxRequiredFee

		if changeAmount != 0 && !txrules.IsDustAmount(changeAmount,
			changeScriptSize, relayFeePerKb) {
			if len(changeScript) > txscript.MaxScriptElementSize {
				return nil, fmt.Errorf("script size exceed maximum bytes " +
					"pushable to the stack")
			}
			change := &wire.TxOut{
				Value:    int64(changeAmount),
				Version:  changeScriptVersion,
				PkScript: changeScript,
			}
			l := len(outputs)
			unsignedTransaction.TxOut = append(outputs[:l:l], change)
			changeIndex = l
		} else {
			maxSignedSize = txsizes.EstimateSerializeSize(scriptSizes,
				unsignedTransaction.TxOut, 0)
		}

		return &txauthor.AuthoredTx{
			Tx:                           unsignedTransaction,
			PrevScripts:                  inputDetail.Scripts,
			TotalInput:                   inputDetail.Amount,
			ChangeIndex:                  changeIndex,
			EstimatedSignedSerializeSize: maxSignedSize,
		}, nil
	}
}

//func (b *Bridge) estimateSize(scripts [][]byte, txOuts []*wireTxOutType, addChangeOutput, isAggregate bool) int {
//	if !isAggregate {
//		return txsizes.EstimateSerializeSize(len(scripts), txOuts, addChangeOutput)
//	}
//
//	var p2sh, p2pkh int
//	for _, pkScript := range scripts {
//		switch {
//		case b.IsPayToScriptHash(pkScript):
//			p2sh++
//		default:
//			p2pkh++
//		}
//	}
//
//	size := txsizes.EstimateSerializeSize(p2pkh, txOuts, addChangeOutput)
//	if p2sh > 0 {
//		size += p2sh * redeemAggregateP2SHInputSize
//	}
//
//	return size
//}
