package location

import "github.com/lightningnetwork/lnd/tlv"

type Location interface {
	Type() tlv.Type
	Tlv() []tlv.Record
	EncodedTlv() []byte
	ToObject() map[string]interface{}
}
