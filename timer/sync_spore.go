package timer

import (
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/lib/common"
	molecule "mybtckb-svr/lib/molecule/spore"
	"mybtckb-svr/tables"
)

func (t *TxTimer) doSyncSpore() error {
	searchKey := &indexer.SearchKey{
		Script: &types.Script{
			CodeHash: t.contracts.SporeType.CodeHash,
			HashType: types.HashTypeData1,
		},
		ScriptType: indexer.ScriptTypeType,
	}
	sporeList := make([]*tables.TableSpore, 0)
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
			owner := ""
			btcOutpoint := ""
			var addrType uint8
			if v.Output.Lock.CodeHash == t.contracts.RGBPP.CodeHash && v.Output.Lock.HashType == t.contracts.RGBPP.HashType {
				addrType = 1
				//fmt.Println("rgb ++ tx: ", v.OutPoint.TxHash.Hex())
				args := v.Output.Lock.Args
				index, txHash := common.GetOutpointByargs(args)
				//fmt.Println(txHash, "---", index)
				btcOutpoint = common.OutPoint2String(txHash, uint(index))
				//owner = "btc"
				owner, err = t.contracts.GetBtcAddressByOutpoint(index, txHash)
				//if err != nil {
				//	continue
				//	//return fmt.Errorf("GetBtcAddressByOutpoint err %s", err.Error())
				//}
			} else {
				owner, err = address.ConvertScriptToAddress(t.contracts.Mode(), v.Output.Lock)
				if err != nil {
					return fmt.Errorf("ConvertScriptToAddress err: %s", err.Error())
				}
			}
			sporeId := common.Bytes2Hex(v.Output.Type.Args)
			data, err := molecule.SporeDataFromSlice(v.OutputData, true)
			if err != nil {
				return fmt.Errorf("SporeDataFromSlice err :", err.Error())
			}
			contentType := data.ContentType().RawData()
			content := data.Content().RawData()
			clusterId := ""
			if !data.ClusterId().IsNone() {
				s, err := data.ClusterId().IntoBytes()
				if err != nil {
					return fmt.Errorf("data.ClusterId().IntoBytes() err: ", err.Error())
				}
				clusterId = types.BytesToHash(s.RawData()).Hex()
			}
			sporeList = append(sporeList, &tables.TableSpore{
				SporeId:     sporeId,
				Address:     owner,
				ClusterId:   clusterId,
				ContentType: string(contentType),
				Content:     content,
				Outpoint:    common.OutPointStruct2String(v.OutPoint),
				BtcOutpoint: btcOutpoint,
				AddrType:    addrType,
				BlockNum:    v.BlockNumber,
			})
		}
	}
	if len(sporeList) == 0 {
		return nil
	}

	if err := t.dbDao.CreateSporeInfo(sporeList); err != nil {
		return fmt.Errorf("CreateSporeInfo err: %s", err.Error())
	}
	return nil
}
