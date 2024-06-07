package txtool

import (
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/status-im/keycard-go/hexutils"
	"mybtckb-svr/lib/common"
	"mybtckb-svr/lib/contract"
	"mybtckb-svr/lib/math/uint128"
	"mybtckb-svr/lib/txbuilder"
)

type TransferParams struct {
	TokenId        string
	Amount         uint64
	Address        string
	ReceiptAddress string
}

func (txTool *TxTool) Transfer(params *TransferParams) (*txbuilder.TxBuilder, error) {
	txBuilder := txbuilder.NewTxBuilder(txbuilder.WithTxBuilderBase(txTool.txBuilderBase), txbuilder.WithCache(txTool.cache))
	txParams := &txbuilder.BuildTransactionParams{}

	//txParams.CellDeps = append(txParams.CellDeps, txTool.contract.ConfigCell.ToCellDep())

	requestAddress, err := address.Parse(params.Address)
	if err != nil {
		return nil, err
	}
	xudtLiveCells, total, xudtType, err := txTool.GetXudtCells(&ParamGetXudtCells{
		LockScript: requestAddress.Script,
		TokenId:    params.TokenId,
		AmountNeed: params.Amount,
	})
	if err != nil {
		return nil, err
	}
	inputCapacity := uint64(0)
	fmt.Println("xudt live cell: ", len(xudtLiveCells))
	fmt.Println("xudt total amount: ", total)

	//input:sender xudt cell
	for _, v := range xudtLiveCells {
		inputCapacity += v.Output.Capacity
		txParams.Inputs = append(txParams.Inputs, &types.CellInput{
			PreviousOutput: v.OutPoint,
		})
	}
	//output: receipt xudt cell
	receiptAddress, err := address.Parse(params.ReceiptAddress)
	if err != nil {
		return nil, err
	}
	receiptXudt := &types.CellOutput{
		Lock: receiptAddress.Script,
		Type: xudtType,
	}
	fmt.Println("receipt args: ", hexutils.BytesToHex(receiptAddress.Script.Args))
	transferAmount := uint128.From64(params.Amount)
	receiptXudt.Capacity = receiptXudt.OccupiedCapacity(transferAmount.Bytes()) * common.OneCkb
	outputCapacity := receiptXudt.Capacity
	txParams.Outputs = append(txParams.Outputs, receiptXudt)
	txParams.OutputsData = append(txParams.OutputsData, transferAmount.Bytes())
	fmt.Println(hexutils.BytesToHex(txParams.Outputs[0].Lock.Args))
	fmt.Println("999", len(txParams.Outputs))
	//fmt.Println(hexutils.BytesToHex(txParams.Outputs[1].Lock.Args))
	// xudt change fee
	var changeXudt *types.CellOutput
	if total.Cmp(transferAmount) > 0 {
		fmt.Println()
		changeXudt = &types.CellOutput{
			Lock: requestAddress.Script,
			Type: xudtType,
		}
		changeAmount := total.Sub(transferAmount).Bytes()
		fmt.Println("xudt change: ", total.Sub(transferAmount).String())
		changeXudt.Capacity = changeXudt.OccupiedCapacity(changeAmount) * common.OneCkb
		outputCapacity += changeXudt.Capacity
		txParams.Outputs = append(txParams.Outputs, changeXudt)
		txParams.OutputsData = append(txParams.OutputsData, changeAmount)
	}
	fmt.Println("outputCapacity: ", outputCapacity, "inputCapacity: ", inputCapacity)
	capacityNeed := common.OneCkb
	inputChange := uint64(0)
	if outputCapacity > inputCapacity {
		capacityNeed += outputCapacity - inputCapacity
	} else {
		inputChange = inputCapacity - outputCapacity
	}

	// fees
	fmt.Println("capacityNeed: ", capacityNeed, "inputChange: ", inputCapacity)
	change, liveCells, err := txTool.contract.GetBalanceCellWithLock(&contract.ParamGetBalanceCells{
		LockScript:        requestAddress.Script,
		CapacityNeed:      capacityNeed,
		CapacityForChange: common.MinCellOccupiedCkb,
		SearchOrder:       indexer.SearchOrderDesc,
	})
	if err != nil {
		return nil, fmt.Errorf("GetBalanceCells err: %s", err.Error())
	}
	for _, v := range liveCells {
		txParams.Inputs = append(txParams.Inputs, &types.CellInput{
			PreviousOutput: v.OutPoint,
		})
	}
	// change fee
	txParams.Outputs = append(txParams.Outputs, &types.CellOutput{
		Capacity: change + inputChange + common.OneCkb,
		Lock:     requestAddress.Script,
	})
	txParams.OutputsData = append(txParams.OutputsData, []byte{})

	//celldeps
	txParams.CellDeps = append(txParams.CellDeps, &types.CellDep{
		OutPoint: txTool.contract.XudtType.OutPoint,
		DepType:  types.DepTypeCode,
	})

	if err := txBuilder.BuildTransaction(txParams); err != nil {
		return nil, err
	}

	sizeInBlock, _ := txBuilder.SizeInBlock()
	//fee := uint64(math.Ceil(float64(sizeInBlock) / float64(1000)))
	fee := sizeInBlock + 1000
	latestIndex := len(txBuilder.Outputs) - 1
	changeCapacity := txBuilder.Outputs[latestIndex].Capacity - fee
	txBuilder.Outputs[latestIndex].Capacity = changeCapacity
	for _, v := range txBuilder.Outputs {
		fmt.Println("00000: ", hexutils.BytesToHex(v.Lock.Args))
	}
	return txBuilder, nil
}
