package txbuilder

import (
	"errors"
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/nervosnetwork/ckb-sdk-go/utils"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/lib/common"
)

func EqualArgs(src, dst string) bool {
	if common.Has0xPrefix(src) {
		src = src[2:]
	}
	if common.Has0xPrefix(dst) {
		dst = dst[2:]
	}
	return src == dst
}

func (d *TxBuilder) newTx() error {
	systemScriptCell, err := utils.NewSystemScripts(d.Contracts.Client())
	if err != nil {
		return err
	}
	baseTx := transaction.NewSecp256k1SingleSigTx(systemScriptCell)
	d.Transaction = baseTx
	return nil
}

func (d *TxBuilder) addInputsForTx(inputs []*types.CellInput) error {
	if len(inputs) == 0 {
		return fmt.Errorf("inputs is nil")
	}
	startIndex := len(d.Transaction.Inputs)
	_, _, err := d.addInputsForTransaction(d.Transaction, inputs)
	if err != nil {
		return fmt.Errorf("AddInputsForTransaction err: %s", err.Error())
	}

	systemScriptCell, err := utils.NewSystemScripts(d.Contracts.Client())
	if err != nil {
		return err
	}

	var cellDepList []*types.CellDep
	for i, v := range inputs {
		if v == nil {
			return fmt.Errorf("input is nil")
		}
		item, err := d.getInputCell(v.PreviousOutput)
		if err != nil {
			return fmt.Errorf("getInputCell err: %s", err.Error())
		}

		if item.Cell.Output.Lock != nil &&
			item.Cell.Output.Lock.CodeHash.Hex() == systemScriptCell.SecpSingleSigCell.CellHash.Hex() &&
			EqualArgs(common.Bytes2Hex(item.Cell.Output.Lock.Args), d.ServerArgs) {
			d.ServerSignGroup = append(d.ServerSignGroup, startIndex+i)
		}

	}
	d.addCellDepListIntoMapCellDep(cellDepList)
	return nil
}

func (d *TxBuilder) addInputsForTransaction(transaction *types.Transaction, inputs []*types.CellInput) ([]int, *types.WitnessArgs, error) {
	if len(inputs) == 0 {
		return nil, nil, errors.New("input cells empty")
	}
	group := make([]int, len(inputs))
	start := len(transaction.Inputs)
	for i := 0; i < len(inputs); i++ {
		input := inputs[i]
		transaction.Inputs = append(transaction.Inputs, input)
		transaction.Witnesses = append(transaction.Witnesses, []byte{})
		group[i] = start + i
	}
	witnessArgs := &types.WitnessArgs{Lock: common.SignaturePlaceholder}
	witnessArgsData, _ := witnessArgs.Serialize()
	transaction.Witnesses[start] = witnessArgsData
	return group, witnessArgs, nil
}

func (d *TxBuilder) addOutputsForTx(outputs []*types.CellOutput, outputsData [][]byte) error {
	lo := len(outputs)
	lod := len(outputsData)
	if lo == 0 || lod == 0 || lo != lod {
		return fmt.Errorf("outputs[%d] or outputDatas[%d]", lo, lod)
	}

	var cellDepList []*types.CellDep
	for i := 0; i < lo; i++ {
		output := outputs[i]
		outputData := outputsData[i]
		d.Transaction.Outputs = append(d.Transaction.Outputs, output)
		d.Transaction.OutputsData = append(d.Transaction.OutputsData, outputData)
		if output.Type == nil {
			continue
		}

	}

	d.addCellDepListIntoMapCellDep(cellDepList)
	return nil
}

func (d *TxBuilder) checkTxWitnesses() error {
	if len(d.Transaction.Witnesses) == 0 {
		return fmt.Errorf("witness is nil")
	}
	lenI := len(d.Transaction.Inputs)
	lenW := len(d.Transaction.Witnesses)
	if lenW < lenI {
		return fmt.Errorf("len witness[%d]<len inputs[%d]", lenW, lenI)
	}
	return nil
}

func (d *TxBuilder) addCellDepListIntoMapCellDep(cellDepList []*types.CellDep) {
	for i, v := range cellDepList {
		k := fmt.Sprintf("%s-%d", v.OutPoint.TxHash.Hex(), v.OutPoint.Index)
		d.mapCellDep[k] = cellDepList[i]
	}
}

