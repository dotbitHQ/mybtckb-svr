package contract

import (
	"context"
	"github.com/btcsuite/btcd/rpcclient"
	"log"
	"mybtckb-svr/config"
	"mybtckb-svr/lib/cache"
	"mybtckb-svr/lib/contract"
	"mybtckb-svr/lib/outpoint_cache"

	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

func Init(ctx context.Context, wg *sync.WaitGroup, configPath string, rc *cache.RedisCache) (*contract.Contracts, error) {
	logrus.SetOutput(os.Stdout)
	if err := config.InitCfg(configPath); err != nil {
		return nil, err
	}
	ckbCli, err := rpc.DialWithIndexer(config.Cfg.Chain.Ckb.Url, config.Cfg.Chain.Ckb.Url)
	if err != nil {
		return nil, err
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         config.Cfg.Chain.Btc.Url, // 比特币节点的地址和端口
		User:         "root",
		Pass:         "root",
		HTTPPostMode: true,  // 比特币节点使用HTTP POST模式
		DisableTLS:   false, // 如果启用了TLS，需要设置为false，并提供相应的证书
	}
	connCfg.ExtraHeaders = make(map[string]string, 0)
	connCfg.ExtraHeaders["api-key"] = config.Cfg.Chain.Btc.ApiKey
	// 创建新的RPC客户端
	btcClient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatalf("Error creating new RPC client: %v", err)
	}

	txCache := outpoint_cache.NewCache(ctx, wg)
	txCache.RunClearExpiredOutPoint(time.Minute * 5)
	contracts := contract.NewContracts(
		ctx,
		contract.WithWaitGroup(wg),
		contract.WithCkbClient(ckbCli),
		contract.WithBtcClient(*btcClient),
		contract.WithRedisClient(rc),
		contract.WithNetType(config.Cfg.Server.Net),
		contract.WithCache(txCache),
	)
	if err := contracts.Init(time.Minute * 5); err != nil {
		return nil, err
	}
	return contracts, nil
}
