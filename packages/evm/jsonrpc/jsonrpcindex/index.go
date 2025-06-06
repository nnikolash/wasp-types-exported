package jsonrpcindex

import (
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	hivedb "github.com/iotaledger/hive.go/db"
	"github.com/iotaledger/hive.go/kvstore"
	"github.com/nnikolash/wasp-types-exported/packages/database"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/state"
	"github.com/nnikolash/wasp-types-exported/packages/trie"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/blocklog"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/evm/emulator"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/governance"
)

type Index struct {
	store           kvstore.KVStore
	blockchainDB    func(chainState state.State) *emulator.BlockchainDB
	stateByTrieRoot func(trieRoot trie.Hash) (state.State, error)

	mu sync.Mutex
}

func New(
	blockchainDB func(chainState state.State) *emulator.BlockchainDB,
	stateByTrieRoot func(trieRoot trie.Hash) (state.State, error),
	indexDbEngine hivedb.Engine,
	indexDbPath string,
) *Index {
	db, err := database.NewDatabase(indexDbEngine, indexDbPath, true, false, database.CacheSizeDefault)
	if err != nil {
		panic(err)
	}
	return &Index{
		store:           db.KVStore(),
		blockchainDB:    blockchainDB,
		stateByTrieRoot: stateByTrieRoot,
		mu:              sync.Mutex{},
	}
}

func (c *Index) IndexBlock(trieRoot trie.Hash) {
	c.mu.Lock()
	defer c.mu.Unlock()
	state, err := c.stateByTrieRoot(trieRoot)
	if err != nil {
		panic(err)
	}
	blockKeepAmount := governance.NewStateAccess(state).GetBlockKeepAmount()
	if blockKeepAmount == -1 {
		return // pruning disabled, never cache anything
	}
	// cache the block that will be pruned next (this way reorgs are okay, as long as it never reorgs more than `blockKeepAmount`, which would be catastrophic)
	if state.BlockIndex() < uint32(blockKeepAmount-1) {
		return
	}
	blockIndexToCache := state.BlockIndex() - uint32(blockKeepAmount-1)
	cacheUntil := uint32(0)
	lastBlockIndexed := c.lastBlockIndexed()
	if lastBlockIndexed != nil {
		cacheUntil = *lastBlockIndexed
	}

	// we need to look at the next block to get the trie commitment of the block we want to cache
	nextBlockInfo, found := blocklog.NewStateAccess(state).BlockInfo(blockIndexToCache + 1)
	if !found {
		panic(fmt.Errorf("block %d not found on active state %d", blockIndexToCache, state.BlockIndex()))
	}

	// start in the active state of the block to cache
	activeStateToCache, err := c.stateByTrieRoot(nextBlockInfo.PreviousL1Commitment().TrieRoot())
	if err != nil {
		panic(err)
	}

	for i := blockIndexToCache; i >= cacheUntil; i-- {
		// walk back and save all blocks between [lastBlockIndexCached...blockIndexToCache]

		blockinfo, found := blocklog.NewStateAccess(activeStateToCache).BlockInfo(i)
		if !found {
			panic(fmt.Errorf("block %d not found on active state %d", i, state.BlockIndex()))
		}

		db := c.blockchainDB(activeStateToCache)
		blockTrieRoot := activeStateToCache.TrieRoot()
		c.setBlockTrieRootByIndex(i, blockTrieRoot)

		evmBlock := db.GetCurrentBlock()
		c.setBlockIndexByHash(evmBlock.Hash(), i)

		blockTransactions := evmBlock.Transactions()
		for _, tx := range blockTransactions {
			c.setBlockIndexByTxHash(tx.Hash(), i)
		}
		// walk backwards until all blocks are cached
		if i == 0 {
			// nothing more to cache, don't try to walk back further
			break
		}
		activeStateToCache, err = c.stateByTrieRoot(blockinfo.PreviousL1Commitment().TrieRoot())
		if err != nil {
			panic(err)
		}
	}
	c.setLastBlockIndexed(blockIndexToCache)
	c.store.Flush()
}

func (c *Index) BlockByNumber(n *big.Int) *types.Block {
	if n == nil {
		return nil
	}
	db := c.evmDBFromBlockIndex(uint32(n.Uint64()))
	if db == nil {
		return nil
	}
	return db.GetBlockByNumber(n.Uint64())
}

func (c *Index) BlockByHash(hash common.Hash) *types.Block {
	blockIndex := c.blockIndexByHash(hash)
	if blockIndex == nil {
		return nil
	}
	return c.evmDBFromBlockIndex(*blockIndex).GetBlockByHash(hash)
}

func (c *Index) BlockTrieRootByIndex(n uint32) *trie.Hash {
	return c.blockTrieRootByIndex(n)
}

func (c *Index) TxByHash(hash common.Hash) (tx *types.Transaction, blockHash common.Hash, blockNumber, txIndex uint64) {
	blockIndex := c.blockIndexByTxHash(hash)
	if blockIndex == nil {
		return nil, common.Hash{}, 0, 0
	}
	tx, blockHash, blockNumber, txIndex, err := c.evmDBFromBlockIndex(*blockIndex).GetTransactionByHash(hash)
	if err != nil {
		panic(err)
	}
	return tx, blockHash, blockNumber, txIndex
}

