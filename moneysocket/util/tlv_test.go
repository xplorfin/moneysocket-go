package util

import (
	"bytes"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/dustin/go-humanize"
	"github.com/kr/pretty"
	"github.com/lightningnetwork/lnd/tlv"
)

// asserts tlv encoder that produces bytes works correctly.
func TestTlvRecordToBytes(t *testing.T) {
	gofakeit.Seed(0)
	if !PrngIsAvailable() {
		t.Skip("running this test requires PRNG is available on your machine")
	}
	for i := 1; i < 100; i++ {
		tlvBytes, err := GenerateRandomBytes(i)
		if err != nil {
			panic(err)
		}
		encoder := tlv.StubEncoder(tlvBytes)
		record := tlv.MakeStaticRecord(1, nil, 4, encoder, nil)
		producedBytes := TLVRecordToBytes(record)
		if !bytes.Equal(tlvBytes, producedBytes[2:]) {
			t.Errorf("bytes not correctly produced on %s run", humanize.Ordinal(i))
			t.Error(pretty.Diff(tlvBytes, producedBytes))
		}
	}
}
