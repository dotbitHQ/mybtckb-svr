package txbuilder

import (
	"encoding/binary"
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/lib/common"
	"sort"
)

type SignData struct {
	SignType  common.AlgorithmId `json:"sign_type"`
	SignMsg   string             `json:"sign_msg"`
	Signature string             `json:"signature,omitempty"`
}

func (d *TxBuilder) GenerateMultiSignDigest(group []int, firstN uint8, signatures [][]byte, sortArgsList [][]byte) ([]byte, error) {
	if len(group) == 0 {
		return nil, fmt.Errorf("group is nil")
	}

	wa := GenerateMultiSignWitnessArgs(firstN, signatures, sortArgsList)
	data, _ := wa.Serialize()
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(data)))

	hash, _ := d.Transaction.ComputeHash()
	message := append(hash.Bytes(), length...)
	message = append(message, data...)

	// hash the other witnesses in the group
	if len(group) > 1 {
		for i := 1; i < len(group); i++ {
			data = d.Transaction.Witnesses[group[i]]
			lengthTmp := make([]byte, 8)
			binary.LittleEndian.PutUint64(lengthTmp, uint64(len(data)))
			message = append(message, lengthTmp...)
			message = append(message, data...)
		}
	}

	// hash witnesses which do not in any input group
	for _, wit := range d.Transaction.Witnesses[len(d.Transaction.Inputs):] {
		lengthTmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(lengthTmp, uint64(len(wit)))
		message = append(message, lengthTmp...)
		message = append(message, wit...)
	}
	return blake2b.Blake256(message)
}

func (d *TxBuilder) GenerateDigestListFromTx(skipGroups []int) ([]SignData, error) {
	fmt.Println("serverSignGroup: ", d.ServerSignGroup)
	skipGroups = append(skipGroups, d.ServerSignGroup...)
	signList, err := d.generateDigestListFromTx(skipGroups)
	if err != nil {
		return nil, err
	}
	res := make([]SignData, 0)
	for _, v := range signList {
		if v.SignMsg != "" {
			res = append(res, v)
		}
	}
	return res, nil
}

func (d *TxBuilder) generateDigestListFromTx(skipGroups []int) ([]SignData, error) {
	fmt.Println("skipGroups: ", skipGroups)
	groups, err := d.getGroupsFromTx()
	if err != nil {
		return nil, fmt.Errorf("getGroupsFromTx err: %s", err.Error())
	}
	log.Info("groups:", len(groups), groups[0])
	var digestList []SignData
	for _, group := range groups {
		digest, err := d.generateDigestByGroup(group, skipGroups)
		if err != nil {
			return nil, fmt.Errorf("generateDigestByGroup err: %s", err.Error())
		}
		fmt.Println(456783, digest.SignType, digest.SignMsg)
		digestList = append(digestList, *digest)
	}
	return digestList, nil
}

func (d *TxBuilder) getGroupsFromTx() ([][]int, error) {
	//input
	//code1-1
	//code1-2
	//code2-1
	var tmpMapForGroup = make(map[string][]int)
	var sortList []string
	for i, v := range d.Transaction.Inputs {
		item, err := d.getInputCell(v.PreviousOutput)
		if err != nil {
			return nil, fmt.Errorf("getInputCell err: %s", err.Error())
		}
		cellHash, _ := item.Cell.Output.Lock.Hash()
		indexList, okTmp := tmpMapForGroup[cellHash.String()]
		if !okTmp {
			sortList = append(sortList, cellHash.String())
		}
		indexList = append(indexList, i)
		tmpMapForGroup[cellHash.String()] = indexList
	}
	//sortList = [code1,code2]
	//tmpMapForGroup = [
	//	code1=>[0,1]
	//	code2=>[2]
	//]
	sort.Strings(sortList)
	var list [][]int
	for _, v := range sortList {
		item, _ := tmpMapForGroup[v]
		list = append(list, item)
	}
	//list = [[0,1], [2]]
	return list, nil
}

func (d *TxBuilder) generateDigestByGroup(group []int, skipGroups []int) (*SignData, error) {
	if group == nil || len(group) < 1 {
		return nil, fmt.Errorf("invalid param")
	}

	signData := &SignData{}

	witnessArgs := &types.WitnessArgs{}
	emptyWitnessArgs := &types.WitnessArgs{}
	emptyWitnessArgsLock := make([]byte, 65)

	groupWitness := d.Transaction.Witnesses[group[0]]
	if len(groupWitness) > 0 {
		decodeWitnessArgs, err := types.DeserializeWitnessArgs(groupWitness)
		if err != nil {
			return nil, err
		}
		emptyWitnessArgs = decodeWitnessArgs
		witnessArgs = decodeWitnessArgs
	}
	witnessArgsData, _ := witnessArgs.Serialize()
	d.Transaction.Witnesses[group[0]] = witnessArgsData

	emptyWitnessArgs.Lock = emptyWitnessArgsLock
	data, _ := emptyWitnessArgs.Serialize()

	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(data)))

	hash, _ := d.Transaction.ComputeHash()
	//fmt.Println("tx_hash:", hash.Hex())

	message := append(hash.Bytes(), length...)
	message = append(message, data...)
	//fmt.Println("init witness:", common.Bytes2Hex(message))
	// hash the other witnesses in the group
	if len(group) > 1 {
		for i := 1; i < len(group); i++ {
			data = d.Transaction.Witnesses[group[i]]
			lengthTmp := make([]byte, 8)
			binary.LittleEndian.PutUint64(lengthTmp, uint64(len(data)))
			message = append(message, lengthTmp...)
			message = append(message, data...)
			//fmt.Println("add group other witness:", common.Bytes2Hex(message))
		}
	}
	//fmt.Println("add group other witness:", common.Bytes2Hex(message))
	// hash witnesses which do not in any input group
	for _, wit := range d.Transaction.Witnesses[len(d.Transaction.Inputs):] {
		lengthTmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(lengthTmp, uint64(len(wit)))
		message = append(message, lengthTmp...)
		message = append(message, wit...)
	}
	//fmt.Println("add other witness:", common.Bytes2Hex(message))

	message, _ = blake2b.Blake256(message)
	signData.SignMsg = common.Bytes2Hex(message)
	log.Info("digest:", signData.SignMsg)

	// skip useless signature
	if len(skipGroups) != 0 {
		skip := false
		for i := range group {
			for j := range skipGroups {
				if group[i] == skipGroups[j] {
					skip = true
					break
				}
			}
			if skip {
				break
			}
		}
		if skip {
			signData.SignMsg = ""
		}
	}
	return signData, nil
}

func (d *TxBuilder) getInputCell(o *types.OutPoint) (*types.CellWithStatus, error) {
	if o == nil {
		return nil, fmt.Errorf("OutPoint is nil")
	}
	key := fmt.Sprintf("%s-%d", o.TxHash.Hex(), o.Index)
	if item, ok := d.MapInputsCell[key]; ok {
		if item.Cell != nil && item.Cell.Output != nil && item.Cell.Output.Lock != nil {
			return item, nil
		}
	}
	if cell, err := d.Contracts.Client().GetLiveCell(d.Ctx, o, true); err != nil {
		return nil, fmt.Errorf("GetLiveCell err: %s", err.Error())
	} else if cell.Cell.Output.Lock != nil {
		d.MapInputsCell[key] = cell
		return cell, nil
	} else {
		log.Warn("GetLiveCell:", key, cell.Status)
		if !d.notCheckInputs {
			return nil, fmt.Errorf("cell [%s] not live", key)
		} else {
			d.MapInputsCell[key] = cell
			return cell, nil
		}
	}
}
