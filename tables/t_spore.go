package tables

import "time"

type TableSpore struct {
	Id          uint64    `json:"id" gorm:"column:id; primaryKey; type:bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '';"`
	Address     string    `json:"address" gorm:"column:address; uniqueIndex:unique_address_sporeid; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	SporeId     string    `json:"spore_id" gorm:"column:spore_id; uniqueIndex:unique_address_sporeid; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	ClusterId   string    `json:"cluster_id" gorm:"column:cluster_id; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	ContentType string    `json:"content_type" gorm:"column:content_type; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	Content     []byte    `json:"content" gorm:"column:content; type:longblob NOT NULL  COMMENT '';"`
	Outpoint    string    `json:"outpoint" gorm:"column:outpoint; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	BlockNum    uint64    `json:"block_num" gorm:"column:block_num;type:bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT ''"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at; type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '';"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at; type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '';"`
}

func (t *TableSpore) TableName() string {
	return "t_spore"
}