func (c *Index) GetReceiptByTxHash(hash common.Hash) *types.Receipt {
	blockIndex := c.blockIndexByTxHash(hash)
	if blockIndex == nil {
		return nil
	}
	return c.evmDBFromBlockIndex(*blockIndex).GetReceiptByTxHash(hash)
}

func (c *Index) TxByBlockHashAndIndex(blockHash common.Hash, txIndex uint64) (tx *types.Transaction, blockNumber uint64) {
	blockIndex := c.blockIndexByHash(blockHash)
	if blockIndex == nil {
		return nil, 0
	}
	block := c.evmDBFromBlockIndex(*blockIndex).GetBlockByHash(blockHash)
	if block == nil {
		return nil, 0
	}
	txs := block.Transactions()
	if txIndex > uint64(len(txs)) {
		return nil, 0
	}
	return txs[txIndex], block.NumberU64()
}

func (c *Index) TxByBlockNumberAndIndex(blockNumber *big.Int, txIndex uint64) (tx *types.Transaction, blockHash common.Hash) {
	if blockNumber == nil {
		return nil, common.Hash{}
	}
	db := c.evmDBFromBlockIndex(uint32(blockNumber.Uint64()))
	if db == nil {
		return nil, common.Hash{}
	}
	block := db.GetBlockByHash(blockHash)
	if block == nil {
		return nil, common.Hash{}
	}
	txs := block.Transactions()
	if txIndex > uint64(len(txs)) {
		return nil, common.Hash{}
	}
	return txs[txIndex], block.Hash()
}

func (c *Index) TxsByBlockNumber(blockNumber *big.Int) types.Transactions {
	if blockNumber == nil {
		return nil
	}
	db := c.evmDBFromBlockIndex(uint32(blockNumber.Uint64()))
	if db == nil {
		return nil
	}
	block := db.GetBlockByNumber(blockNumber.Uint64())
	if block == nil {
		return nil
	}
	return block.Transactions()
}

// internals

const (
	PrefixLastBlockIndexed = iota
	PrefixBlockTrieRootByIndex
	PrefixBlockIndexByTxHash
	PrefixBlockIndexByHash
)

func KeyLastBlockIndexed() kvstore.Key {
	return []byte{PrefixLastBlockIndexed}
}

func KeyBlockTrieRootByIndex(i uint32) kvstore.Key {
	key := []byte{PrefixBlockTrieRootByIndex}
	key = append(key, codec.EncodeUint32(i)...)
	return key
}

func KeyBlockIndexByTxHash(hash common.Hash) kvstore.Key {
	key := []byte{PrefixBlockIndexByTxHash}
	key = append(key, hash[:]...)
	return key
}

func KeyBlockIndexByHash(hash common.Hash) kvstore.Key {
	key := []byte{PrefixBlockIndexByHash}
	key = append(key, hash[:]...)
	return key
}

func (c *Index) get(key kvstore.Key) []byte {
	ret, err := c.store.Get(key)
	if err != nil {
		if errors.Is(err, kvstore.ErrKeyNotFound) {
			return nil
		}
		panic(err)
	}
	return ret
}

func (c *Index) set(key kvstore.Key, value []byte) {
	err := c.store.Set(key, value)
	if err != nil {
		panic(err)
	}
}

func (c *Index) setLastBlockIndexed(n uint32) {
	c.set(KeyLastBlockIndexed(), codec.EncodeUint32(n))
}

func (c *Index) lastBlockIndexed() *uint32 {
	bytes := c.get(KeyLastBlockIndexed())
	if bytes == nil {
		return nil
	}
	ret := codec.MustDecodeUint32(bytes)
	return &ret
}

func (c *Index) setBlockTrieRootByIndex(i uint32, hash trie.Hash) {
	c.set(KeyBlockTrieRootByIndex(i), hash.Bytes())
}

func (c *Index) blockTrieRootByIndex(i uint32) *trie.Hash {
	bytes := c.get(KeyBlockTrieRootByIndex(i))
	if bytes == nil {
		return nil
	}
	hash, err := trie.HashFromBytes(bytes)
	if err != nil {
		panic(err)
	}
	return &hash
}

func (c *Index) setBlockIndexByTxHash(txHash common.Hash, blockIndex uint32) {
	c.set(KeyBlockIndexByTxHash(txHash), codec.EncodeUint32(blockIndex))
}

func (c *Index) blockIndexByTxHash(txHash common.Hash) *uint32 {
	bytes := c.get(KeyBlockIndexByTxHash(txHash))
	if bytes == nil {
		return nil
	}
	ret := codec.MustDecodeUint32(bytes)
	return &ret
}

func (c *Index) setBlockIndexByHash(hash common.Hash, blockIndex uint32) {
	c.set(KeyBlockIndexByHash(hash), codec.EncodeUint32(blockIndex))
}

func (c *Index) blockIndexByHash(hash common.Hash) *uint32 {
	bytes := c.get(KeyBlockIndexByHash(hash))
	if bytes == nil {
		return nil
	}
	ret := codec.MustDecodeUint32(bytes)
	return &ret
}

func (c *Index) evmDBFromBlockIndex(n uint32) *emulator.BlockchainDB {
	trieRoot := c.blockTrieRootByIndex(n)
	if trieRoot == nil {
		return nil
	}
	state, err := c.stateByTrieRoot(*trieRoot)
	if err != nil {
		panic(err)
	}
	return c.blockchainDB(state)
}
