package txtool

import (
	"mybtckb-svr/dao"
	"mybtckb-svr/lib/cache"
	"mybtckb-svr/lib/contract"
	"mybtckb-svr/lib/txbuilder"
)

type TxTool struct {
	cache         *cache.Cache
	dbDao         *dao.DbDao
	contract      *contract.Contracts
	txBuilderBase *txbuilder.TxBuilderBase
}

type ToolOpts func(*TxTool)

func WithCache(c *cache.Cache) ToolOpts {
	return func(tool *TxTool) {
		tool.cache = c
	}
}

func WithDb(db *dao.DbDao) ToolOpts {
	return func(tool *TxTool) {
		tool.dbDao = db
	}
}

func WithContract(contract *contract.Contracts) ToolOpts {
	return func(tool *TxTool) {
		tool.contract = contract
	}
}

func WithTxBuilderBase(txBuilderBase *txbuilder.TxBuilderBase) ToolOpts {
	return func(tool *TxTool) {
		tool.txBuilderBase = txBuilderBase
	}
}

func NewTxTool(opts ...ToolOpts) *TxTool {
	txTool := &TxTool{}
	for _, opt := range opts {
		opt(txTool)
	}
	return txTool
}

func (txTool *TxTool) TxBuilderBase() *txbuilder.TxBuilderBase {
	return txTool.txBuilderBase
}

func (txTool *TxTool) Cache() *cache.Cache {
	return txTool.cache
}
