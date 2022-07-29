package eth

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/anyswap/CrossChain-Bridge/common"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/params"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/eth/abicoder"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

const VerfyMultiSigPath = "../../tokens/eth/abicoder/abi/VerfyMultiSig.abi"

var (
	errInvalidReceiverAddress = errors.New("invalid receiver address")
)

func (b *Bridge) buildSwapoutTxInput(args *tokens.BuildTxArgs) (err error) {
	token := b.GetTokenConfig(args.PairID)
	if token == nil {
		return tokens.ErrUnknownPairID
	}

	receiver := common.HexToAddress(args.Bind)
	if receiver == (common.Address{}) || !common.IsHexAddress(args.Bind) {
		log.Warn("swapout to wrong address", "receiver", args.Bind)
		return errInvalidReceiverAddress
	}

	swapValue := tokens.CalcSwappedValue(args.PairID, args.OriginValue, false, args.OriginFrom, args.OriginTxTo)
	swapValue, err = b.adjustSwapValue(args, swapValue)
	if err != nil {
		return err
	}
	args.SwapValue = swapValue // swap value

	if token.ContractAddress == "" {
		input := b.getUnlockCoinMemo(args)
		args.Input = &input    // input
		args.To = args.Bind    // to
		args.Value = swapValue // value
		return nil
	}

	if token.MultiSignContractDepositAddress == "" {
		funcHash := erc20CodeParts["transfer"]
		input := abicoder.PackDataWithFuncHash(funcHash, receiver, swapValue)
		args.Input = &input             // input
		args.To = token.ContractAddress // to
		return b.checkBalance(token.ContractAddress, token.DcrmAddress, swapValue)
	} else {
		//新增的使用多签方式进行withdraw
		erc20contract := common.HexToAddress(token.ContractAddress)
		signs, err := multiSigns(MultiSignParam{
			Erc20Contract: erc20contract,
			Receiver:      receiver,
			Value:         *swapValue,
		})
		if err != nil {
			log.Errorf("multiSigns err", err)
			return err
		}
		rs, ss, vs := sigs2rsv(signs)

		abiJson, err := ioutil.ReadFile(VerfyMultiSigPath)
		if err != nil {
			log.Errorf("read abi file err", err)
			return err
		}

		myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
		if err != nil {
			log.Errorf("abi json err", err)
			return err
		}
		input, _ := myAbi.Pack("spendERC20", receiver, erc20contract, swapValue, vs, rs, ss)
		args.Input = &input // input
		//调用（提钱）的是多签合约（钱包），上面的ContractAddress地址是ERC20代币合约的地址
		args.To = token.MultiSignContractDepositAddress
		return b.checkBalance(token.ContractAddress, token.DcrmAddress, swapValue)
	}
}

func (b *Bridge) getUnlockCoinMemo(args *tokens.BuildTxArgs) (input []byte) {
	if params.IsNullSwapoutNativeMemo() {
		return input
	}
	isContract, err := b.IsContractAddress(args.Bind)
	if err == nil && !isContract {
		input = []byte(tokens.UnlockMemoPrefix + args.SwapID)
	}
	return input
}
