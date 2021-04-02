package location

import (
	"fmt"

	tlvHelper "github.com/lightningnetwork/lnd/tlv"
	beaconUtil "github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
	moneysocketUtil "github.com/xplorfin/moneysocket-go/moneysocket/util"
)

// see https://git.io/JmVpM

const (
	// DefaultBluetoothPlaceholder temporarily defines an hrp header
	// for the as-of-yet unused BluetoothLocation type.
	DefaultBluetoothPlaceholder = "bluetooth herpader"
	// BluetoothType defines the name of the BluetoothLocation type.
	BluetoothType = "Bluetooth"
)

// BluetoothLocation beacon - this is not currently implemented and is reserved for future use.
type BluetoothLocation struct {
	// PlaceholderString is the
	PlaceholderString string
}

// NewBluetoothLocation Creates a new bluetooth location with the default header.
func NewBluetoothLocation() BluetoothLocation {
	return BluetoothLocation{
		PlaceholderString: DefaultBluetoothPlaceholder,
	}
}

// Type gets the BluetoothLocation tlv type.
func (loc BluetoothLocation) Type() tlvHelper.Type {
	return beaconUtil.BluetoothLocationTLVType
}

// TLV gets the tlv of a given BluetoothLocation.
func (loc BluetoothLocation) TLV() []tlvHelper.Record {
	placeHolder := EncodedPlaceHolderTLV(loc.PlaceholderString)
	return []tlvHelper.Record{tlvHelper.MakePrimitiveRecord(beaconUtil.BluetoothLocationTLVType, &placeHolder)}
}

// EncodedTLV gets the encoded tlv of a given BluetoothLocation.
func (loc BluetoothLocation) EncodedTLV() []byte {
	res := loc.TLV()
	return moneysocketUtil.TLVRecordToBytes(res...)
}

// ToObject converts the BluetoothLocation to a json-serializable map.
func (loc BluetoothLocation) ToObject() map[string]interface{} {
	m := make(map[string]interface{})
	m["type"] = BluetoothType
	m["placeholder_string"] = loc.PlaceholderString
	return m
}

// BluetoothLocationFromTLV converts a util.TLV from a tlv object.
func BluetoothLocationFromTLV(tlv beaconUtil.TLV) (loc BluetoothLocation, err error) {
	if tlv.Type() != beaconUtil.BluetoothLocationTLVType {
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

// statically assert that this type binds to location interface.
var _ Location = BluetoothLocation{}
