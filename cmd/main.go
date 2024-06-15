package main

import (
	"context"
	"fmt"
	"github.com/scorpiotzh/toolib"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	_ "mybtckb-svr/cmd/docs"
	"mybtckb-svr/config"
	"mybtckb-svr/contract"
	"mybtckb-svr/dao"
	"mybtckb-svr/http_server"
	"mybtckb-svr/http_server/handle"
	"mybtckb-svr/lib/cache"
	"mybtckb-svr/lib/common"
	"mybtckb-svr/lib/signal"
	"mybtckb-svr/lib/txbuilder"
	"mybtckb-svr/parser"
	"mybtckb-svr/tables"
	"mybtckb-svr/timer"
	"mybtckb-svr/txtool"
	"os"
	"sync"
	"time"
)

var (
	//log               = logger.NewLogger("main", logger.LevelDebug)
	exit              = make(chan struct{})
	ctxServer, cancel = context.WithCancel(context.Background())
	wgServer          = sync.WaitGroup{}
)

func main() {
	log.Debugf("server start：")
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
		Action: runServer,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runServer(ctx *cli.Context) error {

	// config file
	configFilePath := ctx.String("config")
	if err := config.InitCfg(configFilePath); err != nil {
		return err
	}

	// ============= service start =============

	// redis
	red, err := toolib.NewRedisClient(config.Cfg.DB.Redis.Addr, config.Cfg.DB.Redis.Password, config.Cfg.DB.Redis.DbNum)
	if err != nil {
		return fmt.Errorf("NewRedisClient err:%s", err.Error())
	} else {
		log.Info("redis ok")
	}
	rc := &cache.RedisCache{
		Ctx: ctxServer,
		Red: red,
	}
	contracts, err := contract.Init(ctxServer, &wgServer, configFilePath, rc)
	if err != nil {
		return err
	}

	// db
	dbDao, err := dao.NewGormDB(config.Cfg.DB.Mysql)
	if err != nil {
		return fmt.Errorf("dao.NewGormDB err: %s", err.Error())
	}

	serverArgs, err := common.GetLockFromPk(config.Cfg.Chain.Ckb.PrivateKey)
	if err != nil {
		return err
	}

	txBuilderBase := &txbuilder.TxBuilderBase{
		Ctx:        ctxServer,
		Contracts:  contracts,
		ServerArgs: common.Bytes2Hex(serverArgs.Args),
	}

	localCache := cache.NewCache(ctxServer, &wgServer)
	txTool := txtool.NewTxTool(
		txtool.WithCache(localCache),
		txtool.WithDb(dbDao),
		txtool.WithContract(contracts),
		txtool.WithTxBuilderBase(txBuilderBase),
	)
	blockParser := parser.NewBlockParser(
		parser.WithContext(ctxServer, cancel),
		parser.WithWg(&wgServer),
		parser.WithCurrentBlockNumber(config.Cfg.Chain.Ckb.CurrentBlockNumber),
		parser.WithConcurrentNum(config.Cfg.Chain.Ckb.ConcurrencyNum),
		parser.WithConfirmNum(config.Cfg.Chain.Ckb.ConfirmNum),
		parser.WithParserType(tables.ParserTypeCKB),
		parser.WithContracts(contracts),
		parser.WithDbDao(dbDao),
		parser.WithTxTool(txTool),
	)
	if err := blockParser.Run(); err != nil {
		return err
	}

	txTimer := timer.NewTxTimer(timer.TxTimerParam{
		Ctx:       ctxServer,
		Wg:        &wgServer,
		DbDao:     dbDao,
		Contracts: contracts,
		Txtool:    txTool,
	})
	if err = txTimer.Run(); err != nil {
		return fmt.Errorf("txTimer.Run() err: %s", err.Error())
	}

	log.Info("timer ok")

	hs := http_server.HttpServer{
		Ctx:     ctxServer,
		Address: config.Cfg.Server.HttpServerAddr,
		H: &handle.HttpHandle{
			Ctx:          ctxServer,
			DbDao:        dbDao,
			RC:           rc,
			Contracts:    contracts,
			TxTool:       txTool,
			ServerScript: serverArgs,
		},
	}
	hs.Run()
	log.Info("http server ok")

	// ============= service end =============
	signal.AddInterruptHandler(func() {
		cancel()
		wgServer.Wait()
		log.Warn("success exit server. bye bye!")
		time.Sleep(time.Second)
		exit <- struct{}{}
	})
	<-exit
	return nil
}
