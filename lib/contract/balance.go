package contract

import (
	"context"
	"errors"
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/lib/common"
	"sync"
)

var (
	ErrRejectedOutPoint  = errors.New("RejectedOutPoint")
	ErrInsufficientFunds = errors.New("InsufficientFunds")
	ErrNotEnoughChange   = errors.New("NotEnoughChange")
)

type ParamGetBalanceCells struct {
	LockScript         *types.Script
	TypeScript         *types.Script
	CapacityNeed       uint64
	CapacityForChange  uint64
	CurrentBlockNumber uint64
	SearchOrder        indexer.SearchOrder
	OutputDataLenRange *[2]uint64
}

func (c *Contracts) GetBalanceCells(p *ParamGetBalanceCells) ([]*indexer.LiveCell, uint64, error) {
	if p.OutputDataLenRange == nil {
		p.OutputDataLenRange = &[2]uint64{0, 1}
	}
	return c.GetBalanceCellsFilter(p)
}

func (c *Contracts) GetBalanceCellsFilter(p *ParamGetBalanceCells) ([]*indexer.LiveCell, uint64, error) {
	if c.client == nil {
		return nil, 0, fmt.Errorf("client is nil")
	}
	if p == nil {
		return nil, 0, fmt.Errorf("param is nil")
	}

	searchKey := &indexer.SearchKey{
		Script:     p.LockScript,
		ScriptType: indexer.ScriptTypeLock,
		Filter: &indexer.CellsFilter{
			Script: p.TypeScript,
		},
	}
	if p.OutputDataLenRange != nil {
		searchKey.Filter.OutputDataLenRange = p.OutputDataLenRange
	}
	if p.CurrentBlockNumber > 0 {
		searchKey.Filter.BlockRange = &[2]uint64{0, p.CurrentBlockNumber - 20}
	}

	var cells []*indexer.LiveCell
	total := uint64(0)
	hasCache := false
	lastCursor := ""

	ok := false
	for {
		liveCells, err := c.client.GetCells(context.Background(), searchKey, p.SearchOrder, indexer.SearchLimit, lastCursor)
		if err != nil {
			return nil, 0, err
		}
		log.Info("liveCells:", liveCells.LastCursor, len(liveCells.Objects))
		if len(liveCells.Objects) == 0 || lastCursor == liveCells.LastCursor {
			break
		}
		lastCursor = liveCells.LastCursor

		for _, liveCell := range liveCells.Objects {
			if liveCell.Output.Type != nil {
				continue
			}
			if p.CapacityNeed > 0 && c.cache != nil && c.cache.ExistOutPoint(common.OutPointStruct2String(liveCell.OutPoint)) {
				hasCache = true
				continue
			}
			cells = append(cells, liveCell)
			total += liveCell.Output.Capacity
			if p.CapacityNeed > 0 {
				if total >= p.CapacityNeed+p.CapacityForChange+common.OneCkb {
					ok = true
					break
				}
			}
		}
		if ok {
			break
		}
	}

	if p.CapacityNeed > 0 {
		if total < p.CapacityNeed {
			if hasCache {
				return cells, total, ErrRejectedOutPoint
			} else {
				return cells, total, ErrInsufficientFunds
			}
		} else if total < p.CapacityNeed+p.CapacityForChange {
			if hasCache {
				return cells, total, ErrRejectedOutPoint
			} else {
				return cells, total, ErrNotEnoughChange
			}
		}
	}
	return cells, total, nil
}

var balanceLock sync.Mutex

func (c *Contracts) GetBalanceCellWithLock(p *ParamGetBalanceCells) (uint64, []*indexer.LiveCell, error) {
	balanceLock.Lock()
	defer balanceLock.Unlock()
	if p.SearchOrder == "" {
		p.SearchOrder = indexer.SearchOrderAsc
	}
	if p.CapacityForChange == 0 {
		p.CapacityForChange = common.MinCellOccupiedCkb
	}
	liveCells, total, err := c.GetBalanceCells(p)
	if err != nil {
		return 0, nil, fmt.Errorf("GetBalanceCells err: %s", err.Error())
	}

	var outpoints []string
	for _, v := range liveCells {
		outpoints = append(outpoints, common.OutPointStruct2String(v.OutPoint))
	}
	c.cache.AddOutPoint(outpoints)

	return total - p.CapacityNeed, liveCells, nil
}

func SplitOutputCell(total, base, limit uint64, lockScript, typeScript *types.Script, splitOrder indexer.SearchOrder) ([]*types.CellOutput, error) {
	log.Info("SplitOutputCell:", "total: ", total, "base: ", base, "limit: ", limit)
	formatCell := &types.CellOutput{
		Capacity: base,
		Lock:     lockScript,
		Type:     typeScript,
	}
	realBase := formatCell.OccupiedCapacity(nil) * 1e8
	if total < realBase || base < realBase {
		return nil, fmt.Errorf("total(%d) or base(%d) should not less than real base(%d)", total, base, realBase)
	}
	if base < realBase+common.OneCkb { // 1 ckb for tx fee
		base += common.OneCkb
		formatCell.Capacity = base
	}
	log.Info("realBase:", realBase, base)

	var cellList []*types.CellOutput
	splitTotal := uint64(0)

	if splitOrder == indexer.SearchOrderDesc {
		for i := uint64(0); i < limit && splitTotal+2*base < total; i++ {
			tmp := *formatCell
			cellList = append(cellList, &tmp)
			splitTotal += base
		}
		cellList = append(cellList, &types.CellOutput{
			Capacity: total - splitTotal,
			Lock:     lockScript,
			Type:     typeScript,
		})
	} else {
		cellList = append(cellList, &types.CellOutput{
			Capacity: 0,
			Lock:     lockScript,
			Type:     typeScript,
		})
		for i := uint64(0); i < limit && splitTotal+2*base < total; i++ {
			tmp := *formatCell
			cellList = append(cellList, &tmp)
			splitTotal += base
		}
		cellList[0].Capacity = total - splitTotal
	}

	return cellList, nil
}
