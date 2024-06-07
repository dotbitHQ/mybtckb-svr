package common

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/tron-us/go-common/crypto"
)

type NetType = int

const (
	NetTypeMainNet  NetType = 1
	NetTypeTestnet2 NetType = 2
	NetTypeTestnet3 NetType = 3
)

type AlgorithmId int

const (
	AlgorithmIdCkb              AlgorithmId = 0
	AlgorithmIdCkbMulti         AlgorithmId = 1
	AlgorithmIdCkbSingle        AlgorithmId = 2
	AlgorithmIdCkbOmniLockMulti AlgorithmId = 3
)

func (d AlgorithmId) ToCoinType() CoinType {
	switch d {
	case AlgorithmIdCkb, AlgorithmIdCkbMulti, AlgorithmIdCkbSingle, AlgorithmIdCkbOmniLockMulti:
		return CoinTypeCKB
	default:
		return ""
	}
}

func (d AlgorithmId) Bytes() []byte {
	return []byte{uint8(d)}
}

const (
	OneCkb             = uint64(1e8)
	MinCellOccupiedCkb = uint64(61 * 1e8)
	PercentRateBase    = 1e4
	UsdRateBase        = 1e6
	UserCellTxFeeLimit = 1e4

	OneYearSec = int64(3600 * 24 * 365)

	BlackHoleAddress = "0x0000000000000000000000000000000000000000"
)

func Has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func Hex2Bytes(s string) []byte {
	if Has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	h, _ := hex.DecodeString(s)
	return h
}

func Bytes2Hex(b []byte) string {
	h := hex.EncodeToString(b)
	if len(h) == 0 {
		h = "0"
	}
	return "0x" + h
}

func GetScript(codeHash, args string) *types.Script {
	return &types.Script{
		CodeHash: types.HexToHash(codeHash),
		HashType: types.HashTypeType,
		Args:     Hex2Bytes(args),
	}
}

func TronHexToBase58(address string) (string, error) {
	tAddr, err := crypto.Encode58Check(&address)
	if err != nil {
		return "", fmt.Errorf("Encode58Check:%v", err)
	}
	return *tAddr, nil
}

func TronBase58ToHex(address string) (string, error) {
	addr, err := crypto.Decode58Check(&address)
	if err != nil {
		return "", fmt.Errorf("Decode58Check:%v", err)
	}
	return *addr, nil
}

func Base58CheckDecode(addr string, version byte) (string, error) {
	payload, v, err := base58.CheckDecode(addr)
	if err != nil {
		return "", fmt.Errorf("base58.CheckDecode err: %s[%s]", err.Error(), addr)
	} else if v != version {
		return "", fmt.Errorf("base58.CheckDecode version diff: %d[%s]", v, addr)
	}
	return hex.EncodeToString(payload), nil
}

func Base58CheckEncode(payload string, version byte) (string, error) {
	bys, err := hex.DecodeString(payload)
	if err != nil {
		return "", fmt.Errorf("payload DecodeString err: %s", err.Error())
	}
	return base58.CheckEncode(bys, version), nil
}

// GetNormalLockScript normal script
func GetNormalLockScript(args string) *types.Script {
	return GetScript(transaction.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH, args)
}

// GetNormalLockScriptByMultiSig multi sig
func GetNormalLockScriptByMultiSig(args string) *types.Script {
	return GetScript(transaction.SECP256K1_BLAKE160_MULTISIG_ALL_TYPE_HASH, args)
}
