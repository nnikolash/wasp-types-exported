package coreerrors

import (
	"github.com/nnikolash/wasp-types-exported/packages/isc"
)

type ErrorCollection interface {
	Get(errorID uint16) (*isc.VMErrorTemplate, bool)
	Register(messageFormat string) (*isc.VMErrorTemplate, error)
}
