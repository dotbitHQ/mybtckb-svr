package parser

import (
	"context"
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	log "github.com/sirupsen/logrus"
	"mybtckb-svr/dao"
	"mybtckb-svr/lib/contract"
	"mybtckb-svr/tables"
	"mybtckb-svr/txtool"
	"sync"
	"sync/atomic"
	"time"
)

var Handlers = make(map[string]func(*types.Transaction) (interface{}, error))

type BlockParser struct {
	ctx                context.Context
	cancel             context.CancelFunc
	wg                 *sync.WaitGroup
	currentBlockNumber uint64
	concurrencyNum     uint64
	confirmNum         uint64

	parserType           tables.ParserType
	contracts            *contract.Contracts
	dbDao                *dao.DbDao
	txTool               *txtool.TxTool
	mapTransactionHandle map[string]FuncTransactionHandle
}

type BlockParserOpts func(parser *BlockParser)

func WithContext(ctx context.Context, cancel context.CancelFunc) BlockParserOpts {
	return func(parser *BlockParser) {
		parser.ctx = ctx
		parser.cancel = cancel
	}
}

func WithWg(wg *sync.WaitGroup) BlockParserOpts {
	return func(parser *BlockParser) {
		parser.wg = wg
	}
}

func WithCurrentBlockNumber(currentBlockNumber uint64) BlockParserOpts {
	return func(parser *BlockParser) {
		parser.currentBlockNumber = currentBlockNumber
	}
}

func WithConcurrentNum(concurrencyNum uint64) BlockParserOpts {
	return func(parser *BlockParser) {
		parser.concurrencyNum = concurrencyNum
	}
}

func WithConfirmNum(confirmNum uint64) BlockParserOpts {
	return func(parser *BlockParser) {
		parser.confirmNum = confirmNum

	}
}

func WithParserType(parserType tables.ParserType) BlockParserOpts {
	return func(parser *BlockParser) {
		parser.parserType = parserType
	}
}

func WithContracts(contracts *contract.Contracts) BlockParserOpts {
	return func(parser *BlockParser) {
		parser.contracts = contracts
	}
}

func WithDbDao(dbDao *dao.DbDao) BlockParserOpts {
	return func(parser *BlockParser) {
		parser.dbDao = dbDao
	}
}

func WithTxTool(txTool *txtool.TxTool) BlockParserOpts {
	return func(parser *BlockParser) {
		parser.txTool = txTool
	}
}

func NewBlockParser(opts ...BlockParserOpts) *BlockParser {
	parser := &BlockParser{}
	for _, opt := range opts {
		opt(parser)
	}
	return parser
}

func (b *BlockParser) Run() error {
	b.registerTransactionHandle()
	currentBlockNumber, err := b.contracts.Client().GetTipBlockNumber(b.ctx)
	if err != nil {
		return fmt.Errorf("GetTipBlockNumber err: %s", err.Error())
	}
	if err := b.initCurrentBlockNumber(currentBlockNumber); err != nil {
		return fmt.Errorf("initCurrentBlockNumber err: %s", err.Error())
	}
	atomic.AddUint64(&b.currentBlockNumber, 1)
	b.wg.Add(1)
	go func() {
		for {
			select {
			default:
				time.Sleep(time.Millisecond * 300)
				latestBlockNumber, err := b.contracts.Client().GetTipBlockNumber(b.ctx)
				if err != nil {
					log.Error("GetTipBlockNumber err:", err.Error())
					continue
				}
				if b.concurrencyNum > 1 && b.currentBlockNumber < (latestBlockNumber-b.confirmNum-b.concurrencyNum) {
					nowTime := time.Now()
					if err = b.parserConcurrencyMode(); err != nil {
						log.Error("parserConcurrencyMode err:", err.Error(), b.currentBlockNumber)
					}
					log.Debug("parserConcurrencyMode time:", time.Since(nowTime).Seconds())
				} else if b.currentBlockNumber < (latestBlockNumber - b.confirmNum) { // check rollback
					nowTime := time.Now()
					if err = b.parserSubMode(); err != nil {
						log.Error("parserSubMode err:", err.Error(), b.currentBlockNumber)
					}
					log.Debug("parserSubMode time:", time.Since(nowTime).Seconds())
				} else {
					log.Debug("RunParser:", b.currentBlockNumber, latestBlockNumber)
					time.Sleep(time.Second * 10)
				}
			case <-b.ctx.Done():
				b.wg.Done()
				return
			}
		}
	}()
	return nil
}

func (b *BlockParser) initCurrentBlockNumber(currentBlockNumber uint64) error {
	if block, err := b.dbDao.FindBlockInfo(b.parserType); err != nil {
		return err
	} else if block.Id > 0 {
		b.currentBlockNumber = block.BlockNumber
	} else if b.currentBlockNumber == 0 && currentBlockNumber > 0 {
		b.currentBlockNumber = currentBlockNumber
	}
	return nil
}

