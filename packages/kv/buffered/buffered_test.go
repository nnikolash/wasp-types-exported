package buffered

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/hive.go/kvstore/mapdb"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
)

func TestBufferedKVStore(t *testing.T) {
	db := mapdb.NewMapDB()
	_ = db.Set([]byte("abcd"), []byte("v1"))

	realm, err := db.WithRealm([]byte("ab"))
	require.NoError(t, err)

	v, err := realm.Get([]byte("cd"))
	require.NoError(t, err)
	require.Equal(t, []byte("v1"), v)

	b := NewBufferedKVStore(kv.NewHiveKVStoreReader(realm))

	v = b.Get("cd")
	require.Equal(t, []byte("v1"), v)

	require.EqualValues(
		t,
		map[kv.Key][]byte{
			"cd": []byte("v1"),
		},
		b.DangerouslyDumpToDict(),
	)

	n := 0
	b.Iterate(kv.EmptyPrefix, func(key kv.Key, value []byte) bool {
		require.Equal(t, kv.Key("cd"), key)
		require.Equal(t, []byte("v1"), value)
		n++
		return true
	})
	require.Equal(t, 1, n)

	n = 0
	b.Iterate("c", func(key kv.Key, value []byte) bool {
		require.Equal(t, kv.Key("cd"), key)
		require.Equal(t, []byte("v1"), value)
		n++
		return true
	})
	require.Equal(t, 1, n)

	n = 0
	b.IterateKeys(kv.EmptyPrefix, func(key kv.Key) bool {
		require.Equal(t, kv.Key("cd"), key)
		n++
		return true
	})
	require.Equal(t, 1, n)

	b.Set("cd", []byte("v2"))

	// not committed to DB
	v, err = realm.Get([]byte("cd"))
	require.NoError(t, err)
	require.Equal(t, []byte("v1"), v)

	v = b.Get("cd")
	require.Equal(t, []byte("v2"), v)

	require.EqualValues(
		t,
		map[kv.Key][]byte{
			"cd": []byte("v2"),
		},
		b.DangerouslyDumpToDict(),
	)
}

func TestIterateSorted(t *testing.T) {
	db := mapdb.NewMapDB()
	_ = db.Set([]byte("1246"), []byte("v1246"))
	_ = db.Set([]byte("1248"), []byte("v1248"))
	_ = db.Set([]byte("1345"), []byte("v1345"))
	_ = db.Set([]byte("1259"), []byte("v1259"))
	_ = db.Set([]byte("2345"), []byte("v2345"))
	_ = db.Set([]byte("1247"), []byte("v1247"))
	_ = db.Set([]byte("3123"), []byte("v3123"))
	_ = db.Set([]byte("1234"), []byte("v1234"))
	_ = db.Set([]byte("1323"), []byte("v1323"))
	_ = db.Set([]byte("1245"), []byte("v1245"))

	realm, err := db.WithRealm([]byte("1"))
	require.NoError(t, err)
	b := NewBufferedKVStore(kv.NewHiveKVStoreReader(realm))

	b.Del("246")
	b.Set("250", []byte("v1250"))

	var seen []kv.Key
	b.IterateKeysSorted("2", func(k kv.Key) bool {
		seen = append(seen, k)
		return true
	})
	require.Equal(t, []kv.Key{"234", "245", "247", "248", "250", "259"}, seen)
}
