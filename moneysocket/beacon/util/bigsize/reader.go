package bigsize

import (
	"bytes"

	"github.com/lightningnetwork/lnd/tlv"
)

func Pop(byteStream []byte) (res uint64, byteRes []byte, err error) {
	var arr [8]byte
	copy(arr[:], byteStream[0:8])
	res, err = tlv.ReadVarInt(bytes.NewReader(byteStream), &arr)
	return res, byteStream[tlv.VarIntSize(res):], err
}

// get tlv from a util
func GetTlv(byteStream []byte) (res tlv.Type, byteRes []byte, err error) {
	rawRes, byteStream, err := Pop(byteStream)
	return tlv.Type(rawRes), byteStream, err
}
