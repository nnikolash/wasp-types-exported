package test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnikolash/wasp-types-exported/packages/trie"
)

func TestProofScenariosBlake2b(t *testing.T) {
	runScenario := func(name string, scenario []string) {
		t.Run(name, func(t *testing.T) {
			store := NewInMemoryKVStore()
			initRoot := trie.MustInitRoot(store)
			tr, err := trie.NewTrieUpdatable(store, initRoot)
			require.NoError(t, err)

			checklist, root := runUpdateScenario(tr, store, scenario)
			trr, err := trie.NewTrieReader(store, root)
			require.NoError(t, err)
			for k, v := range checklist {
				vBin := trr.Get([]byte(k))
				if v == "" {
					require.EqualValues(t, 0, len(vBin))
				} else {
					require.EqualValues(t, []byte(v), vBin)
				}
				p := trr.MerkleProof([]byte(k))
				err = p.Validate(root.Bytes())
				require.NoError(t, err)
				if v != "" {
					cID := trie.CommitToData([]byte(v))
					err = p.ValidateWithTerminal(root.Bytes(), cID.Bytes())
					require.NoError(t, err)
				} else {
					require.True(t, p.IsProofOfAbsence())
				}
			}
		})
	}
	runScenario("1", []string{"a"})
	runScenario("2", []string{"a", "ab"})
	runScenario("3", []string{"a", "ab", "a/"})
	runScenario("4", []string{"a", "ab", "a/", "ab/"})
	runScenario("5", []string{"a", "ab", "abc", "a/", "ab/"})
	runScenario("rnd", genRnd3())

	longData := make([]string, 0)
	for _, k := range []string{"a", "ab", "abc", "bca"} {
		longData = append(longData, k+"/"+strings.Repeat(k, 200))
	}
	runScenario("long", longData)
}
