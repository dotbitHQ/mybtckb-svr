package txbuilder

import (
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"mybtckb-svr/lib/common"
)

func GenerateMultiSignData(firstN uint8, signatures [][]byte, sortArgsList [][]byte) []byte {
	var lock []byte

	lock = append(lock, 0)

	// first n
	lock = append(lock, firstN)

	// signature num
	sigNum := uint8(len(signatures))
	lock = append(lock, sigNum)

	// address num
	addressNum := uint8(len(sortArgsList))
	lock = append(lock, addressNum)

	// args_of_all_addresses
	for _, v := range sortArgsList {
		lock = append(lock, v...)
	}

	// signatures
	for _, v := range signatures {
		if len(v) > 0 {
			lock = append(lock, v...)
		} else {
			lock = append(lock, common.SignaturePlaceholder...)
		}
	}
	return lock
}

type MultiSign struct {
	Reserved        uint8
	FirstN          uint8
	SignNum         uint8
	AddressNum      uint8
	ArgsOfAddresses [][]byte
	Signatures      [][]byte
}

func DecodeMultiSign(bs []byte) *MultiSign {
	if len(bs) <= 4 {
		return nil
	}
	multiSign := &MultiSign{}
	multiSign.Reserved = bs[0]
	multiSign.FirstN = bs[1]
	multiSign.SignNum = bs[2]
	multiSign.AddressNum = bs[3]

	idx := 4
	addressStep := 20
	for i := 0; i < int(multiSign.AddressNum); i++ {
		if idx+addressStep > len(bs) {
			return nil
		}
		multiSign.ArgsOfAddresses = append(multiSign.ArgsOfAddresses, bs[idx:idx+addressStep])
		idx += addressStep
	}

	signStep := len(common.SignaturePlaceholder)
	for i := 0; i < int(multiSign.SignNum); i++ {
		if idx+signStep > len(bs) {
			return nil
		}
		multiSign.Signatures = append(multiSign.Signatures, bs[idx:idx+signStep])
		idx += signStep
	}
	return multiSign
}

func GenerateMultiSignWitnessArgs(firstN uint8, signatures [][]byte, sortArgsList [][]byte) *types.WitnessArgs {
	lock := GenerateMultiSignData(firstN, signatures, sortArgsList)
	wa := types.WitnessArgs{
		Lock:       lock,
		InputType:  nil,
		OutputType: nil,
	}
	return &wa
}

func (d *TxBuilder) AddMultiSignatureForTx(group []int, firstN uint8, signatures [][]byte, sortArgsList [][]byte) error {
	wa := GenerateMultiSignWitnessArgs(firstN, signatures, sortArgsList)
	wab, _ := wa.Serialize()
	d.Transaction.Witnesses[group[0]] = wab
	return nil
}

func (d *TxBuilder) AddSignatureForTx(reqSignList []SignData) error {
	if reqSignList == nil || len(reqSignList) == 0 {
		return fmt.Errorf("signData is nil")
	}
	groups, err := d.getGroupsFromTx()
	if err != nil {
		return fmt.Errorf("getGroupsFromTx err: %s", err.Error())
	}
	signList, err := d.generateDigestListFromTx(d.ServerSignGroup)
	if err != nil {
		return err
	}
	reqSignMap := make(map[string]string)
	for _, v := range reqSignList {
		reqSignMap[v.SignMsg] = v.Signature
	}

	// [[0,1], [2]]
	for i, signData := range signList {
		if signData.SignMsg == "" {
			continue
		}
		signature := reqSignMap[signData.SignMsg]
		if signature == "" {
			continue
		}
		group := groups[i]

		wa, err := types.DeserializeWitnessArgs(d.Transaction.Witnesses[group[0]])
		if err != nil {
			return err
		}

		wa.Lock = common.Hex2Bytes(signature)
		wab, _ := wa.Serialize()
		d.Transaction.Witnesses[group[0]] = wab
	}
	return nil
}

func (d *TxBuilder) ServerSignTx() error {
	if len(d.ServerSignGroup) == 0 {
		return nil
	}
	if d.HandleServerSign == nil {
		return fmt.Errorf("handleServerSign is nil")
	}
	digest, err := d.generateDigestByGroup(d.ServerSignGroup, []int{})
	if err != nil {
		return fmt.Errorf("generateDigestByGroup err: %s", err.Error())
	}
	sig, err := d.HandleServerSign(digest.SignMsg)
	if err != nil {
		return fmt.Errorf("handleServerSign err: %s", err.Error())
	}

	wa := &types.WitnessArgs{
		Lock:       sig,
		InputType:  nil,
		OutputType: nil,
	}
	wab, _ := wa.Serialize()
	d.Transaction.Witnesses[d.ServerSignGroup[0]] = wab
	return nil
}
