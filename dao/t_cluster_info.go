package dao

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mybtckb-svr/tables"
)

func (d *DbDao) CreateClusterInfo(clusterInfo []*tables.TableClusterInfo) (err error) {
	return d.db.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"block_num", "outpoint"}),
	}).Create(clusterInfo).Error
}

func (d *DbDao) GetClusterByAddress(address string, page, pageSize uint64) (list []*tables.TableClusterInfo, total int64, err error) {
	err = d.db.Model(&tables.TableClusterInfo{}).Where("address = ? ", address).Count(&total).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	err = d.db.Where("address = ? ", address).Offset(int(page-1) * int(pageSize)).Limit(int(pageSize)).Find(&list).Error
	return

}
