package location

import "github.com/lightningnetwork/lnd/tlv"

type Location interface {
	// get the tlv.Type for a given location
	Type() tlv.Type
	// return the encoded tlv.Record for a given location
	Tlv() []tlv.Record
	// get the encoded tlv as a byte slice for a given location
	EncodedTlv() []byte
	// get a json-serializable object from the location
	ToObject() map[string]interface{}
}
