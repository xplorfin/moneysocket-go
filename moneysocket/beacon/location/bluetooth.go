package location

import (
	"fmt"

	util2 "github.com/xplorfin/moneysocket-go/moneysocket/util"

	tlvHelper "github.com/lightningnetwork/lnd/tlv"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
)

// see https://github.com/moneysocket/js-moneysocket/blob/76e533b59df1fcf03bd09c3e11813f016811fb71/moneysocket/beacon/location/bluetooth.js#L21

const (
	DEFAULT_BLUETOOTH_PLACEHOLDER = "bluetooth herpader"
	BLUETOOTH_TYPE                = "Bluetooth"
)

// BluetoothLocation beacon - this is not currently implemented and is reserved for future use
type BluetoothLocation struct {
	PlaceholderString string
}

// Create a new bluetooth location with the default header
func NewBluetoothLocation() BluetoothLocation {
	return BluetoothLocation{
		PlaceholderString: DEFAULT_BLUETOOTH_PLACEHOLDER,
	}
}

// Get the BluetoothLocation tlv type
func (loc BluetoothLocation) Type() tlvHelper.Type {
	return util.BluetoothLocationTlvType
}

// get the encoded tlv of a given BluetoothLocation
func (loc BluetoothLocation) Tlv() []tlvHelper.Record {
	placeHolder := EncodedPlaceHolderTlv(loc.PlaceholderString)
	return []tlvHelper.Record{tlvHelper.MakePrimitiveRecord(util.BluetoothLocationTlvType, &placeHolder)}
}

// get the encoded tlv of a given BluetoothLocation
func (loc BluetoothLocation) EncodedTlv() []byte {
	res := loc.Tlv()
	return util2.TlvRecordToBytes(res...)
}

// encode the BluetoothLocation to a json-serializable map
func (loc BluetoothLocation) ToObject() map[string]interface{} {
	m := make(map[string]interface{})
	m["type"] = BLUETOOTH_TYPE
	m["placeholder_string"] = loc.PlaceholderString
	return m
}

// decode a BluetoothLocation from a tlv object
func BluetoothLocationFromTlv(tlv util.Tlv) (loc BluetoothLocation, err error) {
	if tlv.Type() != util.BluetoothLocationTlvType {
		return loc, fmt.Errorf("got unexpected tlv type: %d expected %d", tlv.Type(), util.BluetoothLocationTlvType)
	}
	tlvs, err := util.NamespaceIterTlvs(tlv.Value())
	if err != nil {
		return loc, err
	}
	// wether or not the tlv has a placeholder string
	hasPlaceholder := false
	for _, subTlv := range tlvs {
		if subTlv.Type() == PLACEHOLDER_TLV_TYPE {
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
var _ Location = BluetoothLocation{}
