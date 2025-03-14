package blssig

import (
	"github.com/nnikolash/wasp-types-exported/packages/gpa"
)

const (
	msgTypeSigShare gpa.MessageType = iota
)

func (cc *ccImpl) UnmarshalMessage(data []byte) (gpa.Message, error) {
	return gpa.UnmarshalMessage(data, gpa.Mapper{
		msgTypeSigShare: func() gpa.Message { return new(msgSigShare) },
	})
}
