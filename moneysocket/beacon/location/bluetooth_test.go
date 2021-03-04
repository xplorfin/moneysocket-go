package location

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
)

type BluetoothTestCase struct {
	PlaceholderString string
	EncodedTlv        string
}

var bluetoothTestCase = []BluetoothTestCase{
	{
		PlaceholderString: DEFAULT_BLUETOOTH_PLACEHOLDER,
		EncodedTlv:        "fe000101c7140012626c7565746f6f7468206865727061646572",
	},
}

func TestBluetoothEncoding(t *testing.T) {
	for _, testCase := range bluetoothTestCase {
		blc := NewBluetoothLocation()
		if hex.EncodeToString(blc.EncodedTlv()) != testCase.EncodedTlv {
			t.Errorf("expected %s to equal %s", hex.EncodeToString(blc.EncodedTlv()), testCase.EncodedTlv)
		}
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
		loc, err := BluetoothLocationFromTlv(tlv)
		if err != nil {
			t.Error(err)
		}
		if loc != blc {
			t.Error("expected encoded and decoded tlvs to be identical")
		}

		// compare gneerated objects
		if !reflect.DeepEqual(blc.ToObject(), loc.ToObject()) {
			t.Error("expected encoded and decoded tlvs to be identical")
		}
	}
}
