package tables

import "time"

type TableXudtInfo struct {
	Id        uint64    `json:"id" gorm:"column:id; primaryKey; type:bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '';"`
	TokenId   string    `json:"token_id" gorm:"column:token_id; index:idx_token_id; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	Decimal   uint8     `json:"decimal" gorm:"column:decimal; type:tinyint(4) NOT NULL DEFAULT '0' COMMENT '';"`
	Name      string    `json:"name" gorm:"column:name; index:idx_name; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	Symbol    string    `json:"symbol" gorm:"column:symbol; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	Timestamp uint64    `json:"timestamp" gorm:"column:timestamp; type:bigint(20) NOT NULL DEFAULT '0' COMMENT '';"`
	Outpoint  string    `json:"outpoint" gorm:"column:outpoint; index:idx_outpoint; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	BlockNum  uint64    `json:"block_num" gorm:"column:block_num;type:bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT ''"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at; type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '';"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at; type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '';"`
}

func (t *TableXudtInfo) TableName() string {
	return "t_xudt_info"
}
