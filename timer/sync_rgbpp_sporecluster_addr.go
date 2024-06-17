package timer

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/lib/common"
	"mybtckb-svr/tables"
)

func (t *TxTimer) doSyncRgbppSporeClusterAddr() error {
	rgbppSporeCluster, err := t.dbDao.GetRgbppSporeClusterWithEmptyAddr()
	if err != nil {
		return fmt.Errorf("GetRgbppSporeWithEmptyAddr err :", err.Error())
	}
	fmt.Println("need update: ", len(rgbppSporeCluster))
	sporeClusterList := make([]*tables.TableClusterInfo, 0)
	for _, v := range rgbppSporeCluster {
		if v.BtcOutpoint == "" {
			log.Warn("rgbpp spor_cluster btcoutpoint is empty: ", v.ClusterId, v.Outpoint, v.Id)
			continue
		}
		fmt.Println("btc out point ", v.BtcOutpoint)
		txHash, index := common.String2OutPoint(v.BtcOutpoint)

		owner, err := t.contracts.GetBtcAddressByOutpoint(uint32(index), txHash)
		if err != nil {
			log.Errorf("GetBtcAddressByOutpoint err %s", err.Error())
			continue
		}
		sporeClusterList = append(sporeClusterList, &tables.TableClusterInfo{
			ClusterId:   v.ClusterId,
			BlockNum:    v.BlockNum,
			Outpoint:    v.Outpoint,
			BtcOutpoint: v.BtcOutpoint,
			Address:     owner,
		})
	}
	if len(sporeClusterList) == 0 {
		return nil
	}
	fmt.Println("update : ", len(sporeClusterList))
	if err := t.dbDao.CreateClusterInfo(sporeClusterList); err != nil {
		return fmt.Errorf("CreateSporeInfo err: %s", err.Error())
	}
	return nil
}
