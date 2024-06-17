package dao

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mybtckb-svr/tables"
)

func (d *DbDao) CreateRgbppXudt(xudt []*tables.TableXudt) (err error) {
	return d.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Where("id > 0").Delete(&tables.TableXudt{}).Error; err != nil {
			return err
		}
		if len(xudt) > 0 {
			if err := tx.Create(&xudt).Error; err != nil {
				return err
			}
		}
		return nil
	})

}

func (d *DbDao) UpdateRgbppXudtAddr(xudt []*tables.TableXudt) (err error) {
	return d.db.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"outpoint", "address"}),
	}).Create(xudt).Error
}
func (d *DbDao) GetRgbppXudtWithEmptyAddr() (list []*tables.TableXudt, err error) {
	err = d.db.Model(&tables.TableXudt{}).Where("address = ? and addr_type= ?", "", 1).Limit(5).Find(&list).Error
	return
}

type XudtTokenIdAmount struct {
	TokenId string `json:"token_id"`
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

func (d *DbDao) GetRgbppXudtByAddr(address string) (list []*tables.TableXudt, err error) {
	err = d.db.Model(&tables.TableXudt{}).Where("address = ? ", address).Find(&list).Error
	//err = d.db.Model(tables.TableXudt{}).
	//	Select("address, SUM(amount) amount, token_id").
	//	Where("address=? ", address).
	//	Group("token_id").Find(&list).Error
	return
	//return list, nil
}
