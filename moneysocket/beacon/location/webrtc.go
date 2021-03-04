package location

import (
	"fmt"

	util2 "github.com/xplorfin/moneysocket-go/moneysocket/util"

	"github.com/lightningnetwork/lnd/tlv"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
)

// see https://github.com/moneysocket/js-moneysocket/blob/76e533b59df1fcf03bd09c3e11813f016811fb71/moneysocket/beacon/location/webrtc.js#L21

const (
	DefaultWebrtcPlaceholder = "webrtc herpader"
	webrtcType               = "WebRTC"
)

type WebrtcLocation struct {
	PlaceholderString string
}

// statically assert that this type binds to location interface
var _ Location = WebrtcLocation{}

func NewWebRtcLocation() WebrtcLocation {
	return WebrtcLocation{
		PlaceholderString: DefaultWebrtcPlaceholder,
	}
}

func (loc WebrtcLocation) Type() tlv.Type {
	return util.WebrtcLocationTlvType
}

func (loc WebrtcLocation) Tlv() []tlv.Record {
	placeHolder := EncodedPlaceHolderTlv(loc.PlaceholderString)
	return []tlv.Record{tlv.MakePrimitiveRecord(util.WebrtcLocationTlvType, &placeHolder)}
}

func (loc WebrtcLocation) EncodedTlv() []byte {
	res := loc.Tlv()
	return util2.TlvRecordToBytes(res...)
}

func (loc WebrtcLocation) ToObject() map[string]interface{} {
	m := make(map[string]interface{})
	m["type"] = webrtcType
	m["placeholder_string"] = loc.PlaceholderString
	return m
}

func WebRtcLocationFromTlv(tlv util.Tlv) (loc WebrtcLocation, err error) {
	if tlv.Type() != util.WebrtcLocationTlvType {
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
