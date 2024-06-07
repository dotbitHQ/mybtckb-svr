package txtool

import (
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"mybtckb-svr/lib/common"
	"mybtckb-svr/lib/contract"
	"mybtckb-svr/lib/txbuilder"
)

type TransferCkbParams struct {
	Address        string
	Amount         uint64
	ReceiptAddress string
}

func (txTool *TxTool) TransferCkb(params *TransferCkbParams) (*txbuilder.TxBuilder, error) {
	txBuilder := txbuilder.NewTxBuilder(txbuilder.WithTxBuilderBase(txTool.txBuilderBase), txbuilder.WithCache(txTool.cache))
	txParams := &txbuilder.BuildTransactionParams{}

	//txParams.CellDeps = append(txParams.CellDeps, txTool.contract.ConfigCell.ToCellDep())

	fmt.Println("capacityNeed: ", params.Amount)
	senderAddress, err := address.Parse(params.Address)
	if err != nil {
		return nil, err
	}
	change, liveCells, err := txTool.contract.GetBalanceCellWithLock(&contract.ParamGetBalanceCells{
		LockScript:        senderAddress.Script,
		CapacityNeed:      params.Amount,
		CapacityForChange: common.MinCellOccupiedCkb,
		SearchOrder:       indexer.SearchOrderDesc,
	})
	if err != nil {
		return nil, fmt.Errorf("GetBalanceCells err: %s", err.Error())
	}

	//input
	for _, v := range liveCells {
		txParams.Inputs = append(txParams.Inputs, &types.CellInput{
			PreviousOutput: v.OutPoint,
		})
	}

	receiptAddress, err := address.Parse(params.ReceiptAddress)
	if err != nil {
		return nil, err
	}
	//output
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: params.Amount,
		Lock:     receiptAddress.Script,
	})
	txParams.OutputsData = append(txParams.OutputsData, []byte{})
	//change
	if change > 0 {
		txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
			Capacity: change,
			Lock:     senderAddress.Script,
		})
		txParams.OutputsData = append(txParams.OutputsData, []byte{})
	}

	if err := txBuilder.BuildTransaction(txParams); err != nil {
		return nil, err
	}

	sizeInBlock, _ := txBuilder.SizeInBlock()
	txBuilder.Outputs[1].Capacity -= sizeInBlock + 1000
	return txBuilder, nil
}
