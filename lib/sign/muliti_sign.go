package sign

import (
	"errors"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"mybtckb-svr/lib/common"
)

type MultiSigCfg struct {
	FirstN      uint8    `json:"first_n" yaml:"first_n"`
	SigNum      uint8    `json:"sig_num" yaml:"sig_num"`
	PrivateKey  string   `json:"private_key,omitempty" yaml:"private_key"`
	SortAddress []string `json:"sort_address" yaml:"sort_address"`
}

func (c *MultiSigCfg) Check() error {
	if c.FirstN >= uint8(len(c.SortAddress)) {
		return errors.New("first_n >= len(sort_addresses)")
	}
	if c.SigNum > uint8(len(c.SortAddress)) {
		return errors.New("sig_num > len(sort_addresses)")
	}
	if len(c.SortAddress) == 0 {
		return errors.New("sort_address is empty")
	}
	return nil
}

func genMultiSignArgs(firstN, signNum uint8, sortArgs [][]byte) ([]byte, error) {
	bs := []byte{0, firstN, signNum, uint8(len(sortArgs))}
	for _, v := range sortArgs {
		bs = append(bs, v...)
	}
	hash, _ := blake2b.Blake160(bs)
	return hash, nil
}

func GenOmniLockMultiSignArgs(firstN, signNum uint8, sortArgs [][]byte) ([]byte, error) {
	args, err := genMultiSignArgs(firstN, signNum, sortArgs)
	if err != nil {
		return nil, err
	}
	args = append([]byte{6}, args...)
	args = append(args, 0)
	return args, nil
}

func GenOmniLockMultiSignArgsByAddress(firstN, signNum uint8, address []string) ([]byte, error) {
	sortArgs, err := common.AddressToArgs(address)
	if err != nil {
		return nil, err
	}
	args, err := GenOmniLockMultiSignArgs(firstN, signNum, sortArgs)
	if err != nil {
		return nil, err
	}
	return args, nil
}
