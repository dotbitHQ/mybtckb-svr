package timer

import (
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/lib/common"
	"mybtckb-svr/lib/math/uint128"
	"mybtckb-svr/tables"
)

func (t *TxTimer) doSyncRgbppXudt() error {
	searchKey := &indexer.SearchKey{
		Script: &types.Script{
			CodeHash: t.contracts.RGBPP.CodeHash,
			HashType: t.contracts.RGBPP.HashType,
		},
		ScriptType: indexer.ScriptTypeLock,
		Filter: &indexer.CellsFilter{
			Script: &types.Script{
				CodeHash: t.contracts.XudtType.CodeHash,
				HashType: t.contracts.XudtType.HashType,
			},
		},
	}
	xudtList := make([]*tables.TableXudt, 0)
	lastCursor := ""
	for {
		res, err := t.contracts.Client().GetCells(t.ctx, searchKey, indexer.SearchOrderDesc, 1000, lastCursor)
		if err != nil {
			return fmt.Errorf("GetCells err:", err.Error())
		}
		log.Info("Spore liveCells:", res.LastCursor, " length: ", len(res.Objects))

		if len(res.Objects) == 0 || lastCursor == res.LastCursor {
			break
		}
		lastCursor = res.LastCursor
		for _, v := range res.Objects {
			if v.Output.Lock.CodeHash == t.contracts.RGBPP.CodeHash && v.Output.Lock.HashType == t.contracts.RGBPP.HashType {
				//fmt.Println("rgb ++ tx: ", v.OutPoint.TxHash.Hex())
				args := v.Output.Lock.Args
				index, txHash := common.GetOutpointByargs(args)
				//fmt.Println(txHash, "---", index)
				owner := ""
				btcOutpoint := common.OutPoint2String(txHash, uint(index))

				if addr, err := t.rc.GetBtcAddrCache(txHash, index); err != nil {
					log.Errorf("GetBtcAddrCache err: %s", err.Error())
				} else if addr != "" {
					owner = addr
				}
				tokenId := common.Bytes2Hex(v.Output.Type.Args)
				amount := uint128.FromBytes(v.OutputData)
				xudtList = append(xudtList, &tables.TableXudt{
					Address:     owner,
					TokenId:     tokenId,
					Amount:      amount.String(),
					AddrType:    1,
					Outpoint:    common.OutPointStruct2String(v.OutPoint),
					BtcOutpoint: btcOutpoint,
					BlockNum:    v.BlockNumber,
				})
			}
		}
	}
	if len(xudtList) == 0 {
		return nil
	}
	if err := t.dbDao.CreateRgbppXudt(xudtList); err != nil {
		return fmt.Errorf("CreateSporeInfo err: %s", err.Error())
	}
	return nil
}
