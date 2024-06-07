package dao

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mybtckb-svr/tables"
)

func (d *DbDao) GetPendingCkbTransferRecord() (list []*tables.CkbTransferRecord, err error) {
	err = d.db.Where("status = ? ", tables.StatusPending).Order("id ASC").Limit(100).Find(&list).Error
	return
}
func (d *DbDao) GetAllCkbTransferRecordByAddress(address string, page, pageSize uint64) (list []*tables.CkbTransferRecord, total int64, err error) {
	err = d.db.Model(&tables.CkbTransferRecord{}).Where("address = ? or receipt_address = ?", address, address).Count(&total).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return
	}
	err = d.db.Where("address = ? or receipt_address = ?", address, address).Offset(int(page-1) * int(pageSize)).Limit(int(pageSize)).Find(&list).Error
	return

}

func (d *DbDao) CreateCkbTransferRecord(ckbTransferRecord *tables.CkbTransferRecord) error {
	return d.db.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"block_num", "status", "block_timestamp"}),
	}).Create(ckbTransferRecord).Error
}

func (d *DbDao) UpdateCkbTransferRecord(ckbTransferRecord *tables.CkbTransferRecord) error {
	return d.db.Model(tables.CkbTransferRecord{}).
		Where("tx_hash= ? ", ckbTransferRecord.TxHash).
		Updates(map[string]interface{}{
			"block_num":       ckbTransferRecord.BlockNum,
			"status":          ckbTransferRecord.Status,
			"block_timestamp": ckbTransferRecord.BlockTimestamp,
		}).Error
}
