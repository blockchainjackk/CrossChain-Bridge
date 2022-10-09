package server

import (
	"context"
	"encoding/hex"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/db"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/db/ddl"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/db/types"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/eth"
	"github.com/anyswap/CrossChain-Bridge/tokens/eth/abicoder"
	types2 "github.com/anyswap/CrossChain-Bridge/types"
	"math/big"
	"time"
)

const (
	RetrySendTxTime = 5
	SwapOutFunHash  = "ad54056d"
)

type BscSender struct {
	bscBridge                 *eth.Bridge
	swapServer                string
	account                   string
	gasDistributeAddress      string
	gasDistributeKey          string
	distributeSwapOutInterval int64
	scanBalanceInterval       int64
	scanBalancePerTimeMill    int64
	db                        *db.CrossChainDB
	tokenConfig               *tokens.TokenConfig
}

func NewBscSender(conf *autoSwapConf) *BscSender {
	bridge := newBscBridge(conf)
	sender := &BscSender{
		bscBridge:                 bridge,
		swapServer:                conf.SwapServer,
		account:                   conf.Account,
		db:                        conf.Db,
		tokenConfig:               conf.TokenPairConfig.DestToken,
		gasDistributeAddress:      conf.GasDistributeAddress,
		gasDistributeKey:          conf.GasDistributeKey,
		distributeSwapOutInterval: conf.DistributeSwapOutInterval,
		scanBalanceInterval:       conf.ScanBalanceInterval,
		scanBalancePerTimeMill:    conf.ScanBalancePerTimeMill,
	}
	return sender
}
func newBscBridge(conf *autoSwapConf) *eth.Bridge {

	bridge := eth.NewCrossChainBridge(false)
	bridge.ChainConfig = &tokens.ChainConfig{
		BlockChain:    conf.BridgeConfig.DestChain.BlockChain,
		NetID:         conf.BridgeConfig.DestChain.NetID,
		Confirmations: conf.BridgeConfig.DestChain.Confirmations,
	}
	bridge.GatewayConfig = &tokens.GatewayConfig{
		APIAddress: conf.BridgeConfig.DestGateway.APIAddress,
	}

	return bridge
}

func (b *BscSender) DistributeGas(gasCh chan string, ctx context.Context) {

	// todo panic
	for {
		select {
		case to := <-gasCh:
			log.Infof("[DistributeGas] start send gas to %v\n", to)
			err := b.sendGas2Address(b.gasDistributeAddress, b.gasDistributeKey, to)
			if err != nil {
				log.Errorf("[DistributeGas] distribute gas fail, from: %v , to: %v , err: %v\n", "", "", err)
			}
			continue
		case <-ctx.Done():
			log.Infof("[DistributeGas] distribute gas stop !\n")
			return
		}
	}

}

//给地址发送gas费用
func (b *BscSender) sendGas2Address(from string, key string, to string) error {

	fee, err := b.calculateGas()
	if err != nil {
		log.Errorf("[sendGas2Address] calculate gas fee fail, from: %v , to: %v , err: %v\n", from, to, err)
		return err
	}

	tx, preHash, err := b.buildSignedGasTx(from, key, to, fee)
	if err != nil {
		log.Errorf("[sendGas2Address] build gas tx fail, from: %v , to: %v , err: %v\n", from, to, err)
		return err
	}
	log.Infof("[sendGas2Address] build gas tx success, from: %v , to: %v \n", from, to)

	var hash string
	for i := 0; i < RetrySendTxTime; i++ {
		hash, err = b.bscBridge.SendTransaction(tx)
		if err != nil {
			log.Errorf("[sendGas2Address] send gas tx fail : %v , from : %v, to : %d, hash :%v\n", err, from, to, preHash)
			time.Sleep(time.Millisecond * 500)
			continue
		}
		if hash != "" && err == nil {
			log.Infof("[sendGas2Address] send gas tx success , from : %v, to : %v, hash : %v\n", from, to, hash)
			//todo  应该将失败的地址检查一下地址余额，然后再发送的channel中重新发送
			GasCh <- to
			break
		}

	}
	if err != nil {
		log.Errorf("[sendGas2Address] retry finished,send gas tx fail. Hash : %v ,from : %v , to : %v, err: %v\n",
			hash, from, to, err)
		return err
	}

	return nil
}
func (b *BscSender) calculateGas() (*big.Int, error) {
	//DefaultGasLimit 60000
	limit := b.tokenConfig.DefaultGasLimit
	suggestPrice, err := b.bscBridge.SuggestPrice()
	if err != nil {
		log.Errorf("[calculateGas] get suggest gas price err %v ", err)
		return nil, err
	}
	gasFee := new(big.Int).Mul(suggestPrice, new(big.Int).SetUint64(limit))

	return gasFee, nil
}