func (b *BlockParser) parserSubMode() error {
	log.Debug("parserSubMode:", b.currentBlockNumber)
	block, err := b.contracts.Client().GetBlockByNumber(b.ctx, b.currentBlockNumber)
	if err != nil {
		return fmt.Errorf("GetBlockByNumber err: %s", err.Error())
	} else {
		blockHash := block.Header.Hash.Hex()
		parentHash := block.Header.ParentHash.Hex()
		log.Debug("parserSubMode:", b.currentBlockNumber, blockHash, parentHash)
		fork, err := b.checkFork(parentHash)
		if err != nil {
			return fmt.Errorf("checkFork err: %s", err.Error())
		}
		if fork {
			log.Debug("CheckFork is true:", b.currentBlockNumber, blockHash, parentHash)
			atomic.AddUint64(&b.currentBlockNumber, ^uint64(0))
		} else {
			if err := b.parsingBlockData(block); err != nil {
				return fmt.Errorf("parsingBlockData err: %s", err.Error())
			}
			if err = b.dbDao.CreateBlockInfo(b.parserType, b.currentBlockNumber, blockHash, parentHash); err != nil {
				return fmt.Errorf("CreateBlockInfo err: %s", err.Error())
			} else {
				atomic.AddUint64(&b.currentBlockNumber, 1)
			}
			if err = b.dbDao.DeleteBlockInfo(b.parserType, b.currentBlockNumber-20); err != nil {
				return fmt.Errorf("DeleteBlockInfo err: %s", err.Error())
			}
		}
	}
	return nil
}

func (b *BlockParser) checkFork(parentHash string) (bool, error) {
	block, err := b.dbDao.FindBlockInfoByBlockNumber(b.parserType, b.currentBlockNumber-1)
	if err != nil {
		return false, err
	}
	if block.Id == 0 {
		return false, nil
	}
	if block.BlockHash != parentHash {
		log.Warn("CheckFork:", b.currentBlockNumber, parentHash, block.BlockHash)
		return true, nil
	}
	return false, nil
}

func (b *BlockParser) parserConcurrencyMode() error {
	log.Debug("parserConcurrencyMode:", b.currentBlockNumber, b.concurrencyNum)
	for i := uint64(0); i < b.concurrencyNum; i++ {
		block, err := b.contracts.Client().GetBlockByNumber(b.ctx, b.currentBlockNumber)
		if err != nil {
			return fmt.Errorf("GetBlockByNumber err: %s [%d]", err.Error(), b.currentBlockNumber)
		}
		blockHash := block.Header.Hash.Hex()
		parentHash := block.Header.ParentHash.Hex()
		log.Debug("parserConcurrencyMode:", b.currentBlockNumber, blockHash, parentHash)

		if err := b.parsingBlockData(block); err != nil {
			return fmt.Errorf("parsingBlockData err: %s", err.Error())
		}
		if err := b.dbDao.CreateBlockInfo(b.parserType, b.currentBlockNumber, blockHash, parentHash); err != nil {
			return fmt.Errorf("CreateBlockInfo err: %s", err.Error())
		}
		atomic.AddUint64(&b.currentBlockNumber, 1)
	}
	if err := b.dbDao.DeleteBlockInfo(b.parserType, b.currentBlockNumber-20); err != nil {
		return fmt.Errorf("DeleteBlockInfo err: %s", err.Error())
	}
	return nil
}

func (b *BlockParser) parsingBlockData(block *types.Block) error {
	for idx, tx := range block.Transactions {
		txHash := tx.Hash.Hex()
		blockNumber := block.Header.Number
		blockTimestamp := block.Header.Timestamp
		//fmt.Println("block number: ", blockNumber)
		action, txTypeData, err := b.contracts.TxToDidCellAction(tx)
		if err != nil {
			return fmt.Errorf("TxToDidCellAction err : %s", err.Error())
		}
		if action != "" {
			//fmt.Println("===========: ", action, tx.Hash.Hex())
			//continue
			if handle, ok := b.mapTransactionHandle[action]; ok {
				if err := handle(FuncTransactionHandleReq{
					DbDao:          b.dbDao,
					Tx:             tx,
					TxTypeData:     &txTypeData,
					TxHash:         txHash,
					BlockNumber:    blockNumber,
					BlockIndex:     idx,
					BlockTimestamp: blockTimestamp,
					Action:         action,
				}); err != nil {
					log.WithFields(log.Fields{
						"action":       action,
						"block_number": blockNumber,
						"tx_hash":      tx.Hash.String(),
					}).WithError(err).Error("parse block error")
					return err
				}
			}
		}

	}
	return nil
}
