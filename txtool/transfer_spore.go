package txtool

import (
	"context"
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"mybtckb-svr/lib/common"
	"mybtckb-svr/lib/txbuilder"
)

type TransferSporeParams struct {
	SporeId        string
	ReceiptAddress string
}

func (txTool *TxTool) TransferSpore(params *TransferSporeParams) (*txbuilder.TxBuilder, error) {
	txBuilder := txbuilder.NewTxBuilder(txbuilder.WithTxBuilderBase(txTool.txBuilderBase), txbuilder.WithCache(txTool.cache))
	txParams := &txbuilder.BuildTransactionParams{}

	//txParams.CellDeps = append(txParams.CellDeps, txTool.contract.ConfigCell.ToCellDep())

	sporeLiveCell, err := txTool.GetSporeCell(common.Hex2Bytes(params.SporeId))
	if err != nil {
		return nil, err
	}
	txParams.Inputs = append(txParams.Inputs, &types.CellInput{
		PreviousOutput: sporeLiveCell.OutPoint,
	})

	//output: receipt xudt cell
	receiptAddress, err := address.Parse(params.ReceiptAddress)
	if err != nil {
		return nil, err
	}
	outputSporeCell := &types.CellOutput{
		Lock:     receiptAddress.Script,
		Type:     sporeLiveCell.Output.Type,
		Capacity: sporeLiveCell.Output.Capacity,
	}
	txParams.Outputs = append(txParams.Outputs, outputSporeCell)
	txParams.OutputsData = append(txParams.OutputsData, sporeLiveCell.OutputData)

	//celldeps
	txParams.CellDeps = append(txParams.CellDeps, &types.CellDep{
		OutPoint: txTool.contract.SporeType.OutPoint,
		DepType:  types.DepTypeCode,
	})

	if err := txBuilder.BuildTransaction(txParams); err != nil {
		return nil, err
	}

	sizeInBlock, _ := txBuilder.SizeInBlock()
	txBuilder.Outputs[0].Capacity -= sizeInBlock + 1000
	return txBuilder, nil
}

func (TxTool *TxTool) GetSporeCell(args []byte) (*indexer.LiveCell, error) {
	searchKey := &indexer.SearchKey{
		Script: &types.Script{
			CodeHash: TxTool.contract.SporeType.CodeHash,
			HashType: types.HashTypeData1,
			Args:     args,
		},
		ScriptType: indexer.ScriptTypeType,
	}
	res, err := TxTool.contract.Client().GetCells(context.Background(), searchKey, indexer.SearchOrderDesc, 1, "")
	if err != nil {
		return nil, fmt.Errorf("GetCells err: %s", err.Error())
	}
	return res.Objects[0], nil
}
