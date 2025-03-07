package main

import (
	"context"
	"fmt"
	"time"

	"github.com/iotaledger/hive.go/kvstore"
	"github.com/nnikolash/wasp-types-exported/packages/chaindb"
	"github.com/nnikolash/wasp-types-exported/packages/trie"
)

func trieDiff(ctx context.Context, kvs kvstore.KVStore) {
	if blockIndex > blockIndex2 {
		blockIndex, blockIndex2 = blockIndex2, blockIndex
	}
	state1 := getState(kvs, blockIndex)
	state2 := getState(kvs, blockIndex2)

	start := time.Now()

	onlyOn1, onlyOn2 := trie.Diff(trie.NewHiveKVStoreAdapter(kvs, []byte{chaindb.PrefixTrie}), state1.TrieRoot(), state2.TrieRoot())

	fmt.Printf("Diff between blocks #%d -> #%d\n", blockIndex, blockIndex2)
	fmt.Printf("only on #%d: %d\n", blockIndex, len(onlyOn1))
	fmt.Printf("only on #%d: %d\n", blockIndex2, len(onlyOn2))
	elapsed := time.Since(start)
	fmt.Printf("Elapsed: %s\n", elapsed)
}
