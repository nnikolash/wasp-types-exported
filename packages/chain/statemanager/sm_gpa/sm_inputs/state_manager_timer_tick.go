package sm_inputs

import (
	"time"

	"github.com/nnikolash/wasp-types-exported/packages/gpa"
)

type StateManagerTimerTick struct {
	time time.Time
}

var _ gpa.Input = &StateManagerTimerTick{}

func NewStateManagerTimerTick(timee time.Time) *StateManagerTimerTick {
	return &StateManagerTimerTick{time: timee}
}

func (smttT *StateManagerTimerTick) GetTime() time.Time {
	return smttT.time
}
