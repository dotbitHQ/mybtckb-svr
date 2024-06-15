package contract

import (
	"fmt"
	"mybtckb-svr/lib/common"
	"mybtckb-svr/lib/molecule/xudt"
)

type XudtInfo struct {
	Decimal uint8  `json:"decimal"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
}

// uniquecell
func ParseXudtInfo(outputData []byte) (xudtInfo XudtInfo, err error) {
	fmt.Println("outputData: ", common.Bytes2Hex(outputData))
	if len(outputData) < 4 {
		err = fmt.Errorf("xudt info output length err")
		return
	}
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
	if uint8(len(outputData)) < startIdx+nameLen {
		err = fmt.Errorf("xudt info output length err")
		return
	}
	xudtInfo.Name = string(outputData[startIdx : startIdx+nameLen])
	startIdx += nameLen

	if uint8(len(outputData)) < startIdx+1 {
		err = fmt.Errorf("xudt info output length err")
		return
	}
	symbolLen, err := xudt.Bytes2GoU8(outputData[startIdx : startIdx+1])
	if err != nil {
		err = fmt.Errorf("xudt.Bytes2GoU8 err: %s", err.Error())
		return
	}
	startIdx += 1
	if uint8(len(outputData)) < startIdx+symbolLen {
		err = fmt.Errorf("xudt info output length err")
		return
	}
	xudtInfo.Symbol = string(outputData[startIdx : startIdx+symbolLen])
	return
}
