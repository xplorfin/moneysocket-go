package util

import (
	"io"

	"github.com/lightningnetwork/lnd/tlv"
)

// EVarInt is an Encoder for variable byte slices. An error is returned if val
// is not *[]byte.
func EVarInt(w io.Writer, val interface{}, b *[8]byte) error {
	if c, ok := val.(*uint64); ok {
		return tlv.WriteVarInt(w, *c, b)
	}
	return tlv.NewTypeForEncodingErr(val, "uint64")
}

// DVarBytes is a Decoder for variable byte slices. An error is returned if val
// is not *[]byte.
func DVarInt(r io.Reader, val interface{}, _ *[8]byte, l uint64) error {
	panic("method not yet implemented")
}
