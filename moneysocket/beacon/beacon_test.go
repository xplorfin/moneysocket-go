package beacon

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/lightningnetwork/lnd/tlv"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/location"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
	util2 "github.com/xplorfin/moneysocket-go/moneysocket/util"
)

type BeaconTestCases struct {
	// Seed of string
	Seed []byte
	// Bech32String that should result
	Bech32String string
	// Locations list we use for the beacon
	Locations []location.Location
}

// TestCases defines a list of test-cases to ensure python parity.
var TestCases = []BeaconTestCases{
	{
		Seed:         []byte("hello from light"),
		Bech32String: "moneysocket1lcqqzqdmy0lqqqgph5gxsetvd3hjqenjdakjqmrfva58flsqqyquzpl7qqqsr0cpqqpwt49e",
	},
	{
		Seed:         []byte("hello from light"),
		Bech32String: "moneysocket1lcqqzqdmstlqqqgph5gxsetvd3hjqenjdakjqmrfva58flsqqyquzeh7qqqsr0cpqnlqqqgpcv2qqynjv4kxz7fwwdhkx6m9wshx6mmwv4uluqqpq8z3zqq0wajkyun5vvsxsetjwpskgetjlcqqzqw8zsqpycnvw4jhgmm0w35zq6r9wfcxzer9wtlqqqgpey8qqrrwve3jq6r9wfcxzer9wgmv4l2e",
		Locations: []location.Location{
			location.NewWebsocketLocation("relay.socket.money", true),
			location.NewWebRTCLocation(),
			location.NewBluetoothLocation(),
			location.NewNFCLocation(),
		},
	},
}

// TestBeaconParity tests a generated bech32 string against python
// TODO this can be automated with a python (or js) harness.
func TestBeaconParity(t *testing.T) {
	for _, testCase := range TestCases {
		ss, err := BytesToSharedSeed(testCase.Seed)
		if err != nil {
			t.Error(err)
		}

		beacon := NewBeaconFromSharedSeed(ss)

		for _, loc := range testCase.Locations {
			beacon.AddLocation(loc)
		}

		if beacon.ToBech32Str() != testCase.Bech32String {
			t.Errorf("expected bech32 string %s to match %s", beacon.ToBech32Str(), testCase.Bech32String)
		}

		bec, err := DecodeFromBech32Str(testCase.Bech32String)
		if beacon.ToBech32Str() != bec.ToBech32Str() {
			t.Error(err)
		}
		if err != nil {
			t.Error(err)
		}
	}
}

// TestRecord tests tlv encoding.
func TestRecord(t *testing.T) {
	ss, err := BytesToSharedSeed([]byte("hello from light"))
	if err != nil {
		t.Error(err)
	}
	record := tlv.MakeStaticRecord(util.SharedSeedTLVType, &ss.seedBytes, uint64(len(ss.seedBytes)), tlv.EVarBytes, tlv.DVarBytes)
	res, err := tlv.NewStream(record)
	if err != nil {
		panic(err)
	}
	var w bytes.Buffer
	err = res.Encode(&w)
	if err != nil {
		panic(err)
	}
	h := w.Bytes()
	pj := hex.EncodeToString(h) // this is equal to the output of encode_tlvs in python
	if pj != hex.EncodeToString(util2.TLVRecordToBytes(record)) {
		t.Error("oops")
	}
}
