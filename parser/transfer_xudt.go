package parser

import (
	"fmt"
	"mybtckb-svr/tables"
)

func (b *BlockParser) TransferXudt(req FuncTransactionHandleReq) error {

	//xudtInfo, err := contract.ParseXudtInfo(req.TxTypeData.OutputCell.Unique[0].OutputData)
	//if err != nil {
	//	return fmt.Errorf("contract.ParseXudtInfo err: %s", err.Error())
	//}
	//fmt.Println("-------------tx hash: ", req.TxHash)
	//fmt.Println("-------------xudtInfo : ", xudtInfo)

	//tokenId := hexutils.BytesToHex(req.TxTypeData.OutputCell.Xudt[0].Output.Type.Args)
	//
	//address, err := address.ConvertScriptToAddress(b.contracts.Mode(), req.TxTypeData.InputCell.Xudt[0].Output.Lock)
	//if err != nil {
	//	return fmt.Errorf("ConvertScriptToAddress err: %s", err.Error())
	//}
	//for _, v := range req.TxTypeData.OutputCell.Xudt {
	//
	//}
	xudtTransferRecord := &tables.XudtTransferRecord{
		TxHash:         req.TxHash,
		BlockNum:       req.BlockNumber,
		BlockTimestamp: req.BlockTimestamp,
		Status:         tables.StatusConfirm,
	}
	if err := b.dbDao.CreateXudtTransferRecord(xudtTransferRecord); err != nil {
		return fmt.Errorf("UpdateXudtTransferRecord err : %s", err.Error())
	}
	return nil
}