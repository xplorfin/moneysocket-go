package location

import (
	"github.com/lightningnetwork/lnd/tlv"
	"github.com/xplorfin/moneysocket-go/moneysocket/util"
)

// PlaceholderTLVType is the default tlv type for unimplemented location beacons
const PlaceholderTLVType = 0

// PlaceHolderTLV encodes a placeholder tlv for python-parity tests
func PlaceHolderTLV(placeholder string) tlv.Record {
	byteHolder := []byte(placeholder)
	return tlv.MakePrimitiveRecord(PlaceholderTLVType, &byteHolder)
}

// EncodedPlaceHolderTLV generates a placeholder tlv, return a byte slice representing tlv
func EncodedPlaceHolderTLV(placeholder string) []byte {
	placeholderTlv := PlaceHolderTLV(placeholder)
	return util.TLVRecordToBytes(placeholderTlv)
}
