package server

import (
	"context"
	"fmt"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/server/crsclient"
	"github.com/anyswap/CrossChain-Bridge/log"
	"io/ioutil"
	"time"
)

// 这个需要时间长一点应为要等待DCRN主链出块
const DefaultSwapOutInterval = time.Second * 20
const swapOutApiParams = `/swapout/post/DCRN/`

type SwapOutWorker struct {
	bscSender *BscSender
}

func NewSwapOutWorker(b *BscSender) *SwapOutWorker {

	return &SwapOutWorker{
		bscSender: b,
	}
}

func (s *SwapOutWorker) DoSwapOutWork(ctx context.Context) {

	for true {

		//1、找到可以去跨链转账的交易
		txs, err := s.bscSender.findTxs2SwapOut()
		if err != nil {
			log.Error("[DoSwapOutWork] find txs to swap out fail : ", "err", err)
			return
		}
		txAmount := len(txs)
		if txAmount == 0 {
			time.Sleep(time.Second * 10)
			continue
		}

		for _, tx := range txs {

			//todo 判断交易是否成熟

			//2、发送跨链交易
			swapOutUrl := s.bscSender.swapServer + swapOutApiParams + tx
			resp, err := crsclient.HTTPPost(swapOutUrl, nil, nil, nil, crsclient.DefaultTimeout)
			if err != nil {
				err := fmt.Errorf("send cross chain tx fail, txid:%v , %v\v", tx, err)
				log.Error("[DoSwapOutWork]", "err", err)
				continue
			}

			//3、更改跨链交易状态
			body, err := ioutil.ReadAll(resp.Body)
			str := string(body)

			if str == `"Success"` {
				err := s.bscSender.updateSwapOutTxStatus(tx, int64(TxProcessed))
				if err != nil {
					err := fmt.Errorf("save swapout to db fail, txid:%v , %v\v", tx, err)
					log.Error("[DoSwapOutWork]", "err", err)
					return
				}
				log.Infof("[DoSwapOutWork] swapout success, hash: %v\n", tx)
			} else {
				log.Warnf("[DoSwapOutWork] swapout  fail: %v, hash: %v\n", str, tx)
			}
			time.Sleep(DefaultSwapOutInterval)
		}

	}
}
