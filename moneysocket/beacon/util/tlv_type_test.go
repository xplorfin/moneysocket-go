package util

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/lightningnetwork/lnd/tlv"
)

// types cannot be less than 2^16 according to bolt-1 https://git.io/JLCRq
var Bolt1MinTlv = tlv.Type(math.Pow(2, 16))

func TestBolt1MinSpec(t *testing.T) {
	if Bolt1MinTlv != TLVMinimum {
		t.Errorf("bolt-1 specifies a minimum of 2^16 (%d), got %d", Bolt1MinTlv, TLVMinimum)
	}
}

type TlvSlice []tlv.Type

// check if a tlv slice contains a tlv
func (s TlvSlice) Contains(p tlv.Type) bool {
	for _, t := range s {
		if t == p {
			return true
		}
	}
	return false
}

func TestTypes(t *testing.T) {
	prevItems := make(TlvSlice, 0)
	for i, tlvType := range TlvTypes {
		if prevItems.Contains(tlvType) {
			t.Errorf("tlv's are expected to be unique, found duplciate for value %d", tlvType)
		}

		// this is slightly redundant
		if tlvType < Bolt1MinTlv {
			t.Errorf("bolt-1 specifies a minimum of 2^16 (%d), got %d for %s item in array", Bolt1MinTlv, tlvType, humanize.Ordinal(i+1))
		}

		if tlvType < TLVMinimum {
			t.Errorf("the socket.money specifies a minimum of 2^16+443 (%d), got %d for %s item in array", TLVMinimum, TLVMinimum, humanize.Ordinal(i))
		}
	}
}

func TestTlvTest(t *testing.T) {
	test := []byte("random string")
	var test1 uint32 = 2141242112
	tlvStream, err := tlv.NewStream(
		tlv.MakePrimitiveRecord(WebsocketLocationTLVType, &test1),
		tlv.MakePrimitiveRecord(WebRTCLocationTLVType, &test),
	)
	if err != nil {
		t.Error(err)
	}
	var w bytes.Buffer
	err = tlvStream.Encode(&w)
	if err != nil {
		panic(err)
	}
	fmt.Println(hex.EncodeToString(w.Bytes()))
}
