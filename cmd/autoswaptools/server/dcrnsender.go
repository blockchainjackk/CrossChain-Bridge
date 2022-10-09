package server

import (
	"context"
	"fmt"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/db"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/db/ddl"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/db/types"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/server/rpc"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/anyswap/CrossChain-Bridge/tokens/dcrn"
	"time"
)

type DcrnSender struct {
	dcrnBridge             *dcrn.Bridge
	swapServer             string
	account                string
	db                     *db.CrossChainDB
	tokenConfig            *tokens.TokenConfig
	distributeDcrnInterval int64
}

func NewDcrnSender(conf *autoSwapConf) *DcrnSender {
	bridge := newDcrnBridge(conf)
	sender := &DcrnSender{
		dcrnBridge:             bridge,
		swapServer:             conf.SwapServer,
		account:                conf.Account,
		db:                     conf.Db,
		tokenConfig:            conf.TokenPairConfig.SrcToken,
		distributeDcrnInterval: conf.DistributeDcrnInterval,
	}
	return sender
}

func newDcrnBridge(conf *autoSwapConf) *dcrn.Bridge {

	bridge := dcrn.NewCrossChainBridge(true)
	bridge.ChainConfig = &tokens.ChainConfig{
		BlockChain:    conf.BridgeConfig.SrcChain.BlockChain,
		NetID:         conf.BridgeConfig.SrcChain.NetID,
		Confirmations: conf.BridgeConfig.SrcChain.Confirmations,
	}
	bridge.GatewayConfig = &tokens.GatewayConfig{
		APIAddress: conf.BridgeConfig.SrcGateway.APIAddress,
	}

	return bridge
}

func (d *DcrnSender) DistributeDcrn(ctx context.Context) error {

	err := d.db.CreateTable(ddl.SwapInTaleName, ddl.CreateSwapInTable)
	if err != nil {
		log.Errorf("[DistributeDcrn] create swapin table fail: %v\n", err)
		return err
	}
	timer := time.NewTicker(time.Second * time.Duration(d.distributeDcrnInterval))
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Infof("[DistributeDcrn] distributeDcrn stop !\n")
			return nil
		case <-timer.C:
			err = d.SendDcrnToDepositAddress()
			if err != nil {
				log.Errorf("[DistributeDcrn] send dcrn to deposit address fail: %v\n", err)
				//return err
				continue
			}
		}

	}
}

// 向跨链地址打钱，生成跨链转账txHash
func (d *DcrnSender) SendDcrnToDepositAddress() error {
	balance, err := rpc.GetAccountBalance(d.dcrnBridge, d.account)
	if err != nil {
		log.Error("[DistributeDcrn] GetAccountBalance fail :", "err", err)
		return err
	}
	//生成转账费用
	var value float64
	_, v := RandomNormalInt64(1, 10, 2, 1)
	if float64(v) >= *d.tokenConfig.MinimumSwap {
		value = float64(v)
	} else {
		value = *d.tokenConfig.MinimumSwap + *d.tokenConfig.MinimumSwapFee*2
	}

	if balance < value {
		//todo 这里是不是最好可以做点通知
		return fmt.Errorf("account %s balance  not enough. balance: %v, need: %v\n", d.account, balance, value)
	}
	hash, err := rpc.SendFromAccount(d.dcrnBridge, d.account, d.tokenConfig.DepositAddress, value)
	if err != nil {
		log.Error("[DistributeDcrn]  SendFromAccount fail:", "err", err)
		return err
	}
	// 入库
	err = d.db.InsertTxInSwapIn(hash, int64(TxNotSwapped))
	if err != nil {
		log.Error("[DistributeDcrn]  SendFromAccount fail:", "err", err)
		return err
	}
	log.Infof("[DistributeDcrn] SendFromAccount success: Hash :%v ,value %v\n", hash, value)
	return nil
}

func (d *DcrnSender) findBscAddrToSwapIn(account int64) ([]*types.AddressInfo, error) {

	txs, err := d.db.RetrieveAddressesToSwapIn(account)

	return txs, err
}

func (d *DcrnSender) updateSwapIn(
	txId string, fromAddress string, toAddress string, signInfo string, status int64) error {

	err := d.db.UpdateSwapIn(txId, fromAddress, toAddress, signInfo, status)

	return err
}
