package beacon

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	util2 "github.com/xplorfin/moneysocket-go/moneysocket/util"

	"github.com/lightningnetwork/lnd/tlv"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
)

const SharedSeedLength = 16

type SharedSeed struct {
	seedBytes []byte
}

// generate a new random shared seed
func NewSharedSeed() SharedSeed {
	randBytes, err := util2.GenerateRandomBytes(16)
	if err != nil {
		panic(err)
	}
	seed, err := BytesToSharedSeed(randBytes)
	if err != nil {
		panic(err)
	}
	return seed
}

// create a shared seed from a byte slice
func BytesToSharedSeed(rawBytes []byte) (seed SharedSeed, err error) {
	if len(rawBytes) != SharedSeedLength {
		return seed, fmt.Errorf("byte slice of length %d does not match expected %d", len(rawBytes), SharedSeedLength)
	}
	return SharedSeed{seedBytes: rawBytes}, nil
}

// create a shared seed from a hex
func HexToSharedSeed(rawHex string) (seed SharedSeed, err error) {
	if len(rawHex) != SharedSeedLength*2 {
		return seed, fmt.Errorf("hex of length %d is invalid, expected length of %d", len(rawHex), SharedSeedLength*2)
	}
	rawBytes, err := hex.DecodeString(rawHex)
	if err != nil {
		return seed, err
	}
	return BytesToSharedSeed(rawBytes)
}

// get bytes
func (s SharedSeed) GetBytes() []byte {
	return s.seedBytes
}

// convert ints to hash
func (s SharedSeed) Hash() uint64 {
	return binary.BigEndian.Uint64(s.seedBytes)
}

// check if two seeds are equal
func (s SharedSeed) Equal(seed SharedSeed) bool {
	return bytes.Equal(s.seedBytes, seed.seedBytes)
}

// convert shared seed to string
func (s SharedSeed) Hex() string {
	return hex.EncodeToString(s.seedBytes)
}

// wrapper around hex
func (s SharedSeed) ToString() string {
	return s.Hex()
}

// TODO: find out why these are part of this struct
// sha256 the seed bytes, format as string
func (s SharedSeed) Sha256(inputBytes []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(inputBytes))
}

// TODO: find out why these are part of this struct
// double sha 256
func (s SharedSeed) DoubleSha256(inputBytes []byte) []byte {
	return util2.CreateDoubleSha256(inputBytes)
}

// returns a derived aes256 key
// wrapper around DoubleSha256
func (s SharedSeed) DeriveAes256Key() []byte {
	return s.DoubleSha256(s.seedBytes)
}

// create a double sha256 of the derived aes 256 key (itself a double sha-256)
func (s SharedSeed) DeriveRendezvousId() []byte {
	return util2.CreateDoubleSha256(s.DeriveAes256Key())
}

// encode the tlv into a static record
func (s SharedSeed) TLV() tlv.Record {
	return tlv.MakeStaticRecord(util.SharedSeedTlvType, &s.seedBytes, 16, tlv.EVarBytes, tlv.DVarBytes)
}

// encode the static record into a
func (s SharedSeed) EncodedTLV() []byte {
	return util2.TlvRecordToBytes(s.TLV())
}
