package contract

import (
	"context"
	"mybtckb-svr/config"
	"mybtckb-svr/lib/contract"
	"mybtckb-svr/lib/outpoint_cache"

	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

func Init(ctx context.Context, wg *sync.WaitGroup, configPath string) (*contract.Contracts, error) {
	logrus.SetOutput(os.Stdout)
	if err := config.InitCfg(configPath); err != nil {
		return nil, err
	}
	ckbCli, err := rpc.DialWithIndexer(config.Cfg.Chain.Ckb.Url, config.Cfg.Chain.Ckb.Url)
	if err != nil {
		return nil, err
	}
	txCache := outpoint_cache.NewCache(ctx, wg)
	txCache.RunClearExpiredOutPoint(time.Minute * 5)
	contracts := contract.NewContracts(
		ctx,
		contract.WithWaitGroup(wg),
		contract.WithClient(ckbCli),
		contract.WithNetType(config.Cfg.Server.Net),
		contract.WithCache(txCache),
	)
	if err := contracts.Init(time.Minute * 5); err != nil {
		return nil, err
	}
	return contracts, nil
}
