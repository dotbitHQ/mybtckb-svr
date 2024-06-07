package sign

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/secp256k1"
	log "github.com/sirupsen/logrus"
	"strings"
)

func LocalSign(privateKey string) HandleSignCkbMessage {
	if strings.HasPrefix(privateKey, "0x") {
		privateKey = strings.TrimPrefix(privateKey, "0x")
	}
	return func(message string) ([]byte, error) {
		log.Info("LocalSign:", message)
		bys := common.Hex2Bytes(message)
		key, err := secp256k1.HexToKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("secp256k1.HexToKey err: %s", err.Error())
		}
		signed, err := key.Sign(bys)
		if err != nil {
			return nil, err
		}
		return signed, nil
	}
}
