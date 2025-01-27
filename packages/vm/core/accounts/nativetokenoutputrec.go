package accounts

import (
	"fmt"
	"io"
	"math/big"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

type NativeTokenOutputRec struct {
	OutputID          iotago.OutputID
	Amount            *big.Int
	StorageBaseTokens uint64 // always storage deposit
}

func NativeTokenOutputRecFromBytes(data []byte) (*NativeTokenOutputRec, error) {
	return rwutil.ReadFromBytes(data, new(NativeTokenOutputRec))
}

func MustNativeTokenOutputRecFromBytes(data []byte) *NativeTokenOutputRec {
	ret, err := NativeTokenOutputRecFromBytes(data)
	if err != nil {
		panic(err)
	}
	return ret
}

func (rec *NativeTokenOutputRec) Bytes() []byte {
	return rwutil.WriteToBytes(rec)
}

func (rec *NativeTokenOutputRec) String() string {
	return fmt.Sprintf("Native Token Account: base tokens: %d, amount: %d, outID: %s",
		rec.StorageBaseTokens, rec.Amount, rec.OutputID)
}

func (rec *NativeTokenOutputRec) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	rr.ReadN(rec.OutputID[:])
	rec.Amount = rr.ReadUint256()
	rec.StorageBaseTokens = rr.ReadAmount64()
	return rr.Err
}

func (rec *NativeTokenOutputRec) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteN(rec.OutputID[:])
	ww.WriteUint256(rec.Amount)
	ww.WriteAmount64(rec.StorageBaseTokens)
	return ww.Err
}
