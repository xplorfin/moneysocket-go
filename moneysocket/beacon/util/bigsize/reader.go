package bigsize

import (
	"bytes"

	"github.com/lightningnetwork/lnd/tlv"
)

// Pop fetches an int res from a tlv byteStream
// returns error if this cannot be done
func Pop(byteStream []byte) (res uint64, byteRes []byte, err error) {
	var arr [8]byte
	copy(arr[:], byteStream[0:8])
	res, err = tlv.ReadVarInt(bytes.NewReader(byteStream), &arr)
	return res, byteStream[tlv.VarIntSize(res):], err
}

// GetTLVType will return the tlv.type of a given byteStream
func GetTLVType(byteStream []byte) (res tlv.Type, byteRes []byte, err error) {
	rawRes, byteStream, err := Pop(byteStream)
	return tlv.Type(rawRes), byteStream, err
}