func (b *BscSender) buildSignedGasTx(from string, key string, to string, gas *big.Int) (interface{}, string, error) {
	// b.buildTx
	args := b.buildTxArgs(from, to, gas, nil)
	tx, err := b.buildTx(args)
	if err != nil {
		return nil, "", err
	}
	privateKey, err := HexToPrivateKey(key)
	if err != nil {
		log.Errorf("[buildGasSignedTx] hex to private key err: %v\n", err)
		return nil, "", err
	}
	id, err := b.bscBridge.GetSignerChainID()
	if err != nil {
		log.Errorf("[buildGasSignedTx] get signer chain id err: %v\n", err)
		return nil, "", err
	}
	b.bscBridge.MakeSigner(id)
	withPrivateKey, hash, err := b.bscBridge.SignTransactionWithPrivateKey(tx, privateKey)

	if err != nil {
		log.Errorf("[buildGasSignedTx] sign tx err: %v\n", err)
		return nil, "", err
	}
	return withPrivateKey, hash, nil
}

type SwapOutInfo struct {
	tx   interface{}
	from string
	bind string
}

func (b *BscSender) DistributeSwapOut(ctx context.Context) {
	//todo  panic

	err := b.db.CreateTable(ddl.SwapOutTaleName, ddl.CreateSwapOutTable)
	if err != nil {
		log.Errorf("[DistributeSwapOut] create %v fail %v\n", ddl.SwapOutTaleName, err)
		return
	}

	SwapOutTxCh := make(chan interface{})
	SwapOutTxChIsClose := make(chan interface{})
	go b.SwapOutTxProcessor(ctx, SwapOutTxCh, SwapOutTxChIsClose)

	for {
		//总是重复抓取，这里应该根据余额获取，并放慢节奏,
		//内部逻辑，同一个地址如果钱不够了就不会再拿出来了
		fromInfo, binds, err := b.findAddrToSwapOut1()
		if err != nil {
			log.Errorf("[DistributeSwapOut] find addr to swap out err %v\n", err)
			time.Sleep(time.Second * 10)
			continue
		}
		log.Infof("[DistributeSwapOut] find  %v addr to distribute swap out.\n", len(fromInfo))
		//signedTxs := make([]interface{}, 0)
		for i, info := range fromInfo {
			tx, _, err := b.buildSignedSwapOutTx(info.Address, binds[i], info.Key, info.Balance)
			if err != nil {
				log.Errorf("[DistributeSwapOut] build signed swapOut tx err %v \n", err)
				continue
			}
			select {
			case <-SwapOutTxChIsClose:
				log.Infof("[DistributeSwapOut] SwapOutTxCh is close.\n")
				return
			default:
				swapInfo := &SwapOutInfo{
					tx:   tx,
					from: info.Address,
					bind: binds[i],
				}
				log.Infof("[DistributeSwapOut] distribute swap out from : %v, bind : %v.\n", swapInfo.from, swapInfo.bind)
				SwapOutTxCh <- swapInfo
			}
			time.Sleep(time.Second * time.Duration(b.distributeSwapOutInterval))
		}
	}
}

func (b *BscSender) SwapOutTxProcessor(ctx context.Context, SwapOutTxCh, isCloseChan chan interface{}) error {
	for {
		select {
		case <-ctx.Done():
			isCloseChan <- struct{}{}
			time.Sleep(time.Second * 3)
			close(SwapOutTxCh)
			close(isCloseChan)
			log.Infof("[SwapOutTxProcessor] stop !\n")
			return nil
		case swapOutInfo := <-SwapOutTxCh:
			info, _ := swapOutInfo.(*SwapOutInfo)
			tx1, _ := info.tx.(*types2.Transaction)
			log.Infof("[SwapOutTxProcessor] recive  %v to send!\n", tx1.Hash().String())
			_ = b.sendAndSaveSwapOutTx(info)
			continue
		}
	}

}

