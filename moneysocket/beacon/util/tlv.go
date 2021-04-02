package util

import (
	"errors"
	"fmt"

	tlvHelper "github.com/lightningnetwork/lnd/tlv"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util/bigsize"
)

// see https://github.com/moneysocket/py-moneysocket/blob/main/moneysocket/utl/third_party/bolt/tlv.py

// TLV is used for encoding/decoding values to/from TLV (Type-Length-Value) byte strings
// as defined in: https://git.io/JmwOl
// note: I tried to use the golang implementation, but didn't quite get the streaming encodes/decodes to work.
// this can be fixed in a future version.
type TLV struct {
	// t represents the type
	t tlvHelper.Type
	// l is the length of v
	l int
	// v is the value
	v []byte
}

// Type get tlv.type.
func (tlv TLV) Type() tlvHelper.Type {
	return tlv.t
}

// Length gets l of tlv.t.
func (tlv TLV) Length() int {
	return tlv.l
}

// Value gets value of tlv.
func (tlv TLV) Value() []byte {
	return tlv.v
}

// TLVPop returns the first TLV from a byteString and returns the remainder
// returns error if this cannot be done.
func TLVPop(byteString []byte) (tlv TLV, remainder []byte, err error) {
	t, byteString, err := bigsize.GetTLVType(byteString)
	if err != nil {
		return tlv, remainder, fmt.Errorf("could not get type %d", t)
	}
	l, byteString, err := bigsize.Pop(byteString)
	if err != nil {
		return tlv, remainder, fmt.Errorf("could not get length %d", l)
	}
	if uint64(len(byteString)) < l {
		return tlv, remainder, errors.New("value truncated")
	}
	tlv.t = t
	tlv.l = int(l)
	tlv.v = byteString[:l]
	return tlv, byteString[l:], err
}
