package location

import (
	"fmt"

	util2 "github.com/xplorfin/moneysocket-go/moneysocket/util"

	"github.com/lightningnetwork/lnd/tlv"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
)

// see https://git.io/JmwkY

const (
	// DefaultWebRTCPlaceholder temporarily defines an hrp header
	// for the as-of-yet unused WebRTCLocation type
	DefaultWebRTCPlaceholder = "webrtc herpader"
	// WebRTCType defines a type
	WebRTCType = "WebRTC"
)

// WebRTCLocation is a beacon location type placeholder
type WebRTCLocation struct {
	// PlaceholderString defines what gets encoded
	PlaceholderString string
}

// statically assert that this type binds to location interface
var _ Location = WebRTCLocation{}

// NewWebRTCLocation creates a new WebRTCLocation
func NewWebRTCLocation() WebRTCLocation {
	return WebRTCLocation{
		PlaceholderString: DefaultWebRTCPlaceholder,
	}
}

// Type gets the WebRTCLocation tlv.Type
func (loc WebRTCLocation) Type() tlv.Type {
	return util.WebRTCLocationTLVType
}

// TLV gets the tlv for a WebRTCLocation
func (loc WebRTCLocation) TLV() []tlv.Record {
	placeHolder := EncodedPlaceHolderTLV(loc.PlaceholderString)
	return []tlv.Record{tlv.MakePrimitiveRecord(util.WebRTCLocationTLVType, &placeHolder)}
}

// EncodedTLV gets the encoded tlv of a WebRTCLocation
func (loc WebRTCLocation) EncodedTLV() []byte {
	res := loc.TLV()
	return util2.TLVRecordToBytes(res...)
}

// ToObject converts WebRTCLocation to a json-serializable map
func (loc WebRTCLocation) ToObject() map[string]interface{} {
	m := make(map[string]interface{})
	m["type"] = WebRTCType
	m["placeholder_string"] = loc.PlaceholderString
	return m
}

// WebRTCLocationFromTLV returns a WebRTCLocation based on an input tlv
// returns an error
func WebRTCLocationFromTLV(tlv util.TLV) (loc WebRTCLocation, err error) {
	if tlv.Type() != util.WebRTCLocationTLVType {
		return loc, fmt.Errorf("got unexpected tlv type: %d expected %d", tlv.Type(), util.BluetoothLocationTLVType)
	}
	tlvs, err := util.NamespaceIterTLVs(tlv.Value())
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
