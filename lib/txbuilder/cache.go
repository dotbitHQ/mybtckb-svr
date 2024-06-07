package txbuilder

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"math/rand"
	"time"
)

var SignCacheNotFound = errors.New("sign cache not found")

const (
	signListCacheKey = "sign:tx:%s"
)

type SignInfoList struct {
	SignKey  string     `json:"sign_key" binding:"required"`
	Action   string     `json:"action" binding:"required"`
	SignList []SignData `json:"sign_list"`
}

type SignInfoCache struct {
	Action          string              `json:"action"`
	TransferXudt    *TransferXudtCache  `json:"transfer_xudt"`
	TransferSpore   *TransferSporeCache `json:"transfer_spore"`
	TransactionBase *TransactionBase    `json:"builder_tx"`
}

func (s *SignInfoCache) SignKey() string {
	key := fmt.Sprintf("%s%d%d", s.Action, rand.Int(), time.Now().UnixNano())
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}

type TransferXudtCache struct {
	Address        string `json:"address"`
	ReceiptAddress string `json:"receipt_address"`
	Amount         uint64 `json:"amount"`
	TokenId        string `json:"token_id"`
}

type TransferSporeCache struct {
	Address        string `json:"address"`
	ReceiptAddress string `json:"receipt_address"`
	SporeId        string `json:"spore_id"`
}

func (d *TxBuilder) GenSignInfo(sic *SignInfoCache) (*SignInfoList, error) {

	sic.TransactionBase = d.TransactionBase
	d.TxString()
	signKey := sic.SignKey()
	fmt.Println("key---", fmt.Sprintf(signListCacheKey, signKey))
	d.cache.Set(fmt.Sprintf(signListCacheKey, signKey), gconv.String(sic), time.Minute*10)

	signListInfo := &SignInfoList{
		Action:  sic.Action,
		SignKey: signKey,
	}
	return signListInfo, nil
}

func (d *TxBuilder) GetTxBuilderFromCache(signKey string) (sic *SignInfoCache, err error) {
	v, ok := d.cache.Get(fmt.Sprintf(signListCacheKey, signKey))
	if !ok {
		err = SignCacheNotFound
		return sic, err
	}
	fmt.Println("redis cache: ", v)
	if err := json.Unmarshal([]byte(fmt.Sprint(v)), &sic); err != nil {
		fmt.Println("333333 ", err.Error())
		return sic, err
	}
	d.TransactionBase = sic.TransactionBase
	return sic, err
}
