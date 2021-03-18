package location

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
)

// NFCTestCase is a test case to test out encoding/decoding an NFCLocation
type NFCTestCase struct {
	// PlaceholderString contains the placeholder string we test against
	PlaceholderString string
	// EncodedTLV contains the encoded tlv we're trying to parity
	EncodedTLV string
}

var nfcTestCases = []NFCTestCase{
	{
		PlaceholderString: DefaultNFCPlaceholder,
		EncodedTLV:        "fe000101c90e000c6e6663206865727061646572",
	},
}

// TestNFCEncoding makes sure that NFCLocation has parity with python test cases
func TestNFCEncoding(t *testing.T) {
	for _, testCase := range nfcTestCases {
		nlc := NewNFCLocation()
		if hex.EncodeToString(nlc.EncodedTLV()) != testCase.EncodedTLV {
			t.Errorf("expected %s to equal %s", hex.EncodeToString(nlc.EncodedTLV()), testCase.EncodedTLV)
		}
		// TODO
		nlc.ToObject()

		// fetch encoded tlv
		decoded, err := hex.DecodeString(testCase.EncodedTLV)
		if err != nil {
			t.Error(err)
		}
		tlv, _, err := util.TLVPop(decoded)
		if err != nil {
			t.Error(err)
		}
		// try to decode tlv
		loc, err := NfcLocationFromTLV(tlv)
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
