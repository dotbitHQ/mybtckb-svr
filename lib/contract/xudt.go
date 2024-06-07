package contract

import (
	"fmt"
	"mybtckb-svr/lib/molecule/xudt"
)

type XudtInfo struct {
	Decimal uint8  `json:"decimal"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
}

// uniquecell
func ParseXudtInfo(outputData []byte) (xudtInfo XudtInfo, err error) {
	decimal, err := xudt.Bytes2GoU8(outputData[0:1])
	if err != nil {
		err = fmt.Errorf("xudt.Bytes2GoU8 err: %s", err.Error())
		return
	}
	xudtInfo.Decimal = decimal

	nameLen, err := xudt.Bytes2GoU8(outputData[1:2])
	if err != nil {
		err = fmt.Errorf("xudt.Bytes2GoU8 err: %s", err.Error())
		return
	}
	startIdx := uint8(2)
	xudtInfo.Name = string(outputData[startIdx : startIdx+nameLen])
	startIdx += nameLen

	symbolLen, err := xudt.Bytes2GoU8(outputData[startIdx : startIdx+1])
	if err != nil {
		err = fmt.Errorf("xudt.Bytes2GoU8 err: %s", err.Error())
		return
	}
	startIdx += 1
	xudtInfo.Symbol = string(outputData[startIdx : startIdx+symbolLen])
	return
}
