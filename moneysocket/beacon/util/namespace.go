package util

import "fmt"

// NamespacePopTLVs Represents a specific namespace of TLVs as referred to in BOLT 1 and
// provides generic pop helpers for the fundamental types defined here:
// https://git.io/JmwOJ and https://git.io/JmwOk
func NamespacePopTLVs(byteString []byte) (t TLV, remainder []byte, err error) {
	return TLVPop(byteString)
}

// NamespaceTLVsAreValid will determine if a bytestring contains valid tlvs.
func NamespaceTLVsAreValid(byteString []byte) bool {
	bs := byteString
	for len(bs) > 0 {
		_, byteString, err := NamespacePopTLVs(bs)
		bs = byteString
		if err != nil {
			return false
		}
	}
	return true
}

// NamespaceIterTLVs will iterate over the namespaces in a byteString
// returns error if all are not valid.
func NamespaceIterTLVs(byteString []byte) (tlvs []TLV, err error) {
	if !NamespaceTLVsAreValid(byteString) {
		return tlvs, fmt.Errorf("namespace tlvs may not be valid")
	}
	bs := byteString
	for len(bs) > 0 {
		// we validated error above
		tlv, byteString, _ := NamespacePopTLVs(bs)
		bs = byteString
		tlvs = append(tlvs, tlv)
	}
	return tlvs, err
}
