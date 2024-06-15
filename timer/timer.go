package timer

import (
	"context"
	"mybtckb-svr/dao"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/lib/contract"
	"mybtckb-svr/txtool"

	"sync"
	"time"
)

type TxTimer struct {
	ctx       context.Context
	wg        *sync.WaitGroup
	contracts *contract.Contracts
	dbDao     *dao.DbDao
	txTool    *txtool.TxTool
	cron      *cron.Cron
}

type TxTimerParam struct {
	Ctx       context.Context
	Wg        *sync.WaitGroup
	DbDao     *dao.DbDao
	Contracts *contract.Contracts
	Txtool    *txtool.TxTool
}

func NewTxTimer(p TxTimerParam) *TxTimer {
	var t TxTimer
	t.ctx = p.Ctx
	t.wg = p.Wg
	t.dbDao = p.DbDao
	t.contracts = p.Contracts
	t.txTool = p.Txtool
	return &t
}

func (t *TxTimer) Run() error {

	//if err := t.doSyncRgbppSporeAddr(); err != nil {
	//	log.Error("doSyncRgbppSporeAddr err: ", err.Error())
	//}
	//return nil
	tickerXudtTxStatus := time.NewTicker(time.Second * 5)
	tickerSporeTxStatus := time.NewTicker(time.Second * 5)
	tickerCkbTxStatus := time.NewTicker(time.Second * 5)

	tickerSyncCluster := time.NewTicker(time.Minute * 10)
	tickerSyncSpore := time.NewTicker(time.Minute * 12)
	tickerSyncRgbppSporeAddr := time.NewTicker(time.Second * 20)

	t.wg.Add(1)
	go func() {
		for {
			select {
			case <-tickerXudtTxStatus.C:
				log.Debug("doCheckTxStatus start ...")
				if err := t.doCheckXudtTxStatus(); err != nil {
					log.Error("doCheckTxStatus err: ", err.Error())
				}
				log.Debug("doCheckTxStatus end ...")
			case <-tickerCkbTxStatus.C:
				log.Debug("doCheckCkbTxStatus start ...")
				if err := t.doCheckCkbTxStatus(); err != nil {
					log.Error("doCheckCkbTxStatus err: ", err.Error())
				}
				log.Debug("doCheckCkbTxStatus end ...")
			case <-tickerSporeTxStatus.C:
				log.Debug("doCheckSporeTxStatus start ...")
				if err := t.doCheckSporeTxStatus(); err != nil {
					log.Error("doCheckSporeTxStatus err: ", err.Error())
				}
				log.Debug("doCheckSporeTxStatus end ...")
			case <-tickerSyncCluster.C:
				log.Debug("doSyncCluster start ...")
				if err := t.doSyncCluster(); err != nil {
					log.Error("doSyncCluster err: ", err.Error())
				}
				log.Debug("doSyncCluster end ...")
			case <-tickerSyncSpore.C:
				log.Debug("doSyncSpore start ...")
				if err := t.doSyncSpore(); err != nil {
					log.Error("doSyncSpore err: ", err.Error())
				}
				log.Debug("doSyncSpore end ...")
			case <-tickerSyncRgbppSporeAddr.C:
				log.Debug("doSyncRgbppSporeAddr start ...")
				if err := t.doSyncRgbppSporeAddr(); err != nil {
					log.Error("doSyncRgbppSporeAddr err: ", err.Error())
				}
				log.Debug("doSyncRgbppSporeAddr end ...")
			case <-t.ctx.Done():
				log.Debug("timer done")
				t.wg.Done()
				return

			}
		}
	}()

	return nil
}
