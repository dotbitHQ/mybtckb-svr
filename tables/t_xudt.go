package tables

import "time"

type TableXudt struct {
	Id        uint64    `json:"id" gorm:"column:id; primaryKey; type:bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '';"`
	TokenId   string    `json:"token_id" gorm:"column:token_id; uniqueIndex:unique_token_address; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	Address   string    `json:"address" gorm:"column:address; uniqueIndex:unique_token_address; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	Amount    uint8     `json:"amount" gorm:"column:amount; type:tinyint(4) NOT NULL DEFAULT '0' COMMENT '';"`
	BlockNum  uint64    `json:"block_num" gorm:"column:block_num;type:bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT ''"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at; type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '';"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at; type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '';"`
}

func (t *TableXudt) TableName() string {
	return "t_xudt"
}
