package util

import (
	"bytes"

	"github.com/lightningnetwork/lnd/tlv"
)

func TlvRecordToBytes(record ...tlv.Record) []byte {
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
