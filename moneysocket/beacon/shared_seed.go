package beacon

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/lightningnetwork/lnd/tlv"
	beaconUtil "github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
	moneysocketUtil "github.com/xplorfin/moneysocket-go/moneysocket/util"
)

// SharedSeedLength defines the length of the shared seed.
const SharedSeedLength = 16

// SharedSeed is a seed used for end-to-end encryption.
type SharedSeed struct {
	seedBytes []byte
}

// NewSharedSeed generates a new random shared seed.
func NewSharedSeed() SharedSeed {
	randBytes, err := moneysocketUtil.GenerateRandomBytes(16)
	if err != nil {
		panic(err)
	}
	seed, err := BytesToSharedSeed(randBytes)
	if err != nil {
		panic(err)
	}
	return seed
}

// BytesToSharedSeed creates a SharedSeed from a []byte slice
// returns an error when seed cannot be decoded.
func BytesToSharedSeed(rawBytes []byte) (seed SharedSeed, err error) {
	if len(rawBytes) != SharedSeedLength {
		return seed, fmt.Errorf("byte slice of length %d does not match expected %d", len(rawBytes), SharedSeedLength)
	}
	return SharedSeed{seedBytes: rawBytes}, nil
}

// HexToSharedSeed creates a shared seed from an input hex
// returns an error when this is not possible.
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

// GetBytes gets the seedBytes from the shared seed.
func (s SharedSeed) GetBytes() []byte {
	return s.seedBytes
}

// Hash converts the seed bytes to a big endian int.
func (s SharedSeed) Hash() uint64 {
	return binary.BigEndian.Uint64(s.seedBytes)
}

// Equal check if two seeds are equal.
func (s SharedSeed) Equal(seed SharedSeed) bool {
	return bytes.Equal(s.seedBytes, seed.seedBytes)
}

// Hex converts the shared seed to string via hex encoding.
func (s SharedSeed) Hex() string {
	return hex.EncodeToString(s.seedBytes)
}

// ToString is a wrapper around Hex() for usability.
func (s SharedSeed) ToString() string {
	return s.Hex()
}

// SHA256 the seed bytes, format as string.
func (s SharedSeed) SHA256(inputBytes []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(inputBytes))
}

// DoubleSHA256 will generate SHA256(SHA256(inputBytes)).
func (s SharedSeed) DoubleSHA256(inputBytes []byte) []byte {
	return moneysocketUtil.CreateDoubleSha256(inputBytes)
}

// DeriveAES256Key will generate a DoubleSHA256 of the shared seed
// see https://en.wikipedia.org/wiki/Advanced_Encryption_Standard for details.
func (s SharedSeed) DeriveAES256Key() []byte {
	return s.DoubleSHA256(s.seedBytes)
}

// DeriveRendezvousID will generate a DoubleSHA256 of the DeriveAES256Key.
func (s SharedSeed) DeriveRendezvousID() []byte {
	return moneysocketUtil.CreateDoubleSha256(s.DeriveAES256Key())
}

// TLV will encode the SharedSeed into a tlv.Record.
func (s SharedSeed) TLV() tlv.Record {
	return tlv.MakeStaticRecord(beaconUtil.SharedSeedTLVType, &s.seedBytes, 16, tlv.EVarBytes, tlv.DVarBytes)
}

// EncodedTLV will encode the TLV into a byte-slice
// (See BOLT #1: https://git.io/JLCRq ).
func (s SharedSeed) EncodedTLV() []byte {
	return moneysocketUtil.TLVRecordToBytes(s.TLV())
}
