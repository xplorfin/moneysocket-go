package location

import (
	"encoding/hex"
	"net/url"
	"reflect"
	"testing"

	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
)

type WebsocketTestCase struct {
	Host       string
	UseTLS     bool
	EncodedTlv string
	Port       *int
}

// tested against python. Approximately:
// testCase = WebsocketLocation(host="relay.socket", use_tls=True, port=50)
// print(testCase.encode_tlv().hex()).
var testPort = 50

var websocketTestCases = []WebsocketTestCase{
	{
		Host:       "relay.socket.money",
		UseTLS:     true,
		EncodedTlv: "fe000101c314001272656c61792e736f636b65742e6d6f6e6579",
	},
	{
		Host:       "relay.socket.money",
		UseTLS:     false,
		EncodedTlv: "fe000101c317001272656c61792e736f636b65742e6d6f6e6579010100",
	},
	{
		Host:       "relay.socket",
		UseTLS:     false,
		EncodedTlv: "fe000101c314000c72656c61792e736f636b6574010100020132",
		Port:       &testPort,
	},
	{
		Host:       "relay.socket",
		UseTLS:     true,
		EncodedTlv: "fe000101c311000c72656c61792e736f636b6574020132",
		Port:       &testPort,
	},
}

func TestWebsocketEncoding(t *testing.T) {
	for _, testCase := range websocketTestCases {
		// encode websockets
		var ws WebsocketLocation
		if testCase.Port == nil {
			ws = NewWebsocketLocation(testCase.Host, testCase.UseTLS)
		} else {
			ws = NewWebsocketLocationPort(testCase.Host, testCase.UseTLS, *testCase.Port)
		}
		if hex.EncodeToString(ws.EncodedTLV()) != testCase.EncodedTlv {
			t.Errorf("expected tlv %s to equal %s", hex.EncodeToString(ws.EncodedTLV()), testCase.EncodedTlv)
		}

		// make sure produced uri is avlid
		_, err := url.ParseRequestURI(ws.ToString())
		if err != nil {
			t.Error(err)
		}

		if ws.IsTLS() != testCase.UseTLS {
			t.Errorf("expected %v to equal %v", ws.IsTLS(), testCase.UseTLS)
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
		loc, err := WebsocketLocationFromTLV(tlv)
		if err != nil {
			t.Error(err)
		}
		if loc != ws {
			t.Error("expected encoded and decoded tlvs to be identical")
		}

		// compare gneerated objects
		if !reflect.DeepEqual(loc.ToObject(), ws.ToObject()) {
			t.Error("expected encoded and decoded tlvs to be identical")
		}
	}
}