func (b *BscSender) sendAndSaveSwapOutTx(swapOutInfo *SwapOutInfo) error {

	tx, _ := swapOutInfo.tx.(*types2.Transaction)

	var hash string
	var err error
	for i := 0; i < RetrySendTxTime; i++ {
		hash, err = b.bscBridge.SendTransaction(tx)
		if err != nil {
			log.Errorf("[sendSwapOutTx] send swap out tx err %v , hash :%v\n", err, tx.Hash().String())
			continue
		}
		if hash != "" && err == nil {
			log.Infof("[sendSwapOutTx]send swap out tx success. hash : %v\n", hash)
			break
		}

	}
	if err != nil {
		log.Errorf("[sendSwapOutTx] retry finished,send swap out tx fail. Hash : %v ,from : %v , bind : %v, err: %v\n",
			hash, swapOutInfo.from, swapOutInfo.bind, err)
		return err
	}
	//db 落库
	err = b.db.InsertTxInSwapOut(hash, swapOutInfo.from, swapOutInfo.bind, int64(TxNotSwapped))
	if err != nil {
		log.Errorf("[sendSwapOutTx] save tx fail. Hash : %v ,from : %v , bind : %v, err: %v\n",
			hash, swapOutInfo.from, swapOutInfo.bind, err)

		return err
	}
	log.Infof("[sendSwapOutTx] save tx success. Hash : %v ,from : %v , bind : %v\n",
		hash, swapOutInfo.from, swapOutInfo.bind)

	return nil
}

func (b *BscSender) buildSwapOutTxInput(bind string, balance int64) []byte {

	//MinimumSwap := b.tokenConfig.MinimumSwap
	funcHash, _ := hex.DecodeString(SwapOutFunHash)

	//var value float64
	//_, v := RandomNormalInt64(1, 10, 2, 1)
	//balanceD := toDcrnCoin(balance)
	//if float64(v)-balanceD > *MinimumSwap {
	//	value = float64(v)
	//} else {
	//	value = balanceD
	//}
	//swapValueBig := tokens.ToBits(value, *b.tokenConfig.Decimals)

	input := abicoder.PackDataWithFuncHash(funcHash, big.NewInt(balance), bind)
	return input

}

func (b *BscSender) buildSignedSwapOutTx(from, bind, key string, balance int64) (rawTx interface{}, txHash string, err error) {

	input := b.buildSwapOutTxInput(bind, balance)

	args := b.buildTxArgs(from, b.tokenConfig.ContractAddress, big.NewInt(0), &input)
	tx, err := b.buildTx(args)
	if err != nil {
		return nil, "", err
	}
	privateKey, err := HexToPrivateKey(key)
	if err != nil {
		log.Errorf("[buildSignedSwapOutTx] hex to private key err: %v\n", err)
		return nil, "", err
	}
	id, err := b.bscBridge.GetSignerChainID()
	if err != nil {
		log.Errorf("[buildSignedSwapOutTx] get signer chain id err: %v\n", err)
		return nil, "", err
	}
	b.bscBridge.MakeSigner(id)
	withPrivateKey, hash, err := b.bscBridge.SignTransactionWithPrivateKey(tx, privateKey)

	if err != nil {
		log.Errorf("[buildSignedSwapOutTx] sign tx err: %v\n", err)
		return nil, "", err
	}
	return withPrivateKey, hash, nil

}

func (b *BscSender) buildTx(args *tokens.BuildTxArgs) (rawTx interface{}, err error) {

	return b.bscBridge.BuildRawTransaction(args)
}

func (b *BscSender) buildTxArgs(from, to string,
	Value *big.Int, Input *[]byte) *tokens.BuildTxArgs {

	args := &tokens.BuildTxArgs{
		From:  from,
		To:    to,
		Value: Value,
		Input: Input,
	}
	return args
}

// 找到满足调用SwapOUT 的地址，使用call contract
// 找到DCRN的receviver 可以是原来的from

func (b *BscSender) findAddrToSwapOut() ([]*types.AddressInfo, []string, error) {
	//to作为swapOut的sender、from作为bind
	from, to, err := b.db.RetrieveToAddressFromSwapIn()
	if err != nil {
		log.Errorf("RetrieveToAddressFromSwapIn err %v\n", err)
		return nil, nil, err
	}

	//todo 测试的时候，可以将from写死然后走下面的逻辑
	//from := []string{"Sscaawq6q8HsZ8xVA9cFMAVQxez79VF6Qn5"}
	//to := []string{"0x1b499a530A14c5A06316bB8b551D2d89d2b811b4"}
	toAddressInfo := make([]*types.AddressInfo, 0)
	minSwap := tokens.ToBits(*b.tokenConfig.MinimumSwap, *b.tokenConfig.Decimals)
	binds := make([]string, 0)
	for index, addr := range to {
		//获取余额的方式，到后面这里可以直接查库
		balance, err := b.bscBridge.GetErc20Balance(b.tokenConfig.ContractAddress, addr)
		if err != nil {
			log.Errorf("get erc2o balance fail, addr : %v\n", addr)
			continue
		}

		//todo 这里也需要查询bnb gas费用是否充足

		if balance.Cmp(minSwap) < 0 {
			continue
		}
		key, err := b.db.RetrieveKeyByAddress(addr)
		if err != nil {
			log.Errorf("find private key  fail, addr : %v\n", addr)
			continue
		}
		toInfo := &types.AddressInfo{
			Address: addr,
			Key:     key,
		}
		toAddressInfo = append(toAddressInfo, toInfo)
		binds = append(binds, from[index])
	}
	return toAddressInfo, binds, nil
}

