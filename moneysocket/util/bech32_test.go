package util

import (
	"bytes"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

// make sure bech32 encodes
func TestBech32BidrectionalEncoding(t *testing.T) {
	gofakeit.Seed(0)
	for i := 0; i < 10000; i++ {
		rawHrp := gofakeit.Word()
		rawKey := []byte(gofakeit.BitcoinPrivateKey()[0 : 52-len(rawHrp)])
		// assert arbitary word lenghts don't break the decoding mechanis
		encodedKey, err := Bech32EncodeBytes(rawKey, rawHrp)
		if err != nil {
			t.Error(err)
		}
		decodedHrp, decodedKey, err := Bech32DecodeBytes(encodedKey)
		if err != nil {
			t.Error(err)
		}
		if decodedHrp != rawHrp {
			t.Errorf("expected decodedHrp %s to match decoded decodedHrp %s", rawHrp, decodedHrp)
		}
		if !bytes.Equal(rawKey, decodedKey) {
			t.Errorf("expected decodedHrp %s to match decoded decodedHrp %s", rawKey, decodedKey)
		}
	}
}

// test data that cannot be a bech32
func TestUndecodable(t *testing.T) {
	gofakeit.Seed(0)
	for i := 0; i < 10000; i++ {
		//without a rawHrp these should fail
		encodeTest := gofakeit.BitcoinAddress()
		// assert arbitary word lenghts don't break the decoding mechanis
		decodedHrp, decodedKey, err := Bech32DecodeBytes(encodeTest)
		if err == nil {
			t.Errorf("expected error from hrp %s and valid key: %s", decodedHrp, decodedKey)
		}
	}
}
