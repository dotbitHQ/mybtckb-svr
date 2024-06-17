package timer

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/tables"
)

func (t *TxTimer) doSyncRgbppXudtAddr() error {
	rgbppXudt, err := t.dbDao.GetRgbppXudtWithEmptyAddr()
	if err != nil {
		return fmt.Errorf("GetRgbppSporeWithEmptyAddr err :", err.Error())
	}
	fmt.Println("need update: ", len(rgbppXudt))
	xudtList := make([]*tables.TableXudt, 0)
	for _, v := range rgbppXudt {
		if v.BtcOutpoint == "" {
			log.Warn("rgbpp spore btcoutpoint is empty: ", v.Outpoint, v.Id)
			continue
		}
		fmt.Println("btc out point ", v.BtcOutpoint)
		txHash, index := common.String2OutPoint(v.BtcOutpoint)

		owner, err := t.contracts.GetBtcAddressByOutpoint(uint32(index), txHash)
		if err != nil {
			log.Errorf("GetBtcAddressByOutpoint err %s", err.Error())
			continue
		}
		xudtList = append(xudtList, &tables.TableXudt{
			Outpoint: v.Outpoint,
			Address:  owner,
		})
	}
	if len(xudtList) == 0 {
		return nil
	}
	fmt.Println("update : ", len(xudtList))
	if err := t.dbDao.UpdateRgbppXudtAddr(xudtList); err != nil {
		return fmt.Errorf("CreateSporeInfo err: %s", err.Error())
	}
	return nil
}