// 每隔一分钟触发一次余额更新的任务
func (b *BscSender) BalanceScan(ctx context.Context) {

	// todo panic
	for {

		//找到地址
		addrs, err := b.db.RetrieveAddressFromSwapIn2()
		if err != nil {
			log.Errorf("[BalanceScan] get  addrs from swapin err %v\n", err)
			return
		}
		for _, addr := range addrs {
			balance, err := b.bscBridge.GetErc20Balance(b.tokenConfig.ContractAddress, addr)
			if err != nil {
				log.Errorf("[BalanceScan] get balance fail, address : %v, %v \n", addr, err)
				continue
			}
			err = b.db.UpdateAddrBalance(balance.Int64(), addr)
			if err != nil {
				log.Errorf("[BalanceScan] update balance fail, address : %v, %v \n", addr, err)
				continue
			}
			fbalance := toDcrnCoin(balance.Int64())
			log.Debugf("[BalanceScan] update balance success, address: %v, balance: %v \n", addr, fbalance)
			time.Sleep(time.Millisecond * time.Duration(b.scanBalancePerTimeMill))
		}
		log.Info("[BalanceScan] update balance finished")
		time.Sleep(time.Second * time.Duration(b.scanBalanceInterval))
	}

}

func (b *BscSender) findAddrToSwapOut1() ([]*types.AddressInfo, []string, error) {
	//to作为swapOut的sender、from作为bind
	minSwap := tokens.ToBits(*b.tokenConfig.MinimumSwap, *b.tokenConfig.Decimals)

	AddrInfo, err := b.db.RetrieveAddressToSwapOut(minSwap.Int64())
	if err != nil {
		log.Errorf("RetrieveToAddressFromSwapIn err %v\n", err)
		return nil, nil, err
	}

	//todo 测试的时候，可以将from写死然后走下面的逻辑
	//from := []string{"Sscaawq6q8HsZ8xVA9cFMAVQxez79VF6Qn5"}
	//to := []string{"0x1b499a530A14c5A06316bB8b551D2d89d2b811b4"}
	toAddressInfo := make([]*types.AddressInfo, 0)
	binds := make([]string, 0)
	for _, addr := range AddrInfo {
		//获取余额的方式，到后面这里可以直接查库
		balanceDcrn, err := b.bscBridge.GetErc20Balance(b.tokenConfig.ContractAddress, addr.Address)
		if err != nil {
			log.Errorf("get erc2o balance fail, addr : %v\n", addr)
			continue
		}

		if balanceDcrn.Int64() != addr.Balance {

			err := b.db.UpdateAddrBalance(balanceDcrn.Int64(), addr.Address)
			if err != nil {
				log.Errorf("update balance err %v\n", err)
				continue
			}
			if addr.Balance < balanceDcrn.Int64() {
				continue
			}
		}

		//查询gas费是否够用
		balanceBnb, err := b.bscBridge.GetBalance(addr.Address)
		if err != nil {
			log.Errorf("get bnb balance fail, addr : %v\n", addr)
			continue
		}
		gas, _ := b.calculateGas()
		if gas == nil {
			continue
		}
		if balanceBnb.Cmp(gas) < 0 {
			GasCh <- addr.Address
			continue
		}

		if balanceDcrn.Cmp(minSwap) < 0 {
			continue
		}
		//bind
		from, err := b.db.RetrieveBindByAddress(addr.Address)
		if err != nil {
			log.Errorf("find private key  fail, addr : %v\n", addr)
			continue
		}
		toInfo := &types.AddressInfo{
			Address: addr.Address,
			Key:     addr.Key,
			Balance: balanceDcrn.Int64(),
		}
		toAddressInfo = append(toAddressInfo, toInfo)
		binds = append(binds, from)
	}
	return toAddressInfo, binds, nil
}

func (b *BscSender) findTxs2SwapOut() ([]string, error) {

	txs, err := b.db.RetrieveTxsToSwapOut()

	return txs, err
}
func (b *BscSender) updateSwapOutTxStatus(tx string, status int64) error {
	return b.db.UpdateSwapOutStatus(tx, status)
}
