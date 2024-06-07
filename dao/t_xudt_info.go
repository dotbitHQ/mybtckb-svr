package dao

import "mybtckb-svr/tables"

func (d *DbDao) CreateXudtInfo(xudtInfo *tables.TableXudtInfo) (err error) {
	return d.db.Create(xudtInfo).Error
}

func (d *DbDao) GetXudtInfo() (xudtInfoList []tables.TableXudtInfo, err error) {
	err = d.db.Find(&xudtInfoList).Error
	return
}