func (d *TxBuilder) addMapCellDepWitnessForBaseTx(cellDepList []*types.CellDep) error {

	//xudt, ok := d.Contracts.GetSoScriptByName(contract.ConfigNameXudtTypeArgs)
	//if ok {
	//	cellDepList = append(cellDepList, xudt.ToCellDep())
	//}
	//xudtExt, ok := d.Contracts.GetSoScriptByName(contract.ConfigNameXudtExtensionTypeArgs)
	//if ok {
	//	cellDepList = append(cellDepList, xudtExt.ToCellDep())
	//}
	//xudtOwner, ok := d.Contracts.GetSoScriptByName(contract.ConfigNameXudtOwnerTypeArgs)
	//if ok {
	//	cellDepList = append(cellDepList, xudtOwner.ToCellDep())
	//}
	//xudtInfo, ok := d.Contracts.GetSoScriptByName(contract.ConfigNameXudtInfoTypeArgs)
	//if ok {
	//	cellDepList = append(cellDepList, xudtInfo.ToCellDep())
	//}

	tmpMap := make(map[string]bool)
	var tmpCellDeps []*types.CellDep
	for _, v := range cellDepList {
		k := fmt.Sprintf("%s-%d", v.OutPoint.TxHash.Hex(), v.OutPoint.Index)
		if _, ok := tmpMap[k]; ok {
			continue
		}
		tmpMap[k] = true
		tmpCellDeps = append(tmpCellDeps, &types.CellDep{
			OutPoint: v.OutPoint,
			DepType:  v.DepType,
		})
	}
	if len(tmpCellDeps) > 0 {
		d.Transaction.CellDeps = append(tmpCellDeps, d.Transaction.CellDeps...)
	}
	for k, v := range d.mapCellDep {
		if _, ok := tmpMap[k]; ok {
			continue
		}
		d.Transaction.CellDeps = append(d.Transaction.CellDeps, &types.CellDep{
			OutPoint: v.OutPoint,
			DepType:  v.DepType,
		})
	}

	if len(d.otherWitnesses) > 0 {
		d.Transaction.Witnesses = append(d.Transaction.Witnesses, d.otherWitnesses...)
	}
	return nil
}

func (d *TxBuilder) SendTransactionWithCheck(needCheck bool) (*types.Hash, error) {
	if needCheck {
		err := d.checkTxBeforeSend()
		if err != nil {
			return nil, fmt.Errorf("checkTxBeforeSend err: %s", err.Error())
		}
	}

	err := d.ServerSignTx()
	if err != nil {
		return nil, fmt.Errorf("remoteSignTx err: %s", err.Error())
	}

	fmt.Printf("before send: %s\n", d.TxString())
	txHash, err := d.Contracts.Client().SendTransaction(d.Ctx, d.Transaction)
	if err != nil {
		return nil, fmt.Errorf("SendTransaction err: %v", err)
	}
	log.WithField("tx_hash", txHash.Hex()).Info("SendTransaction success")
	return txHash, nil
}

func (d *TxBuilder) SendTransaction() (*types.Hash, error) {
	return d.SendTransactionWithCheck(true)
}

func (d *TxBuilder) checkTxBeforeSend() error {
	// check total num of inputs and outputs
	if len(d.Transaction.Inputs)+len(d.Transaction.Outputs) > 9000 {
		return fmt.Errorf("checkTxBeforeSend, failed len of inputs: %d, ouputs: %d", len(d.Transaction.Inputs), len(d.Transaction.Outputs))
	}
	// check tx fee < 1 CKB
	totalCapacityFromInputs, err := d.getCapacityFromInputs()
	if err != nil {
		return err
	}
	totalCapacityFromOutputs := d.Transaction.OutputsCapacity()
	txFee := totalCapacityFromInputs - totalCapacityFromOutputs
	log.Info("checkTxBeforeSend:", totalCapacityFromInputs, totalCapacityFromOutputs, txFee)
	if totalCapacityFromInputs <= totalCapacityFromOutputs || txFee >= common.OneCkb {
		return fmt.Errorf("checkTxBeforeSend failed, txFee: %d totalCapacityFromInputs: %d totalCapacityFromOutputs: %d", txFee, totalCapacityFromInputs, totalCapacityFromOutputs)
	}

	// check witness format
	err = d.checkTxWitnesses()
	if err != nil {
		return err
	}
	// check the occupied capacity
	for i, cell := range d.Transaction.Outputs {
		occupied := cell.OccupiedCapacity(d.Transaction.OutputsData[i]) * common.OneCkb
		if cell.Capacity < occupied {
			return fmt.Errorf("checkTxBeforeSend occupied capacity failed, occupied: %d capacity: %d index: %d", occupied, cell.Capacity, i)
		}
	}
	log.Info("check success before sent")
	return nil
}

func (d *TxBuilder) getCapacityFromInputs() (uint64, error) {
	total := uint64(0)
	for _, v := range d.Transaction.Inputs {
		item, err := d.getInputCell(v.PreviousOutput)
		if err != nil {
			return 0, fmt.Errorf("getInputCell err: %s", err.Error())
		}
		total += item.Cell.Output.Capacity
	}
	return total, nil
}
