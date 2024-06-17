package dao

import (
	"gorm.io/gorm/clause"
	"mybtckb-svr/tables"
)

func (d *DbDao) CreateXudtInfo(xudtInfo *tables.TableXudtInfo) (err error) {
	return d.db.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"block_num", "outpoint", "name", "symbol", "decimal"}),
	}).Create(xudtInfo).Error
}

func (d *DbDao) GetXudtInfo() (xudtInfoList []tables.TableXudtInfo, err error) {
	err = d.db.Find(&xudtInfoList).Error
	return
}
