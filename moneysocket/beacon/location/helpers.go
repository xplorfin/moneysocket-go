package location

import (
	"github.com/lightningnetwork/lnd/tlv"
	"github.com/xplorfin/moneysocket-go/moneysocket/util"
)

// default tlv type for unimplemented location beacons
const PLACEHOLDER_TLV_TYPE = 0

// encode a placeholder tlv for python-parity tests
func PlaceHolderTlv(placeholder string) tlv.Record {
	byteHolder := []byte(placeholder)
	return tlv.MakePrimitiveRecord(PLACEHOLDER_TLV_TYPE, &byteHolder)
}

// generate a placeholder tlv, return a byte slice representing tlv
func EncodedPlaceHolderTlv(placeholder string) []byte {
	placeholderTlv := PlaceHolderTlv(placeholder)
	return util.TlvRecordToBytes(placeholderTlv)
}
