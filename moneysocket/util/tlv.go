package util

import (
	"bytes"

	"github.com/lightningnetwork/lnd/tlv"
)

// TLVRecordToBytes encodes a TLV (type-length-value, defined in BOLT #1: https://git.io/JLCRq )
// into a byte slice.
func TLVRecordToBytes(record ...tlv.Record) []byte {
	stream, err := tlv.NewStream(record...)
	if err != nil {
		panic(err)
	}
	var w bytes.Buffer
	err = stream.Encode(&w)
	if err != nil {
		panic(err)
	}

	return w.Bytes()
}
