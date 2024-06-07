package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mybtckb-svr/config"
	"mybtckb-svr/tables"
)

type DbDao struct {
	db *gorm.DB
}

func NewGormDB(conf config.DbMysql) (*DbDao, error) {
	conn := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(conn, conf.User, conf.Password, conf.Addr, conf.DbName)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("gorm open :%v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gorm db :%v", err)
	}
	sqlDB.SetMaxOpenConns(conf.MaxOpenConn)
	sqlDB.SetMaxIdleConns(conf.MaxIdleConn)

	if err := db.AutoMigrate(
		&tables.TableBlockParserInfo{},
		&tables.TableSpore{},
		&tables.TableClusterInfo{},
		&tables.TableXudtInfo{},
		&tables.TableXudt{},
		&tables.TransactionInfo{},
		&tables.XudtTransferRecord{},
	); err != nil {
		return nil, err
	}

	return &DbDao{db: db}, nil
}
