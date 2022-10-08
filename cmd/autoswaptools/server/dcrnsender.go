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
	dcrnBridge  *dcrn.Bridge
	swapServer  string
	account     string
	db          *db.CrossChainDB
	tokenConfig *tokens.TokenConfig
}

const CrossChainDepositAddress = "SccW2X8WJBkPep9fqx7G1ssyBMRaX5w3rjz"
const DistributeDcrnInterval = 30 * time.Second

func NewDcrnSender(conf *autoSwapConf) *DcrnSender {
	bridge := newDcrnBridge(conf)
	sender := &DcrnSender{
		dcrnBridge:  bridge,
		swapServer:  conf.SwapServer,
		account:     conf.Account,
		db:          conf.Db,
		tokenConfig: conf.TokenPairConfig.SrcToken,
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

// todo 定时分发

func (d *DcrnSender) DistributeDcrn(ctx context.Context) error {

	err := d.db.CreateTable(ddl.SwapInTaleName, ddl.CreateSwapInTable)
	if err != nil {
		log.Errorf("[DistributeDcrn] create swapin table fail: %v\n", err)
		return err
	}
	timer := time.NewTicker(DistributeDcrnInterval)
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
	// 0.3 + 加预估0.1的手续费
	if balance < 0.4 {
		//todo 这里是不是最好可以做点通知
		return fmt.Errorf("account %s balance  not enough.", d.account)
	}
	hash, err := rpc.SendFromAccount(d.dcrnBridge, d.account, CrossChainDepositAddress, 0.3)
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
	log.Info("[DistributeDcrn] SendFromAccount success:", "hash", hash)
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
