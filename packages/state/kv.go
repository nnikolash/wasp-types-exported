// Copyright 2022 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/trie"
)

// TrieKVAdapter is a KVStoreReader backed by a TrieReader
type TrieKVAdapter struct {
	*trie.TrieReader
}

var _ kv.KVStoreReader = &TrieKVAdapter{}

func (t *TrieKVAdapter) Get(key kv.Key) []byte {
	return t.TrieReader.Get([]byte(key))
}

func (t *TrieKVAdapter) Has(key kv.Key) bool {
	return t.TrieReader.Has([]byte(key))
}

func (t *TrieKVAdapter) Iterate(prefix kv.Key, f func(kv.Key, []byte) bool) {
	t.TrieReader.Iterator([]byte(prefix)).Iterate(func(k []byte, v []byte) bool {
		return f(kv.Key(k), v)
	})
}

func (t *TrieKVAdapter) IterateKeys(prefix kv.Key, f func(kv.Key) bool) {
	t.TrieReader.Iterator([]byte(prefix)).IterateKeys(func(k []byte) bool {
		return f(kv.Key(k))
	})
}

func (t *TrieKVAdapter) IterateKeysSorted(prefix kv.Key, f func(key kv.Key) bool) {
	t.IterateKeys(prefix, f)
}

func (t *TrieKVAdapter) IterateSorted(prefix kv.Key, f func(key kv.Key, value []byte) bool) {
	t.Iterate(prefix, f)
}
