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

func (t *TxTimer) doSyncCluster() error {
	searchKey := &indexer.SearchKey{
		Script: &types.Script{
			CodeHash: t.contracts.SporeClusterType.CodeHash,
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
		log.Warn("empty cluster")
		return nil
	}
	fmt.Println(len(res.Objects))
	clusterList := make([]*tables.TableClusterInfo, 0)
	for _, v := range res.Objects {
		//fmt.Println("cluster id: ", hexutils.BytesToHex(v.Output.Type.Args))
		owner, err := address.ConvertScriptToAddress(t.contracts.Mode(), v.Output.Lock)
		if err != nil {
			return fmt.Errorf("ConvertScriptToAddress err: %s", err.Error())
		}

		clusterId := common.Bytes2Hex(v.Output.Type.Args)
		//fmt.Println(hexutils.BytesToHex(v.OutputData))
		//molecue 解析data字段获取cluster id和 content，content-type字段
		data, err := molecule.ClusterDataFromSlice(v.OutputData, true)
		if err != nil {
			log.Errorf("SporeDataFromSlice err : %s", err.Error())
		}
		name := data.Name().RawData()
		//fmt.Println("cluster name: ", string(name))
		description := data.Description().RawData()
		//fmt.Println("cluster description: ", string(description))

		clusterList = append(clusterList, &tables.TableClusterInfo{
			ClusterId:   clusterId,
			Address:     owner,
			Name:        string(name),
			Description: string(description),
			Outpoint:    common.OutPointStruct2String(v.OutPoint),
			BlockNum:    v.BlockNumber,
		})
	}
	if err := t.dbDao.CreateClusterInfo(clusterList); err != nil {
		return fmt.Errorf("CreateClusterInfo err: %s", err.Error())
	}
	return nil
}
