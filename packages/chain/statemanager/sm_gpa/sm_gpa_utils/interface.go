package sm_gpa_utils

import (
	"time"

	"github.com/nnikolash/wasp-types-exported/packages/state"
)

type BlockCache interface {
	AddBlock(state.Block)
	GetBlock(*state.L1Commitment) state.Block
	CleanOlderThan(time.Time)
	Size() int
}

type BlockWAL interface {
	Write(state.Block) error
	Contains(state.BlockHash) bool
	Read(state.BlockHash) (state.Block, error)
	ReadAllByStateIndex(cb func(stateIndex uint32, block state.Block) bool) error
}

type TimeProvider interface {
	SetNow(time.Time)
	GetNow() time.Time
	After(time.Duration) <-chan time.Time
}
