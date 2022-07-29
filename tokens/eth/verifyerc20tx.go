package eth

import (
	"bytes"
	"errors"
	"math/big"
	"strings"

	"github.com/anyswap/CrossChain-Bridge/common"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/params"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/tools"
	"github.com/anyswap/CrossChain-Bridge/types"
)

// verifyErc20SwapinTx verify erc20 swapin with pairID
func (b *Bridge) verifyErc20SwapinTx(swapInfo *tokens.TxSwapInfo, allowUnstable bool, token *tokens.TokenConfig, receipt *types.RPCTxReceipt) (*tokens.TxSwapInfo, error) {
	err := b.verifyErc20SwapinTxReceipt(swapInfo, receipt, token)
	if err != nil {
		return swapInfo, err
	}

	err = b.checkSwapinInfo(swapInfo)
	if err != nil {
		return swapInfo, err
	}

	if !allowUnstable {
		log.Info("verify erc20 swapin stable pass",
			"identifier", params.GetIdentifier(), "pairID", swapInfo.PairID,
			"from", swapInfo.From, "to", swapInfo.To, "bind", swapInfo.Bind,
			"value", swapInfo.Value, "txid", swapInfo.Hash,
			"height", swapInfo.Height, "timestamp", swapInfo.Timestamp)
	}
	return swapInfo, nil
}

func (b *Bridge) verifyErc20SwapinTxReceipt(swapInfo *tokens.TxSwapInfo, receipt *types.RPCTxReceipt, token *tokens.TokenConfig) error {
	if receipt.Recipient == nil {
		return tokens.ErrTxWithWrongContract
	}
	//跨链合约地址
	swapInfo.TxTo = strings.ToLower(receipt.Recipient.String()) // TxTo
	//MPC提供的地址
	swapInfo.From = strings.ToLower(receipt.From.String()) // From
	//from ： 跨链token合约中的0x000地址
	//to   ： 跨链金额接收地址
	var depositAddress string
	if token.MultiSignContractDepositAddress == "" {
		depositAddress = token.DepositAddress
	} else {
		//如果配置了多签合约地址，就往多签合约地址进行质押
		depositAddress = token.MultiSignContractDepositAddress
	}
	from, to, value, err := ParseErc20SwapinTxLogs(receipt.Logs, token.ContractAddress, depositAddress)
	if err != nil {
		if !errors.Is(err, tokens.ErrTxWithWrongReceiver) {
			log.Debug(b.ChainConfig.BlockChain+" ParseErc20SwapinTxLogs failed", "tx", swapInfo.Hash, "err", err)
		}
		return err
	}
	swapInfo.To = strings.ToLower(to)     // To
	swapInfo.Value = value                // Value
	swapInfo.Bind = strings.ToLower(from) // Bind

	if !token.AllowSwapinFromContract &&
		!b.ChainConfig.AllowCallByContract &&
		!common.IsEqualIgnoreCase(swapInfo.TxTo, token.ContractAddress) {
		if err := b.checkCallByContract(swapInfo); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bridge) checkCallByContract(swapInfo *tokens.TxSwapInfo) error {
	if b.ChainConfig.IsInCallByContractWhitelist(swapInfo.TxTo) {
		return nil
	}
	if b.ChainConfig.HasCallByContractCodeHashWhitelist() {
		codehash := b.GetContractCodeHash(common.HexToAddress(swapInfo.TxTo))
		if codehash != (common.Hash{}) &&
			b.ChainConfig.IsInCallByContractCodeHashWhitelist(codehash.String()) {
			return nil
		}
	}
	return tokens.ErrTxWithWrongContract
}

// ParseErc20SwapinTxLogs parse erc20 swapin tx logs
func ParseErc20SwapinTxLogs(logs []*types.RPCLog, contractAddress, checkToAddress string) (from, to string, value *big.Int, err error) {
	transferLogExist := false
	for _, log := range logs {
		if log.Removed != nil && *log.Removed {
			continue
		}
		if !common.IsEqualIgnoreCase(log.Address.String(), contractAddress) {
			continue
		}
		if len(log.Topics) != 3 || log.Data == nil {
			continue
		}
		//Transfer 事件
		if !bytes.Equal(log.Topics[0][:], erc20CodeParts["LogTransfer"]) {
			continue
		}
		transferLogExist = true
		//跨链金额接收地址
		to = common.BytesToAddress(log.Topics[2][:]).String()
		if !common.IsEqualIgnoreCase(to, checkToAddress) {
			continue
		}
		//0x000000
		from = common.BytesToAddress(log.Topics[1][:]).String()
		value = common.GetBigInt(*log.Data, 0, 32)
		return from, to, value, nil
	}
	if transferLogExist {
		err = tokens.ErrTxWithWrongReceiver
	} else {
		err = tokens.ErrDepositLogNotFound
	}
	return "", "", nil, err
}

func (b *Bridge) checkSwapinInfo(swapInfo *tokens.TxSwapInfo) error {
	if swapInfo.Bind == swapInfo.To {
		return tokens.ErrTxWithWrongSender
	}
	if !tokens.CheckSwapValue(swapInfo, b.IsSrc) {
		return tokens.ErrTxWithWrongValue
	}
	token := b.GetTokenConfig(swapInfo.PairID)
	if token == nil {
		return tokens.ErrUnknownPairID
	}
	bindAddr := swapInfo.Bind
	if !tokens.DstBridge.IsValidAddress(bindAddr) {
		log.Warn("wrong bind address in swapin", "bind", bindAddr)
		return tokens.ErrTxWithWrongMemo
	}
	if params.MustRegisterAccount() && !tools.IsAddressRegistered(bindAddr) {
		return tokens.ErrTxSenderNotRegistered
	}
	if params.IsSwapServer &&
		token.ContractAddress != "" &&
		params.CheckBindAddrIsContract() &&
		common.IsEqualIgnoreCase(swapInfo.TxTo, token.ContractAddress) {
		isContract, err := b.IsContractAddress(bindAddr)
		if err != nil {
			log.Warn("query is contract address failed", "bindAddr", bindAddr, "err", err)
			return err
		}
		if isContract {
			return tokens.ErrBindAddrIsContract
		}
	}
	return nil
}
