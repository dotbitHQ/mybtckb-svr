package txtool

import (
	"context"
	"errors"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/lib/common"
	"mybtckb-svr/lib/math/uint128"
)

var (
	ErrInsufficientFunds = errors.New("InsufficientFunds")
)

type ParamGetXudtCells struct {
	LockScript *types.Script
	TokenId    string
	AmountNeed uint64
}

func (txTool *TxTool) GetXudtCells(p *ParamGetXudtCells) ([]*indexer.LiveCell, uint128.Uint128, *types.Script, error) {
	total := uint128.Zero
	xudtType := txTool.contract.GetXudtTypeScript(common.Hex2Bytes(p.TokenId))
	searchKey := &indexer.SearchKey{
		Script:     p.LockScript,
		ScriptType: indexer.ScriptTypeLock,
		Filter: &indexer.CellsFilter{
			Script: xudtType,
		},
	}
	var cells []*indexer.LiveCell
	lastCursor := ""

	ok := false
	for {
		liveCells, err := txTool.contract.Client().GetCells(context.Background(), searchKey, indexer.SearchOrderDesc, indexer.SearchLimit, lastCursor)
		if err != nil {
			return nil, total, nil, err
		}
		log.Info("liveCells:", liveCells.LastCursor, len(liveCells.Objects))
		if len(liveCells.Objects) == 0 || lastCursor == liveCells.LastCursor {
			break
		}
		lastCursor = liveCells.LastCursor

		for _, liveCell := range liveCells.Objects {
			cells = append(cells, liveCell)
			total = total.Add(uint128.FromBytes(liveCell.OutputData))
			if p.AmountNeed > 0 && total.Cmp64(p.AmountNeed) >= 0 {
				ok = true
				break
			}
		}
		if ok {
			break
		}
	}

	if p.AmountNeed > 0 {
		if total.Cmp64(p.AmountNeed) < 0 {
			return cells, total, nil, ErrInsufficientFunds
		}
	}
	return cells, total, xudtType, nil
}
