package server

import "github.com/urfave/cli/v2"

type SwapStatus uint16

// swap status values
const (
	TxNotSwapped SwapStatus = iota // 0  区块的确认数还没有达到跨链要求
	TxProcessed                    // 1  向swapserver成功发送了跨链请求
	TxSwapFailed                   // 2  向swapserver发送跨链请求失败了
)

var (
	swapServerFlag = &cli.StringFlag{
		Name:  "swapserver",
		Usage: "swap server, ie. http://13.212.22.91:9545",
		Value: "http://127.0.0.1:11556",
	}
	dcrnNetworkFlag = &cli.StringFlag{
		Name:  "dcrnnet",
		Usage: "network identifier, ie. mainnet, testnet3 ,simnet",
		Value: "simnet",
	}
	dcrnGateWayFlag = &cli.StringFlag{
		Name:  "gateway",
		Usage: "gateway",
		Value: "https://127.0.0.1:19557",
	}
	dcrnAccountFlag = &cli.StringFlag{
		Name:  "account",
		Usage: "dcrn wallet account for distribution",
		Value: "boom",
	}
	bridgeConfigFileFlag = &cli.StringFlag{
		Name:  "bridgeConfigFile",
		Usage: "bridge swapserver config file",
		Value: "/Users/boom/bitmain/projects/crosschain/antpool/CrossChain-Bridge/test-conf/dcrn/config.toml",
	}
	tokenConfigDirFlag = &cli.StringFlag{
		Name:  "tokenConfigDir",
		Usage: "bridge swapserver token config dir",
		Value: "/Users/boom/bitmain/projects/crosschain/antpool/CrossChain-Bridge/test-conf/dcrn/tokenpairs",
	}
	configFileFlag = &cli.StringFlag{
		Name:  "configFile",
		Usage: "dcrn auto swap config file",
		Value: "./config_main.toml",
	}
)
