package location

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
)

type webrtcTestCase struct {
	PlaceholderString string
	EncodedTlv        string
}

var webrtcTestCases = []webrtcTestCase{
	{
		PlaceholderString: DefaultWebRTCPlaceholder,
		EncodedTlv:        "fe000101c511000f776562727463206865727061646572",
	},
}

func TestWebrtcEncoding(t *testing.T) {
	for _, testCase := range webrtcTestCases {
		rlc := NewWebRTCLocation()
		if hex.EncodeToString(rlc.EncodedTLV()) != testCase.EncodedTlv {
			t.Errorf("expected %s to equal %s", hex.EncodeToString(rlc.EncodedTLV()), testCase.EncodedTlv)
		}
		// fetch encoded tlv
		decoded, err := hex.DecodeString(testCase.EncodedTlv)
		if err != nil {
			t.Error(err)
		}
		tlv, _, err := util.TLVPop(decoded)
		if err != nil {
			t.Error(err)
		}
		// try to decode tlv
		loc, err := WebRTCLocationFromTLV(tlv)
		if err != nil {
			t.Error(err)
		}
		if loc != rlc {
			t.Error("expected encoded and decoded tlvs to be identical")
		}

		// compare gneerated objects
		if !reflect.DeepEqual(rlc.ToObject(), loc.ToObject()) {
			t.Error("expected encoded and decoded tlvs to be identical")
		}
	}
}
