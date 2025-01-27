package sm_utils

import (
	"github.com/nnikolash/wasp-types-exported/packages/gpa"
)

type NodeRandomiser interface {
	UpdateNodeIDs([]gpa.NodeID)
	IsInitted() bool
	GetRandomOtherNodeIDs(int) []gpa.NodeID
}
