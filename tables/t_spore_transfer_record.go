package tables

import (
	"time"
)

type SporeTransferRecord struct {
	Id             uint64    `json:"id" gorm:"column:id;primaryKey;type:bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT ''"`
	BlockNum       uint64    `json:"block_num" gorm:"column:block_num;index:idx_block_num_idx;type:bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT ''"`
	TxHash         string    `json:"tx_hash" gorm:"column:tx_hash;uniqueIndex:uk_tx_receipt_addr;type:varchar(255) NOT NULL DEFAULT '' COMMENT ''"`
	Address        string    `json:"address" gorm:"column:address;index:idx_address;type:varchar(255) NOT NULL DEFAULT '' COMMENT ''"`
	ReceiptAddress string    `json:"receipt_address" gorm:"column:receipt_address;uniqueIndex:uk_tx_receipt_addr;type:varchar(255) NOT NULL DEFAULT '' COMMENT ''"`
	Capacity       uint64    `json:"capacity" gorm:"column:capacity;type:bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT ''"`
	SporeId        string    `json:"spore_id" gorm:"column:spore_id;type:varchar(255) NOT NULL DEFAULT '' COMMENT ''"`
	Status         TxStatus  `json:"status" gorm:"column:status;type:smallint(6) NOT NULL DEFAULT '0' COMMENT '0-default 1-rejected'"`
	BlockTimestamp uint64    `json:"block_timestamp" gorm:"column:block_timestamp;type:bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT ''"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''"`
}

func (t *SporeTransferRecord) TableName() string {
	return "t_spore_transfer_record"
}
