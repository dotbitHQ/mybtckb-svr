package txbuilder

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"mybtckb-svr/lib/cache"
	"mybtckb-svr/lib/contract"
	"mybtckb-svr/lib/sign"
	"sync"
)

type TxBuilder struct {
	*TxBuilderBase
	*TransactionBase

	ctx            context.Context
	wg             *sync.WaitGroup
	cache          *cache.Cache
	mapCellDep     map[string]*types.CellDep // for memory
	notCheckInputs bool
	otherWitnesses [][]byte
}

type TxBuilderBase struct {
	Ctx              context.Context
	Contracts        *contract.Contracts
	HandleServerSign sign.HandleSignCkbMessage
	ServerArgs       string
}

type TransactionBase struct {
	*types.Transaction `json:"transaction"`
	MapInputsCell      map[string]*types.CellWithStatus `json:"map_inputs_cell"`
	ServerSignGroup    []int                            `json:"server_sign_group"`
}

type BuildTransactionParams struct {
	CellDeps       []*types.CellDep    `json:"cell_deps"`
	HeadDeps       []types.Hash        `json:"head_deps"`
	Inputs         []*types.CellInput  `json:"inputs"`
	Outputs        []*types.CellOutput `json:"outputs"`
	OutputsData    [][]byte            `json:"outputs_data"`
	Witnesses      [][]byte            `json:"witnesses"`
	OtherWitnesses [][]byte            `json:"other_witnesses"`
}

type Options func(*TxBuilder)

func WithTxBuilderBase(base *TxBuilderBase) Options {
	return func(tx *TxBuilder) {
		tx.TxBuilderBase = base
	}
}

func WithTxTransactionBase(txBase *TransactionBase) Options {
	return func(tx *TxBuilder) {
		tx.TransactionBase = txBase
	}
}

func WithTxTransaction(txBase *types.Transaction) Options {
	return func(tx *TxBuilder) {
		if tx.TransactionBase == nil {
			tx.TransactionBase = &TransactionBase{
				MapInputsCell: make(map[string]*types.CellWithStatus),
			}
		}
		tx.TransactionBase.Transaction = txBase
	}
}

func WithCache(cache *cache.Cache) Options {
	return func(builder *TxBuilder) {
		builder.cache = cache
	}
}

func NewTxBuilder(opts ...Options) *TxBuilder {
	b := &TxBuilder{
		mapCellDep: make(map[string]*types.CellDep),
	}
	for _, v := range opts {
		v(b)
	}
	if b.TransactionBase == nil {
		b.TransactionBase = &TransactionBase{}
	}
	if b.TransactionBase.MapInputsCell == nil {
		b.TransactionBase.MapInputsCell = make(map[string]*types.CellWithStatus)
	}
	return b
}

func (d *TxBuilder) BuildTransactionWithCheckInputs(p *BuildTransactionParams, notCheckInputs bool) error {
	d.notCheckInputs = notCheckInputs
	err := d.newTx()
	if err != nil {
		return fmt.Errorf("newBaseTx err: %s", err.Error())
	}

	err = d.addInputsForTx(p.Inputs)
	if err != nil {
		return fmt.Errorf("addInputsForBaseTx err: %s", err.Error())
	}

	err = d.addOutputsForTx(p.Outputs, p.OutputsData)
	if err != nil {
		return fmt.Errorf("addOutputsForBaseTx err: %s", err.Error())
	}

	d.Transaction.Witnesses = append(d.Transaction.Witnesses, p.Witnesses...)
	d.otherWitnesses = append(d.otherWitnesses, p.OtherWitnesses...)

	if err := d.addMapCellDepWitnessForBaseTx(p.CellDeps); err != nil {
		return fmt.Errorf("addMapCellDepWitnessForBaseTx err: %s", err.Error())
	}

	for _, v := range p.HeadDeps {
		d.Transaction.HeaderDeps = append(d.Transaction.HeaderDeps, v)
	}

	return nil
}

func (d *TxBuilder) BuildTransaction(p *BuildTransactionParams) error {
	if err := d.BuildTransactionWithCheckInputs(p, false); err != nil {
		return err
	}
	if _, err := d.generateDigestListFromTx(nil); err != nil {
		return err
	}
	return nil
}

func (d *TxBuilder) TxString() string {
	txStr, _ := rpc.TransactionString(d.Transaction)
	return txStr
}

func (d *TxBuilder) GetTxBuilderTransactionString() string {
	bys, err := json.Marshal(d.Transaction)
	if err != nil {
		return ""
	}
	return string(bys)
}
