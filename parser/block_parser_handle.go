package parser

import (
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"mybtckb-svr/dao"
	"mybtckb-svr/lib/contract"
)

func (b *BlockParser) registerTransactionHandle() {
	b.mapTransactionHandle = make(map[string]FuncTransactionHandle)
	b.mapTransactionHandle[contract.ActionMintXudt] = b.MintXudt

}

type FuncTransactionHandleReq struct {
	DbDao          *dao.DbDao
	Tx             *types.Transaction
	TxTypeData     *contract.TxTypeData
	TxHash         string
	BlockNumber    uint64
	BlockIndex     int
	BlockTimestamp uint64
	Action         string
}

type FuncTransactionHandle func(FuncTransactionHandleReq) error
