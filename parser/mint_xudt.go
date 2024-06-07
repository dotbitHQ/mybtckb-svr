package parser

import (
	"fmt"
	"github.com/status-im/keycard-go/hexutils"
	"mybtckb-svr/lib/common"
	"mybtckb-svr/lib/contract"
	"mybtckb-svr/tables"
)

func (b *BlockParser) MintXudt(req FuncTransactionHandleReq) error {

	xudtInfo, err := contract.ParseXudtInfo(req.TxTypeData.OutputCell.Unique[0].OutputData)
	if err != nil {
		return fmt.Errorf("contract.ParseXudtInfo err: %s", err.Error())
	}
	//fmt.Println("-------------tx hash: ", req.TxHash)
	//fmt.Println("-------------xudtInfo : ", xudtInfo)
	tokenId := hexutils.BytesToHex(req.TxTypeData.OutputCell.Xudt[0].Output.Type.Args)
	xudtInfooutPoint := common.OutPoint2String(req.TxHash, req.TxTypeData.OutputCell.Unique[0].Index)
	xudtInfoRecord := &tables.TableXudtInfo{
		TokenId:   tokenId,
		Decimal:   xudtInfo.Decimal,
		Name:      xudtInfo.Name,
		Symbol:    xudtInfo.Symbol,
		Timestamp: req.BlockTimestamp,
		Outpoint:  xudtInfooutPoint,
		BlockNum:  req.BlockNumber,
	}
	if err = b.dbDao.CreateXudtInfo(xudtInfoRecord); err != nil {
		return fmt.Errorf("CreateXudtInfo err : %s", err.Error())
	}

	return nil
	//return b.dbDao.Transaction(func(tx *gorm.DB) error {
	//	if err := tx.Save(&tick).Error; err != nil {
	//		return err
	//	}
	//	return tx.Clauses(clause.OnConflict{
	//		DoUpdates: clause.AssignmentColumns([]string{"tx_hash"}),
	//	}).Create(&tables.TransactionInfo{
	//		BlockNum:       req.BlockNumber,
	//		BlockIdx:       req.BlockIndex,
	//		TxHash:         req.TxHash,
	//		Action:         contract.ActionConfirmBurn,
	//		Capacity:       req.Tx.Outputs[0].Capacity,
	//		BlockTimestamp: req.BlockTimestamp,
	//	}).Error
	//})
}
