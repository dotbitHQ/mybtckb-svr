package dao

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mybtckb-svr/tables"
)

func (d *DbDao) CreateSporeInfo(sporeInfo []*tables.TableSpore) (err error) {
	return d.db.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"block_num", "outpoint", "btc_outpoint", "address"}),
	}).Create(sporeInfo).Error
}

func (d *DbDao) GetSporeByClusterId(clusterId string, page, pageSize uint64) (list []*tables.TableSpore, total int64, err error) {
	err = d.db.Model(&tables.TableSpore{}).Where("cluster_id = ? ", clusterId).Count(&total).Error
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	err = d.db.Where("cluster_id = ? ", clusterId).Offset(int(page-1) * int(pageSize)).Limit(int(pageSize)).Find(&list).Error
	return

}
func (d *DbDao) GetSporeByAddress(address string, page, pageSize uint64) (list []*tables.TableSpore, total int64, err error) {
	err = d.db.Model(&tables.TableSpore{}).Where("address = ? ", address).Count(&total).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	err = d.db.Where("address = ? ", address).Offset(int(page-1) * int(pageSize)).Limit(int(pageSize)).Find(&list).Error
	return

}

func (d *DbDao) GetRgbppSporeWithEmptyAddr() (list []*tables.TableSpore, err error) {
	err = d.db.Model(&tables.TableSpore{}).Where("address = ? and addr_type= ?", "", 1).Limit(5).Find(&list).Error
	return
}
