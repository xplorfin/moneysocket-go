package util

import "fmt"

// Represents a specific namespace of TLVs as referred to in BOLT 1 and
// provides generic pop helpers for the fundamental types defined here:
// https://github.com/lightningnetwork/lightning-rfc/blob/master/01-messaging.md#fundamental-types
// see: https://github.com/moneysocket/py-moneysocket/blob/main/moneysocket/utl/third_party/bolt/namespace.py#L9
func NamespacePopTlvs(byteString []byte) (t Tlv, remainder []byte, err error) {
	return TlvPop(byteString)
}

// assert tlvs are valid
func NamespaceTlvsAreValid(byteString []byte) bool {
	bs := byteString
	for len(bs) > 0 {
		_, byteString, err := NamespacePopTlvs(bs)
		bs = byteString
		if err != nil {
			return false
		}
	}
	return true
}

func NamespaceIterTlvs(byteString []byte) (tlvs []Tlv, err error) {
	if !NamespaceTlvsAreValid(byteString) {
		return tlvs, fmt.Errorf("namespace tlvs may not be valid")
	}
	bs := byteString
	for len(bs) > 0 {
		// we validated error above
		tlv, byteString, _ := NamespacePopTlvs(bs)
		bs = byteString
		tlvs = append(tlvs, tlv)
	}
	return tlvs, err
}
