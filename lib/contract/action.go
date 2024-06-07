package contract

import (
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type Action = string

const (
	ActionMintXudt      Action = "mint_xudt"
	ActionTransferXudt  Action = "transfer_xudt"
	ActionTransferSpore Action = "transfer_spore"
)

type CellData struct {
	Output     *types.CellOutput
	Index      uint
	OutputData []byte
}
type CellNum struct {
	Xudt         []*CellData
	Unique       []*CellData
	Spore        []*CellData
	SporeCluster []*CellData
}
type TxTypeData struct {
	InputCell, OutputCell CellNum
}

func (c *Contracts) CheckCellType(index uint, output *types.CellOutput, outputData []byte, cellNum *CellNum) {
	if output.Type == nil {
		return
	}
	//fmt.Println("code hash: ", output.Type.CodeHash)
	if output.Type.CodeHash == c.XudtType.CodeHash && output.Type.HashType == c.XudtType.HashType {
		cellNum.Xudt = append(cellNum.Xudt, &CellData{
			Output:     output,
			Index:      index,
			OutputData: outputData,
		})
	} else if output.Type.CodeHash == c.UniqueType.CodeHash && output.Type.HashType == c.UniqueType.HashType {
		cellNum.Unique = append(cellNum.Unique, &CellData{
			Output:     output,
			Index:      index,
			OutputData: outputData,
		})
	} else if output.Type.CodeHash == c.SporeType.CodeHash && output.Type.HashType == c.SporeType.HashType {
		cellNum.Spore = append(cellNum.Spore, &CellData{
			Output:     output,
			Index:      index,
			OutputData: outputData,
		})
	} else if output.Type.CodeHash == c.SporeClusterType.CodeHash && output.Type.HashType == c.SporeClusterType.HashType {
		cellNum.SporeCluster = append(cellNum.SporeCluster, &CellData{
			Output:     output,
			Index:      index,
			OutputData: outputData,
		})
	}
}

func (c *Contracts) TxToDidCellAction(tx *types.Transaction) (Action, TxTypeData, error) {
	var txTypeData TxTypeData
	//fmt.Println("tx hash: ", tx.Hash.Hex())
	for k, v := range tx.Outputs {
		c.CheckCellType(uint(k), v, tx.OutputsData[k], &txTypeData.OutputCell)
	}
	if len(txTypeData.OutputCell.Xudt) == 0 && len(txTypeData.OutputCell.Spore) == 0 && len(txTypeData.OutputCell.SporeCluster) == 0 {
		return "", txTypeData, nil
	}
	for _, v := range tx.Inputs {
		res, err := c.client.GetTransaction(c.ctx, v.PreviousOutput.TxHash)
		if err != nil {
			return "", txTypeData, fmt.Errorf("GetTransaction err: ", err.Error())
		}
		c.CheckCellType(v.PreviousOutput.Index, res.Transaction.Outputs[v.PreviousOutput.Index], []byte{}, &txTypeData.InputCell)
	}
	//mint xudt

	//transfer xudt
	if len(txTypeData.InputCell.Xudt) == 0 && len(txTypeData.OutputCell.Xudt) > 0 {
		return ActionMintXudt, txTypeData, nil
	} else if len(txTypeData.InputCell.Xudt) > 0 && len(txTypeData.OutputCell.Xudt) > 0 {
		return ActionTransferXudt, txTypeData, nil
	}
	return "", txTypeData, nil
}

func (c *Contracts) ParseXudtTx() {

}
