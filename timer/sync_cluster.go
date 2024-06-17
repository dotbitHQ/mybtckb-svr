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

	clusterList := make([]*tables.TableClusterInfo, 0)
	lastCursor := ""
	for {
		res, err := t.contracts.Client().GetCells(t.ctx, searchKey, indexer.SearchOrderDesc, 1000, lastCursor)
		if err != nil {
			return fmt.Errorf("GetCells err:", err.Error())
		}
		log.Info("SporeCluster liveCells:", res.LastCursor, " length: ", len(res.Objects))

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
				args := v.Output.Lock.Args
				index, txHash := common.GetOutpointByargs(args)
				btcOutpoint = common.OutPoint2String(txHash, uint(index))
			} else {
				owner, err = address.ConvertScriptToAddress(t.contracts.Mode(), v.Output.Lock)
				if err != nil {
					return fmt.Errorf("ConvertScriptToAddress err: %s", err.Error())
				}
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
				BtcOutpoint: btcOutpoint,
				AddrType:    addrType,
				BlockNum:    v.BlockNumber,
			})
		}
	}

	if err := t.dbDao.CreateClusterInfo(clusterList); err != nil {
		return fmt.Errorf("CreateClusterInfo err: %s", err.Error())
	}
	return nil
}
