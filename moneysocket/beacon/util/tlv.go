package util

import (
	"errors"
	"fmt"

	tlvHelper "github.com/lightningnetwork/lnd/tlv"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util/bigsize"
)

// see https://github.com/moneysocket/py-moneysocket/blob/main/moneysocket/utl/third_party/bolt/tlv.py

// For encoding/decoding values to/from Tlv (Type-Length-Value) byte strings
// as defined in:
// https://github.com/lightningnetwork/lightning-rfc/blob/master/01-messaging.md#type-length-value-format
// note: I tried to use the golang implementation, but didn't quite get the sreaming encodes/decodes to work.
// this can be fixed in a future version
type Tlv struct {
	t tlvHelper.Type
	l int
	v []byte
}

func NewTlv(t tlvHelper.Type, v []byte) (res Tlv) {
	res.t = t
	res.v = v
	res.l = len(res.v)
	return res
}

// get tlv.type
func (tlv Tlv) Type() tlvHelper.Type {
	return tlv.t
}

// get l of tlv.t
func (tlv Tlv) Length() int {
	return tlv.l
}

// get value of tlv
func (tlv Tlv) Value() []byte {
	return tlv.v
}

func TlvPop(byteString []byte) (tlv Tlv, remainder []byte, err error) {
	t, byteString, err := bigsize.GetTlv(byteString)
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
