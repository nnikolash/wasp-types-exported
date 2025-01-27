package isc

import (
	"math"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/nnikolash/wasp-types-exported/packages/util/rwutil"
)

const addressIsNil rwutil.Kind = 0x80

func AddressFromReader(rr *rwutil.Reader) (address iotago.Address) {
	kind := rr.ReadKind()
	if kind == addressIsNil {
		return nil
	}
	addrSize := 0
	if rr.Err == nil {
		address, rr.Err = iotago.AddressSelector(uint32(kind))
		if rr.Err != nil {
			addrSize = 0
		} else {
			addrSize = address.Size()
		}
	}
	rr.PushBack().WriteKind(kind)
	rr.ReadSerialized(address, math.MaxUint16, addrSize)
	return address
}

func AddressToWriter(ww *rwutil.Writer, address iotago.Address) {
	if address == nil {
		ww.WriteKind(addressIsNil)
		return
	}
	ww.WriteSerialized(address, math.MaxUint16, address.Size())
}

func AddressFromBytes(data []byte) (iotago.Address, error) {
	rr := rwutil.NewBytesReader(data)
	return AddressFromReader(rr), rr.Err
}

func AddressToBytes(address iotago.Address) []byte {
	ww := rwutil.NewBytesWriter()
	AddressToWriter(ww, address)
	return ww.Bytes()
}
