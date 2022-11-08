package server

import (
	"fmt"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/server/crsclient"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/server/rpc"
	"github.com/anyswap/CrossChain-Bridge/log"
	"io/ioutil"
	"time"
)

const swapInApiParams = `/swapindcrn/post`

var GasCh = make(chan string)

type SwapinWorker struct {
	swapInPerTxInterval int64
	swapInInterval      int64
	dcrnsender          *DcrnSender
}

type swapInRequest struct {
	TxId        string `json:"txID"`
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	SignInfo    string `json:"signMsg"`
}

func NewSwapInWorker(d *DcrnSender, conf *autoSwapConf) *SwapinWorker {

	return &SwapinWorker{
		dcrnsender:          d,
		swapInInterval:      conf.SwapInInterval,
		swapInPerTxInterval: conf.SwapInPerTxInterval,
	}
}

// DoSwapInWork 定时进行跨链转账
func (s *SwapinWorker) DoSwapInWork() {
	//todo recovery  err  panic

	for {

		//1、找到可以去跨链转账的交易
		txs, fromAddrs, err := s.findTx2SwapIn()
		log.Infof("[DoSwapInWork] find %v tx send swapin request\n", len(txs))
		if err != nil {
			log.Error("[DoSwapInWork] finsd tx to swapin fail : ", "err", err)
			return
		}
		txAmount := len(txs)
		if txAmount == 0 {
			log.Infof("[DoSwapInWork] find no tx to  send swapin request\n")

			time.Sleep(time.Second * time.Duration(s.swapInInterval))
			continue
		}

		//2、查找BSC端的toAddress
		toAddrs, err := s.dcrnsender.findBscAddrToSwapIn(int64(txAmount))
		if err != nil {
			log.Error("[DoSwapInWork] find bsc addr to swapin fail : ", "err", err)
			return
		}
		if len(toAddrs) == 0 {
			time.Sleep(time.Second * time.Duration(s.swapInInterval))
			continue
		}
		if txAmount > len(toAddrs) {
			txAmount = len(toAddrs)
		}
		for i := 0; i < txAmount; i++ {

			//3、构建跨链转账参数
			tx := txs[i]
			fromAddress := fromAddrs[i]
			toAddress := toAddrs[i]
			msg, err := rpc.GetSignMsg(s.dcrnsender.dcrnBridge, fromAddress, toAddress.Address)
			if err != nil {
				err := fmt.Errorf("get sign msg fail, txid:%v , %v\v", tx, err)
				log.Error("[DoSwapInWork]", "err", err)
				continue
			}
			request := &swapInRequest{
				TxId:        tx,
				FromAddress: fromAddress,
				ToAddress:   toAddress.Address,
				SignInfo:    msg,
			}

			//4、调用跨链服务
			swapInUrl := s.dcrnsender.swapServer + swapInApiParams
			resp, err := crsclient.HTTPPost(swapInUrl, request, nil, nil, crsclient.DefaultTimeout)
			if err != nil {
				err := fmt.Errorf("send cross chain tx fail, txid:%v , %v\v", tx, err)
				log.Error("[DoSwapInWork]", "err", err)
				continue
			}

			//5、更改跨链交易状态
			body, err := ioutil.ReadAll(resp.Body)
			str := string(body)

			if str == `"Success"` {
				err := s.dcrnsender.updateSwapIn(tx, fromAddress, toAddress.Address, msg, int64(TxProcessed))
				if err != nil {

					err := fmt.Errorf("save swapin to db fail, txid:%v , %v\v", tx, err)
					log.Error("[DoSwapInWork]", "err", err)
					return
				}

				log.Infof("[DoSwapInWork] swapin success ! tx : %v, fromAddress : %v, toAddress : %v, msg : %v\n", tx, fromAddress, toAddress.Address, msg)
			} else {
				log.Warnf("[DoSwapInWork] swapin  fail: %v ,tx : %v, fromAddress : %v, toAddress : %v, msg : %v\n", str, tx, fromAddress, toAddress.Address, msg)
			}
			time.Sleep(time.Second * time.Duration(s.swapInPerTxInterval))
		}

	}
}

func (s *SwapinWorker) findTx2SwapIn() ([]string, []string, error) {

	txs, err := s.dcrnsender.db.RetrieveTx2SwapIn()
	if err != nil {
		log.Error("[findTx2SwapIn] fail :", "err", err)
		if err != nil {

		}
	}
	txIds := make([]string, 0)
	fromAddress := make([]string, 0)
	confirmations := s.dcrnsender.dcrnBridge.ChainConfig.Confirmations
	for _, tx := range txs {
		//查询状态
		result, err := rpc.GetDcrnTransactionByHash(s.dcrnsender.dcrnBridge, tx)
		if err != nil {
			log.Error("[findTx2SwapIn] get tx status fail :", "err", err)
			continue
		}

		if uint64(result.Confirmations) < *confirmations {
			continue
		}

		address, err := s.getFromAddress(result.Vin[0].Txid, int64(result.Vin[0].Vout))
		if err != nil {
			log.Error("[findTx2SwapIn] get from address fail :", "err", err)
			continue
		}
		txIds = append(txIds, tx)
		fromAddress = append(fromAddress, address)
	}
	return txIds, fromAddress, nil
}

func (s *SwapinWorker) getFromAddress(txId string, vout int64) (string, error) {
	result, err := rpc.GetDcrnTransactionByHash(s.dcrnsender.dcrnBridge, txId)
	if err != nil {
		return "", err
	}
	amount := len(result.Vout)
	if int64(amount) <= vout {
		err := fmt.Errorf("outputs of %s err\n", txId)
		log.Errorf("[getFromAddress] get from address fail :", "err", err)
		return "", err
	}
	address := result.Vout[vout].ScriptPubKey.Addresses[0]

	return address, err
}
