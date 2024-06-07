package common

import (
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/secp256k1"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
)

var SignaturePlaceholder = make([]byte, 65)

func GetLockFromPk(privateKey string) (*types.Script, error) {
	if strings.HasPrefix(privateKey, "0x") {
		privateKey = strings.TrimPrefix(privateKey, "0x")
	}
	key, err := secp256k1.HexToKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("secp256k1.HexToKey err: %s", err.Error())
	}
	pubKey, _ := blake2b.Blake160(key.PubKey())
	normalLock := GetNormalLockScript(Bytes2Hex(pubKey))
	return normalLock, nil
}

func GetAddressByPk(mode address.Mode, pk string) (string, error) {
	lock, err := GetLockFromPk(pk)
	if err != nil {
		return "", err
	}
	addr, err := address.ConvertScriptToAddress(mode, lock)
	if err != nil {
		return "", err
	}
	return addr, nil
}

func AddressToArgs(addrStr []string) ([][]byte, error) {
	args := make([][]byte, 0)
	for _, v := range addrStr {
		addr, err := address.Parse(v)
		if err != nil {
			return nil, err
		}
		args = append(args, addr.Script.Args)
	}
	return args, nil
}

func ArgsToAddress(mode address.Mode, args [][]byte) ([]string, error) {
	addresses := make([]string, 0)
	for _, v := range args {
		normalLock := GetNormalLockScript(Bytes2Hex(v))
		addr, err := address.ConvertScriptToAddress(mode, normalLock)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}
	return addresses, nil
}
