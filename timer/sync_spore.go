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
	res, err := t.contracts.Client().GetCells(t.ctx, searchKey, indexer.SearchOrderDesc, 1000, "")
	if err != nil {
		return fmt.Errorf("GetCells err:", err.Error())
	}
	//发行xudt FAIR 0x79fa45b10a34cf767e9be38394ba62eb551988b3c8363b7d5ac5422af1071f13
	//transfer xudt TEST 0x89230a87dc2b66b44a5dd8d3ef9e2bdbddfb776f5ec754298e76c11971ad7cad
	if len(res.Objects) == 0 {
		log.Warn("empty spore")
		return nil
	}
	fmt.Println(len(res.Objects))
	sporeList := make([]*tables.TableSpore, 0)
	for _, v := range res.Objects {
		owner, err := address.ConvertScriptToAddress(t.contracts.Mode(), v.Output.Lock)
		if err != nil {
			return fmt.Errorf("ConvertScriptToAddress err: %s", err.Error())
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
			BlockNum:    v.BlockNumber,
		})
	}
	if err := t.dbDao.CreateSporeInfo(sporeList); err != nil {
		return fmt.Errorf("CreateSporeInfo err: %s", err.Error())
	}
	return nil
}
