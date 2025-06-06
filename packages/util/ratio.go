package util

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/iotaledger/hive.go/serializer/v2"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

// A + B
const RatioByteSize = serializer.UInt32ByteSize + serializer.UInt32ByteSize

// Ratio32 represents a ratio (a:b) between two quantities, expressed as two uint32 values.
type Ratio32 struct {
	A uint32 `json:"a" swagger:"min(0),required"`
	B uint32 `json:"b" swagger:"min(0),required"`
}

func Ratio32FromBytes(data []byte) (ret Ratio32, err error) {
	_, err = rwutil.ReadFromBytes(data, &ret)
	return ret, err
}

func Ratio32FromString(s string) (ret Ratio32, err error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return ret, errors.New("invalid Ratio32 string")
	}
	a, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return ret, err
	}
	b, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return ret, err
	}
	ret.A = uint32(a)
	ret.B = uint32(b)
	return ret, nil
}

func (ratio Ratio32) Bytes() []byte {
	return rwutil.WriteToBytes(&ratio)
}

func (ratio Ratio32) String() string {
	return fmt.Sprintf("%d:%d", ratio.A, ratio.B)
}

func ceil(x, dividend, divisor uint64) uint64 {
	return (x*dividend + (divisor - 1)) / divisor
}

// YFloor64 computes y = floor(x * b / a)
func (ratio Ratio32) YFloor64(x uint64) uint64 {
	return x * uint64(ratio.B) / uint64(ratio.A)
}

// YCeil64 computes y = ceil(x * b / a)
func (ratio Ratio32) YCeil64(x uint64) uint64 {
	return ceil(x, uint64(ratio.B), uint64(ratio.A))
}

// XFloor64 computes x = floor(y * a / b)
func (ratio Ratio32) XFloor64(y uint64) uint64 {
	return y * uint64(ratio.A) / uint64(ratio.B)
}

// XCeil64 computes x = ceil(y * a / b)
func (ratio Ratio32) XCeil64(y uint64) uint64 {
	return ceil(y, uint64(ratio.A), uint64(ratio.B))
}

// Set is part of the pflag.Value interface. It accepts a string in the form "a:b".
func (ratio *Ratio32) Set(s string) error {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return errors.New("invalid format for Ratio32")
	}
	a, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return err
	}
	b, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return err
	}
	ratio.A = uint32(a)
	ratio.B = uint32(b)
	return nil
}

// Type is part of the pflag.Value interface.
func (ratio Ratio32) Type() string {
	return "Ratio32"
}

func (ratio Ratio32) HasZeroComponent() bool {
	return ratio.A == 0 || ratio.B == 0
}

func (ratio Ratio32) IsValid() bool {
	return ratio.IsEmpty() || !ratio.HasZeroComponent()
}

func (ratio Ratio32) IsEmpty() bool {
	ZeroGasFee := Ratio32{}
	return ratio == ZeroGasFee
}

func (ratio *Ratio32) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	ratio.A = rr.ReadUint32()
	ratio.B = rr.ReadUint32()
	if rr.Err == nil && !ratio.IsValid() {
		rr.Err = errors.New("ratio has zero component")
	}
	return rr.Err
}

func (ratio *Ratio32) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteUint32(ratio.A)
	ww.WriteUint32(ratio.B)
	return ww.Err
}
