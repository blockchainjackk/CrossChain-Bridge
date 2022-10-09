package main

import (
	"context"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/server"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/server/crsclient"
	"github.com/anyswap/CrossChain-Bridge/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	server.InitApp()
	if err := server.App.Run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	ctx, cancel := context.WithCancel(context.Background())
	//todo
	server.BscAccountInit(server.AutoSwapConf.Db)
	crsclient.InitHTTPClient()
	sender := server.NewDcrnSender(server.AutoSwapConf)

	go sender.DistributeDcrn(ctx)

	swapinWork := server.NewSwapInWorker(sender, server.AutoSwapConf)

	go swapinWork.DoSwapInWork()

	bscSender := server.NewBscSender(server.AutoSwapConf)
	swapOutWorker := server.NewSwapOutWorker(bscSender, server.AutoSwapConf)

	go bscSender.BalanceScan(ctx)
	// gorutine
	go bscSender.DistributeSwapOut(ctx)

	go bscSender.DistributeGas(server.GasCh, ctx)

	go swapOutWorker.DoSwapOutWork(ctx)

	select {}

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	cancel()

}
