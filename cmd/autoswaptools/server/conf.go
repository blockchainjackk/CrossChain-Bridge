package server

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/anyswap/CrossChain-Bridge/cmd/autoswaptools/db"
	"github.com/anyswap/CrossChain-Bridge/cmd/utils"
	"github.com/anyswap/CrossChain-Bridge/common"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/params"
	"github.com/anyswap/CrossChain-Bridge/tokens"
	"github.com/urfave/cli/v2"
	"sync"
)

var (
	clientIdentifier = "swaptools"
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	gitDate   = ""
	// The App that holds all commands and flags.
	App               = utils.NewApp(clientIdentifier, gitCommit, gitDate, "the swaptools command line interface")
	loadConfigStarter sync.Once
)

const PairID = "dcrn"

type autoSwapConf struct {
	SwapServer       string
	Account          string
	BridgeConfigFile string
	TokenConfigFile  string
	BridgeConfig     *params.BridgeConfig
	Db               *db.CrossChainDB
	TokenPairConfig  *tokens.TokenPairConfig
}

var AutoSwapConf *autoSwapConf

func LoadConfig(configFile string) *autoSwapConf {
	loadConfigStarter.Do(func() {
		if configFile == "" {
			log.Fatalf("LoadConfig error: no config file specified")
		}
		log.Println("Config file is", configFile)
		if !common.FileExist(configFile) {
			log.Fatalf("LoadConfig error: config file %v not exist", configFile)
		}
		config := &autoSwapConf{}
		if _, err := toml.DecodeFile(configFile, &config); err != nil {
			log.Fatalf("LoadConfig error (toml DecodeFile): %v", err)
		}
		bridgeConfig := &params.BridgeConfig{}
		if config.BridgeConfigFile != "" {
			if _, err := toml.DecodeFile(config.BridgeConfigFile, &bridgeConfig); err != nil {
				log.Fatalf("LoadConfig error (toml DecodeFile): %v", err)
			}
		}
		config.BridgeConfig = bridgeConfig

		tokenConfig := &tokens.TokenPairConfig{}
		if config.TokenConfigFile != "" {
			if _, err := toml.DecodeFile(config.TokenConfigFile, &tokenConfig); err != nil {
				log.Fatalf("LoadConfig error (toml DecodeFile): %v", err)
			}
		}
		config.TokenPairConfig = tokenConfig
		AutoSwapConf = config
		//为了初始化rpc服务
		params.SetConfig(config.BridgeConfig)
	})
	return AutoSwapConf
}

func InitApp() {
	// Initialize the CLI app and start action
	App.Action = autoSwapDcrn
	App.HideVersion = true // we have a command to print the version
	App.Copyright = "Copyright 2020-2022 The CrossChain-Bridge Authors"
	App.Commands = []*cli.Command{
		utils.LicenseCommand,
		utils.VersionCommand,
	}
	App.Flags = []cli.Flag{
		utils.VerbosityFlag,
		utils.JSONFormatFlag,
		utils.ColorFormatFlag,
		//swapServerFlag,
		//dcrnNetworkFlag,
		//dcrnGateWayFlag,
		//dcrnAccountFlag,
		//bridgeConfigFileFlag,
		//tokenConfigDirFlag,
		configFileFlag,
	}
}

func autoSwapDcrn(ctx *cli.Context) error {

	utils.SetLogger(ctx)
	if ctx.NArg() > 0 {
		return fmt.Errorf("invalid command: %q", ctx.Args().Get(0))
	}
	_ = cli.ShowAppHelp(ctx)
	fmt.Println()
	//swapServer := ctx.String(swapServerFlag.Name)
	//dcrnNetwork := ctx.String(dcrnNetworkFlag.Name)
	//dcrnGateWay := ctx.String(dcrnGateWayFlag.Name)
	//account := ctx.String(dcrnAccountFlag.Name)
	//bridgeconfigFile := ctx.String(bridgeConfigFileFlag.Name)

	//tokenConfigDir := ctx.String(tokenConfigDirFlag.Name)
	configFile := ctx.String(configFileFlag.Name)
	db, err := db.NewCrossChainDB(db.DefaultDBHost, db.DefaultDBPort, db.DefaultDBUser, db.DefaultDBPass, db.DefaultDBName)
	if err != nil {
		log.Errorf("autoSwapDcrn NewChainDB err %v\n", err)
		return err
	}
	LoadConfig(configFile)
	AutoSwapConf.Db = db
	//config := params.LoadConfig(bridgeconfigFile, true)
	//tokenPairConfigMap, err := tokens.LoadTokenPairsConfigInDir(tokenConfigDir, false)
	//tokenPairConfig := tokenPairConfigMap[PairID]
	//if err != nil {
	//	log.Errorf("Load token pairs config fail: %v\n", err)
	//	return err
	//}
	//
	//AutoSwapConf = &autoSwapConf{
	//	SwapServer: swapServer,
	//	//dcrnNetID:       dcrnNetwork,
	//	//dcrnGateway:     dcrnGateWay,
	//	Account:         account,
	//	BridgeConfig:    config,
	//	Db:              db,
	//	TokenPairConfig: tokenPairConfig,
	//}
	return nil
}
