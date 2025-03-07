package worker

import (
	"time"

	"github.com/anyswap/CrossChain-Bridge/params"
	"github.com/anyswap/CrossChain-Bridge/rpc/client"
	"github.com/anyswap/CrossChain-Bridge/tokens/bridge"
)

const interval = 10 * time.Millisecond

// StartWork start swap server work
func StartWork(isServer bool) {
	if isServer {
		logWorker("worker", "start server worker")
	} else {
		logWorker("worker", "start oracle worker")
	}

	client.InitHTTPClient()
	bridge.InitCrossChainBridge(isServer)

	if params.IsTestMode() {
		if isServer {
			StartTestWork()
		} else {
			StartAcceptSignJob()
		}
		return
	}

	StartScanJob(isServer)
	time.Sleep(interval)

	StartUpdateLatestBlockHeightJob()
	time.Sleep(interval)

	if !isServer {
		StartAcceptSignJob()
		time.Sleep(interval)
		AddTokenPairDynamically()
		time.Sleep(interval)
		StartReportStatJob()
		return
	}

	StartSwapJob()
	time.Sleep(interval)

	StartVerifyJob()
	time.Sleep(interval)

	StartStableJob()
	time.Sleep(interval)

	StartReplaceJob()
	time.Sleep(interval)
	//如果配置文件EnablePassBigValue打开，会对bigValue交易做验证，从而生成swapresult表放到数据库中，但是bigValue 要12小时后才做这个操作。
	StartPassBigValueJob()
	time.Sleep(interval)
	//针对于BTC的
	StartAggregateJob()
	time.Sleep(interval)

	StartCheckFailedSwapJob()
}
