package errors

import (
	"errors"

	"github.com/nnikolash/wasp-types-exported/packages/isc"
	"github.com/nnikolash/wasp-types-exported/packages/kv"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/kv/collections"
	"github.com/nnikolash/wasp-types-exported/packages/vm/core/errors/coreerrors"
)

func ErrorTemplateKey(contractID isc.Hname) string {
	return PrefixErrorTemplateMap + string(contractID.Bytes())
}

// StateErrorCollectionWriter implements ErrorCollection. Is used for contract internal errors.
// It requires a reference to a KVStore such as the vmctx and the hname of the caller.
type StateErrorCollectionWriter struct {
	partition kv.KVStore
	hname     isc.Hname
}

func NewStateErrorCollectionWriter(partition kv.KVStore, hname isc.Hname) coreerrors.ErrorCollection {
	errorCollection := StateErrorCollectionWriter{
		partition: partition,
		hname:     hname,
	}

	return &errorCollection
}

func (e *StateErrorCollectionWriter) getErrorTemplateMap() *collections.Map {
	return collections.NewMap(e.partition, ErrorTemplateKey(e.hname))
}

func (e *StateErrorCollectionWriter) Get(errorID uint16) (*isc.VMErrorTemplate, bool) {
	errorMap := e.getErrorTemplateMap()
	errorIDKey := codec.EncodeUint16(errorID)

	errorBytes := errorMap.GetAt(errorIDKey)
	if errorBytes == nil {
		return nil, false
	}

	template, err := isc.VMErrorTemplateFromBytes(errorBytes)
	if err != nil {
		panic(err)
	}

	return template, true
}

func (e *StateErrorCollectionWriter) Register(messageFormat string) (*isc.VMErrorTemplate, error) {
	errorMap := e.getErrorTemplateMap()
	errorID := isc.GetErrorIDFromMessageFormat(messageFormat)

	if len(messageFormat) > isc.VMErrorMessageLimit {
		return nil, coreerrors.ErrErrorMessageTooLong
	}

	if t, ok := e.Get(errorID); ok && messageFormat != t.MessageFormat() {
		return nil, coreerrors.ErrErrorTemplateConflict.Create(errorID)
	}

	newError := isc.NewVMErrorTemplate(isc.NewVMErrorCode(e.hname, errorID), messageFormat)

	errorMap.SetAt(codec.EncodeUint16(errorID), newError.Bytes())

	return newError, nil
}

// StateErrorCollectionReader implements ErrorCollection partially. Is used for contract internal error readings only.
// It requires a reference to a KVStoreReader such as the vmctx and the hname of the caller.
type StateErrorCollectionReader struct {
	partition kv.KVStoreReader
	hname     isc.Hname
}

func (e *StateErrorCollectionReader) getErrorTemplateMap() *collections.ImmutableMap {
	return collections.NewMapReadOnly(e.partition, ErrorTemplateKey(e.hname))
}

func NewStateErrorCollectionReader(partition kv.KVStoreReader, hname isc.Hname) coreerrors.ErrorCollection {
	errorCollection := StateErrorCollectionReader{
		partition: partition,
		hname:     hname,
	}

	return &errorCollection
}

func (e *StateErrorCollectionReader) Get(errorID uint16) (*isc.VMErrorTemplate, bool) {
	errorMap := e.getErrorTemplateMap()
	errorIDKey := codec.EncodeUint16(errorID)

	errorBytes := errorMap.GetAt(errorIDKey)
	if errorBytes == nil {
		return nil, false
	}

	template, err := isc.VMErrorTemplateFromBytes(errorBytes)
	if err != nil {
		panic(err)
	}

	return template, true
}

func (e *StateErrorCollectionReader) Register(messageFormat string) (*isc.VMErrorTemplate, error) {
	return nil, errors.New("registering in read only maps is unsupported")
}
