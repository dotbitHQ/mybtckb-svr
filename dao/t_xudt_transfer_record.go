package dao

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mybtckb-svr/tables"
)

func (d *DbDao) GetPendingXudtTransferRecord() (list []*tables.XudtTransferRecord, err error) {
	err = d.db.Where("status = ? ", tables.StatusPending).Order("id ASC").Limit(100).Find(&list).Error
	return
}
func (d *DbDao) GetAllXudtTransferRecordByAddress(address string, page, pageSize uint64) (list []*tables.XudtTransferRecord, total int64, err error) {
	err = d.db.Model(&tables.XudtTransferRecord{}).Where("address = ? or receipt_address = ?", address, address).Count(&total).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	err = d.db.Where("address = ? or receipt_address = ?", address, address).Offset(int(page-1) * int(pageSize)).Limit(int(pageSize)).Find(&list).Error
	return

}

func (d *DbDao) CreateXudtTransferRecord(xudtTransferRecord *tables.XudtTransferRecord) error {
	return d.db.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"block_num", "status", "block_timestamp"}),
	}).Create(xudtTransferRecord).Error
}

func (d *DbDao) UpdateXudtTransferRecord(xudtTransferRecord *tables.XudtTransferRecord) error {
	return d.db.Model(tables.XudtTransferRecord{}).
		Where("tx_hash= ? ", xudtTransferRecord.TxHash).
		Updates(map[string]interface{}{
			"block_num":       xudtTransferRecord.BlockNum,
			"status":          xudtTransferRecord.Status,
			"block_timestamp": xudtTransferRecord.BlockTimestamp,
		}).Error
}
