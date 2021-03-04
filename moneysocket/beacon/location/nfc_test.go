package location

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
)

type NfcTestCase struct {
	PlaceholderString string
	EncodedTlv        string
}

var nfcTestCases = []NfcTestCase{
	{
		PlaceholderString: DefaultNfcPlaceholder,
		EncodedTlv:        "fe000101c90e000c6e6663206865727061646572",
	},
}

func TestNfcEncoding(t *testing.T) {
	for _, testCase := range nfcTestCases {
		nlc := NewNfcLocation()
		if hex.EncodeToString(nlc.EncodedTlv()) != testCase.EncodedTlv {
			t.Errorf("expected %s to equal %s", hex.EncodeToString(nlc.EncodedTlv()), testCase.EncodedTlv)
		}
		// TODO
		nlc.ToObject()

		// fetch encoded tlv
		decoded, err := hex.DecodeString(testCase.EncodedTlv)
		if err != nil {
			t.Error(err)
		}
		tlv, _, err := util.TlvPop(decoded)
		if err != nil {
			t.Error(err)
		}
		// try to decode tlv
		loc, err := NfcLocationFromTlv(tlv)
		if err != nil {
			t.Error(err)
		}
		if loc != nlc {
			t.Error("expected encoded and decoded tlvs to be identical")
		}

		// compare generated objects
		if !reflect.DeepEqual(nlc.ToObject(), loc.ToObject()) {
			t.Error("expected encoded and decoded tlvs to be identical")
		}
	}
}
