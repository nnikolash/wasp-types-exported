package trie

import (
	"io"

	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

func (tr *TrieReader) TakeSnapshot(w io.Writer) error {
	// Some duplicated nodes and values might be written more than once in the snapshot;
	// Using a size-capped map to prevent this.
	// If the cap is reached, the generated snapshot will contain duplicate information,
	// but will still be correct.
	seenNodes := make(map[Hash]struct{})
	seenValues := make(map[string]struct{})
	const mapSizeCap = 2_000_000 / HashSizeBytes // 2 MB max for each map

	ww := rwutil.NewWriter(w)
	tr.IterateNodes(func(_ []byte, n *NodeData, depth int) IterateNodesAction {
		if _, seen := seenNodes[n.Commitment]; seen {
			return IterateContinue
		}
		if len(seenNodes) < mapSizeCap {
			seenNodes[n.Commitment] = struct{}{}
		}

		ww.WriteBytes(n.Bytes())
		if n.Terminal != nil && !n.Terminal.IsValue {
			valueKey := n.Terminal.Bytes()
			if _, seen := seenValues[string(valueKey)]; !seen {
				ww.WriteBool(true)
				value := tr.nodeStore.valueStore.Get(valueKey)
				ww.WriteBytes(value)
				if len(seenValues) < mapSizeCap {
					seenValues[string(valueKey)] = struct{}{}
				}
			} else {
				ww.WriteBool(false)
			}
		}
		if ww.Err != nil {
			return IterateStop
		}
		return IterateContinue
	})
	return ww.Err
}

func RestoreSnapshot(r io.Reader, store KVStore) error {
	triePartition := makeWriterPartition(store, partitionTrieNodes)
	valuePartition := makeWriterPartition(store, partitionValues)
	refcounts := newRefcounts(store)
	rr := rwutil.NewReader(r)
	for rr.Err == nil {
		nodeBytes := rr.ReadBytes()
		if rr.Err == io.EOF {
			return nil
		}
		n, err := nodeDataFromBytes(nodeBytes)
		if err != nil {
			return err
		}
		n.updateCommitment()
		nodeKey := n.Commitment.Bytes()

		var valueKey, value []byte
		if n.Terminal != nil && !n.Terminal.IsValue {
			if rr.ReadBool() {
				value = rr.ReadBytes()
				if rr.Err != nil {
					break
				}
				valueKey = n.Terminal.Bytes()
			}
		}

		if refcounts.GetNode(n.Commitment) == 0 {
			// node is new -- save it and set node and value refcounts to 1
			triePartition.Set(nodeKey, nodeBytes)
			if valueKey != nil {
				valuePartition.Set(valueKey, value)
			}
			refcounts.incNodeAndValue(n)

			// Increment the refcounts of the children that already exist
			// (for the others, their refcount will be set to 1 in a
			// later iteration, when they are read from the snapshot).
			n.iterateChildren(func(i byte, commitment Hash) bool {
				if refcounts.GetNode(commitment) > 0 {
					refcounts.incNode(commitment)
				}
				return true
			})
		}

		if rr.Err != nil {
			break
		}
	}
	return rr.Err
}
