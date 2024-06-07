package contract

import (
	"context"
	"errors"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"mybtckb-svr/lib/common"
	"mybtckb-svr/lib/outpoint_cache"
	"sync"
	"time"
)

var (
	ConfigCellNotFound = errors.New("config cell not found")
	CellNotFound       = errors.New("cell not found")
)

type ContractsOption func(*Contracts)

func WithClient(client rpc.Client) ContractsOption {
	return func(dc *Contracts) {
		dc.client = client
	}
}

func WithNetType(net common.NetType) ContractsOption {
	return func(dc *Contracts) {
		dc.net = net
	}
}

func WithWaitGroup(wg *sync.WaitGroup) ContractsOption {
	return func(dc *Contracts) {
		dc.wg = wg
	}
}

func WithCache(cache *outpoint_cache.Cache) ContractsOption {
	return func(dc *Contracts) {
		dc.cache = cache
	}
}

type Contracts struct {
	ctx              context.Context
	client           rpc.Client
	net              common.NetType
	wg               *sync.WaitGroup
	cache            *outpoint_cache.Cache
	XudtType         *Info
	UniqueType       *Info
	SporeType        *Info
	SporeClusterType *Info
}

func NewContracts(ctx context.Context, opts ...ContractsOption) *Contracts {
	contracts := &Contracts{
		ctx: ctx,
	}
	for _, opt := range opts {
		opt(contracts)
	}
	return contracts
}

func (c *Contracts) Init(t time.Duration) error {
	if err := c.initContracts(); err != nil {
		return err
	}
	return nil
}

func (c *Contracts) Client() rpc.Client {
	return c.client
}

func (c *Contracts) NetType() common.NetType {
	return c.net
}

func (c *Contracts) Mode() address.Mode {
	if c.net == common.NetTypeMainNet {
		return address.Mainnet
	}
	return address.Testnet
}

func (c *Contracts) Env() Env {
	if c.net == common.NetTypeMainNet {
		return EnvMain
	} else {
		return EnvTest
	}
}

func (c *Contracts) Cache() *outpoint_cache.Cache {
	return c.cache
}

func (c *Contracts) GetScriptByTypeIdAndLock(lock *types.Script, typeId types.Hash) (*indexer.LiveCell, error) {
	cells, err := c.Client().GetCells(context.Background(), &indexer.SearchKey{
		Script: &types.Script{
			CodeHash: typeId,
			HashType: types.HashTypeType,
		},
		ScriptType: indexer.ScriptTypeType,
		Filter: &indexer.CellsFilter{
			Script: lock,
		},
	}, indexer.SearchOrderDesc, 1, "")
	if err != nil {
		return nil, err
	}
	if len(cells.Objects) == 0 {
		return nil, CellNotFound
	}
	return cells.Objects[0], nil
}

func (c *Contracts) GetCellByLockAndTypeId(lock *types.Script, typeId types.Hash) (*indexer.LiveCell, error) {
	configs, err := c.Client().GetCells(context.Background(), &indexer.SearchKey{
		Script: &types.Script{
			CodeHash: typeId,
			HashType: types.HashTypeType,
		},
		ScriptType: indexer.ScriptTypeType,
		Filter: &indexer.CellsFilter{
			Script: lock,
		},
	}, indexer.SearchOrderDesc, 1, "")
	if err != nil {
		return nil, err
	}
	if len(configs.Objects) == 0 {
		return nil, CellNotFound
	}
	return configs.Objects[0], nil
}

func (c *Contracts) initContracts() error {
	env := c.Env()
	c.XudtType = &Info{
		OutPoint: &types.OutPoint{
			TxHash: types.HexToHash(env.Xudt.TxHash),
			Index:  env.Xudt.Index,
		},
		CodeHash: types.HexToHash(env.Xudt.CodeHash),
		HashType: env.Xudt.HashType,
	}
	c.UniqueType = &Info{
		OutPoint: &types.OutPoint{
			TxHash: types.HexToHash(env.UniqueCell.TxHash),
			Index:  env.UniqueCell.Index,
		},
		CodeHash: types.HexToHash(env.UniqueCell.CodeHash),
		HashType: env.UniqueCell.HashType,
	}
	c.SporeType = &Info{
		OutPoint: &types.OutPoint{
			TxHash: types.HexToHash(env.Spore.TxHash),
			Index:  env.Spore.Index,
		},
		CodeHash: types.HexToHash(env.Spore.CodeHash),
		HashType: env.Spore.HashType,
	}
	c.SporeClusterType = &Info{
		OutPoint: &types.OutPoint{
			TxHash: types.HexToHash(env.SporeCluster.TxHash),
			Index:  env.SporeCluster.Index,
		},
		CodeHash: types.HexToHash(env.SporeCluster.CodeHash),
		HashType: env.SporeCluster.HashType,
	}
	return nil
}

type Info struct {
	OutPoint *types.OutPoint // contract outpoint
	//OutPut   *types.CellOutput // contract script
	CodeHash types.Hash           // code hash
	HashType types.ScriptHashType // hash type
}

func (d *Info) ToCellDep() *types.CellDep {
	return &types.CellDep{
		OutPoint: d.OutPoint,
		DepType:  types.DepTypeCode,
	}
}

func (c *Contracts) GetXudtTypeScript(args []byte) *types.Script {
	xudtType := &types.Script{
		CodeHash: c.XudtType.CodeHash,
		HashType: c.XudtType.HashType,
		Args:     args,
	}
	return xudtType
}

func (c *Contracts) GetSporeTypeScript(args []byte) *types.Script {
	sporeType := &types.Script{
		CodeHash: c.SporeType.CodeHash,
		HashType: c.SporeType.HashType,
		Args:     args,
	}
	return sporeType
}

//func (d *Info) ToScript(args []byte) *types.Script {
//	return &types.Script{
//		CodeHash: d.ContractTypeId,
//		HashType: types.HashTypeType,
//		Args:     args,
//	}
//}
//
//func (d *Info) IsSameTypeId(codeHash types.Hash) bool {
//	return d.ContractTypeId.Hex() == codeHash.Hex()
//}
