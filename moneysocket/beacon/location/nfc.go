package location

import (
	"fmt"

	moneysocketUtil "github.com/xplorfin/moneysocket-go/moneysocket/util"

	"github.com/lightningnetwork/lnd/tlv"
	beaconUtil "github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
)

// see https://git.io/JmVpX

const (
	// DefaultNFCPlaceholder defines the default value used by a placeholder
	DefaultNFCPlaceholder = "nfc herpader"
	// NFCType defines the name of the NFCLocation type
	NFCType = "NFC"
)

// NFCLocation defines a Location type for nfc
type NFCLocation struct {
	PlaceholderString string
}

// NewNFCLocation creates a new NFCLocation
func NewNFCLocation() NFCLocation {
	return NFCLocation{
		PlaceholderString: DefaultNFCPlaceholder,
	}
}

// Type gets the type of NFCLocation
func (loc NFCLocation) Type() tlv.Type {
	return beaconUtil.NFCLocationTLVType
}

// TLV gets the tlv of a given NFCLocation
func (loc NFCLocation) TLV() []tlv.Record {
	placeHolder := EncodedPlaceHolderTLV(loc.PlaceholderString)
	return []tlv.Record{tlv.MakePrimitiveRecord(beaconUtil.NFCLocationTLVType, &placeHolder)}
}

// EncodedTLV gets the encoded tlv of a given NFCLocation
func (loc NFCLocation) EncodedTLV() []byte {
	res := loc.TLV()
	return moneysocketUtil.TLVRecordToBytes(res...)
}

// ToObject converts the NFCLocation to a json-serializable map
func (loc NFCLocation) ToObject() map[string]interface{} {
	m := make(map[string]interface{})
	m["type"] = NFCType
	m["placeholder_string"] = loc.PlaceholderString
	return m
}

// NfcLocationFromTLV converts a util.TLV from a tlv object
func NfcLocationFromTLV(tlv beaconUtil.TLV) (loc NFCLocation, err error) {
	if tlv.Type() != beaconUtil.NFCLocationTLVType {
		return loc, fmt.Errorf("got unexpected tlv type: %d expected %d", tlv.Type(), beaconUtil.BluetoothLocationTLVType)
	}
	tlvs, err := beaconUtil.NamespaceIterTLVs(tlv.Value())
	if err != nil {
		return loc, err
	}
	// wether or not the tlv has a placeholder string
	hasPlaceholder := false
	for _, subTlv := range tlvs {
		if subTlv.Type() == PlaceholderTLVType {
			loc.PlaceholderString = string(subTlv.Value())
			hasPlaceholder = true
		}
	}
	if !hasPlaceholder {
		return loc, fmt.Errorf("error decoding placeholder string")
	}
	return loc, err
}

// statically assert that this type binds to location interface
var _ Location = NFCLocation{}
