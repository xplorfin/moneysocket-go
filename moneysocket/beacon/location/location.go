package location

import "github.com/lightningnetwork/lnd/tlv"

// Location defines an interface for various Location types.
type Location interface {
	// Type gets the tlv.Type for a given location
	Type() tlv.Type
	// TLV returns the encoded tlv.Record for a given location
	TLV() []tlv.Record
	// EncodedTLV gets the encoded tlv as a byte slice for a given location
	EncodedTLV() []byte
	// ToObject gets a json-serializable object from the Location
	ToObject() map[string]interface{}
}
