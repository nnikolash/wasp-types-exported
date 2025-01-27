package accounts

import (
	"io"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/packages/kv/codec"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

// FoundryOutputRec contains information to reconstruct output
type FoundryOutputRec struct {
	OutputID    iotago.OutputID
	Amount      uint64 // always storage deposit
	TokenScheme iotago.TokenScheme
	Metadata    []byte
}

func (rec *FoundryOutputRec) Bytes() []byte {
	return rwutil.WriteToBytes(rec)
}

func FoundryOutputRecFromBytes(data []byte) (*FoundryOutputRec, error) {
	return rwutil.ReadFromBytes(data, new(FoundryOutputRec))
}

func MustFoundryOutputRecFromBytes(data []byte) *FoundryOutputRec {
	ret, err := FoundryOutputRecFromBytes(data)
	if err != nil {
		panic(err)
	}
	return ret
}

func (rec *FoundryOutputRec) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	rr.ReadN(rec.OutputID[:])
	rec.Amount = rr.ReadUint64()
	tokenScheme := rr.ReadBytes()
	if rr.Err == nil {
		rec.TokenScheme, rr.Err = codec.DecodeTokenScheme(tokenScheme)
	}
	rec.Metadata = rr.ReadBytes()
	return rr.Err
}

func (rec *FoundryOutputRec) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteN(rec.OutputID[:])
	ww.WriteUint64(rec.Amount)
	if ww.Err == nil {
		tokenScheme := codec.EncodeTokenScheme(rec.TokenScheme)
		ww.WriteBytes(tokenScheme)
	}
	ww.WriteBytes(rec.Metadata)
	return ww.Err
}
