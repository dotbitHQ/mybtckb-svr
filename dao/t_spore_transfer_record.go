package dao

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mybtckb-svr/tables"
)

func (d *DbDao) GetPendingSporeTransferRecord() (list []*tables.SporeTransferRecord, err error) {
	err = d.db.Where("status = ? ", tables.StatusPending).Order("id ASC").Limit(100).Find(&list).Error
	return
}

func (d *DbDao) GetAllSporeTransferRecordByAddress(address string, page, pageSize uint64) (list []*tables.SporeTransferRecord, total int64, err error) {
	err = d.db.Model(&tables.SporeTransferRecord{}).Where("address = ? or receipt_address = ?", address, address).Count(&total).Error
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	err = d.db.Where("address = ? or receipt_address = ?", address, address).Offset(int(page-1) * int(pageSize)).Limit(int(pageSize)).Find(&list).Error
	return

}

func (d *DbDao) CreateSporeTransferRecord(sporeTransferRecord *tables.SporeTransferRecord) error {
	return d.db.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"block_num", "status", "block_timestamp"}),
	}).Create(sporeTransferRecord).Error
}

func (d *DbDao) UpdateSporeTransferRecord(xudtTransferRecord *tables.SporeTransferRecord) error {
	return d.db.Model(tables.SporeTransferRecord{}).
		Where("tx_hash= ? ", xudtTransferRecord.TxHash).
		Updates(map[string]interface{}{
			"block_num":       xudtTransferRecord.BlockNum,
			"status":          xudtTransferRecord.Status,
			"block_timestamp": xudtTransferRecord.BlockTimestamp,
		}).Error
}
