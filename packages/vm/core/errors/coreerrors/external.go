package coreerrors

import (
	"github.com/nnikolash/wasp-types-exported/packages/isc"
)

var coreErrorCollection = NewCoreErrorCollection()

func Register(messageFormat string) *isc.VMErrorTemplate {
	template, err := coreErrorCollection.Register(messageFormat)
	if err != nil {
		panic(err)
	}
	return template
}

func All() ErrorCollection {
	return coreErrorCollection
}
