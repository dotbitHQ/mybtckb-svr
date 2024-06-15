package timer

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/tables"
)

func (t *TxTimer) doSyncRgbppSporeAddr() error {
	rgbppSpore, err := t.dbDao.GetRgbppSporeWithEmptyAddr()
	if err != nil {
		return fmt.Errorf("GetRgbppSporeWithEmptyAddr err :", err.Error())
	}
	fmt.Println("need update: ", len(rgbppSpore))
	sporeList := make([]*tables.TableSpore, 0)
	for _, v := range rgbppSpore {
		if v.Outpoint == "" {
			log.Warn("rgbpp spore btcoutpoint is empty: ", v.SporeId, v.Outpoint, v.Id)
			continue
		}
		fmt.Println("btc out point ", v.BtcOutpoint)
		txHash, index := common.String2OutPoint(v.BtcOutpoint)

		owner, err := t.contracts.GetBtcAddressByOutpoint(uint32(index), txHash)
		if err != nil {
			log.Errorf("GetBtcAddressByOutpoint err %s", err.Error())
			continue
		}
		sporeList = append(sporeList, &tables.TableSpore{
			SporeId:     v.SporeId,
			BlockNum:    v.BlockNum,
			Outpoint:    v.Outpoint,
			BtcOutpoint: v.BtcOutpoint,
			Address:     owner,
			Content:     []byte{},
		})
	}
	if len(sporeList) == 0 {
		return nil
	}
	fmt.Println("update : ", len(sporeList))
	if err := t.dbDao.CreateSporeInfo(sporeList); err != nil {
		return fmt.Errorf("CreateSporeInfo err: %s", err.Error())
	}
	return nil
}
