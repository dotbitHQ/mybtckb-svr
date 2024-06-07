package timer

import (
	"context"
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/tables"
)

func (t *TxTimer) doCheckXudtTxStatus() error {
	pendingXudtTransferRecords, err := t.dbDao.GetPendingXudtTransferRecord()
	if err != nil {
		return fmt.Errorf("GetPendingXudtTransferRecord err : ", err.Error())
	}
	for _, v := range pendingXudtTransferRecords {
		res, err := t.contracts.Client().GetTransaction(context.Background(), types.HexToHash(v.TxHash))
		if err != nil {
			log.Warn("GetTransaction err: ", err.Error())
			continue
		}

		xudtTransferRecord := &tables.XudtTransferRecord{
			TxHash: res.Transaction.Hash.Hex(),
		}
		if res.TxStatus.Status == types.TransactionStatusCommitted {
			if res.TxStatus.Status == types.TransactionStatusCommitted {
				if block, err := t.contracts.Client().GetBlock(t.ctx, *res.TxStatus.BlockHash); err == nil && block != nil && block.Header != nil {
					log.Info("UpdatePendingToConfirm:", v.Id)
					xudtTransferRecord.BlockNum = block.Header.Number
					xudtTransferRecord.BlockTimestamp = block.Header.Timestamp
					xudtTransferRecord.Status = tables.StatusConfirm
				}
			}
		} else if res.TxStatus.Status == types.TransactionStatusRejected {
			xudtTransferRecord.Status = tables.StatusRejected
		}
		if err := t.dbDao.UpdateXudtTransferRecord(xudtTransferRecord); err != nil {
			return fmt.Errorf("UpdateXudtTransferRecord err : %s", err.Error())
		}
	}
	return nil
}

func (t *TxTimer) doCheckSporeTxStatus() error {
	pendingSporeTransferRecords, err := t.dbDao.GetPendingSporeTransferRecord()
	if err != nil {
		return fmt.Errorf("GetPendingXudtTransferRecord err : ", err.Error())
	}
	for _, v := range pendingSporeTransferRecords {
		res, err := t.contracts.Client().GetTransaction(context.Background(), types.HexToHash(v.TxHash))
		if err != nil {
			log.Warn("GetTransaction err: ", err.Error())
			continue
		}

		sporeTransferRecord := &tables.SporeTransferRecord{
			TxHash: res.Transaction.Hash.Hex(),
		}
		if res.TxStatus.Status == types.TransactionStatusCommitted {
			if res.TxStatus.Status == types.TransactionStatusCommitted {
				if block, err := t.contracts.Client().GetBlock(t.ctx, *res.TxStatus.BlockHash); err == nil && block != nil && block.Header != nil {
					log.Info("UpdatePendingToConfirm:", v.Id)
					sporeTransferRecord.BlockNum = block.Header.Number
					sporeTransferRecord.BlockTimestamp = block.Header.Timestamp
					sporeTransferRecord.Status = tables.StatusConfirm
				}
			}
		} else if res.TxStatus.Status == types.TransactionStatusRejected {
			sporeTransferRecord.Status = tables.StatusRejected
		}
		if err := t.dbDao.UpdateSporeTransferRecord(sporeTransferRecord); err != nil {
			return fmt.Errorf("UpdateSporeTransferRecord err : %s", err.Error())
		}
	}
	return nil
}

func (t *TxTimer) doCheckCkbTxStatus() error {
	pendingCkbTransferRecords, err := t.dbDao.GetPendingCkbTransferRecord()
	if err != nil {
		return fmt.Errorf("GetPendingXudtTransferRecord err : ", err.Error())
	}
	for _, v := range pendingCkbTransferRecords {
		res, err := t.contracts.Client().GetTransaction(context.Background(), types.HexToHash(v.TxHash))
		if err != nil {
			log.Warn("GetTransaction err: ", err.Error())
			continue
		}

		ckbTransferRecord := &tables.CkbTransferRecord{
			TxHash: res.Transaction.Hash.Hex(),
		}
		if res.TxStatus.Status == types.TransactionStatusCommitted {
			if res.TxStatus.Status == types.TransactionStatusCommitted {
				if block, err := t.contracts.Client().GetBlock(t.ctx, *res.TxStatus.BlockHash); err == nil && block != nil && block.Header != nil {
					log.Info("UpdatePendingToConfirm:", v.Id)
					ckbTransferRecord.BlockNum = block.Header.Number
					ckbTransferRecord.BlockTimestamp = block.Header.Timestamp
					ckbTransferRecord.Status = tables.StatusConfirm
				}
			}
		} else if res.TxStatus.Status == types.TransactionStatusRejected {
			ckbTransferRecord.Status = tables.StatusRejected
		}
		if err := t.dbDao.UpdateCkbTransferRecord(ckbTransferRecord); err != nil {
			return fmt.Errorf("UpdateSporeTransferRecord err : %s", err.Error())
		}
	}
	return nil
}
