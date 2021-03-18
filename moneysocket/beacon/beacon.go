package beacon

import (
	"fmt"

	"github.com/lightningnetwork/lnd/tlv"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/location"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util/bigsize"
	encodeUtils "github.com/xplorfin/moneysocket-go/moneysocket/util"
)

// MoneysocketHrp is the human readable part (https://en.bitcoin.it/wiki/BIP_0173#Bech32 )
// of the TLV (type-length-value, defined in BOLT #1: https://git.io/JLCRq )
const MoneysocketHrp = "moneysocket"

// Beacon contains a SharedSeed for end-to-end encryption and a list of location.Location's
type Beacon struct {
	seed      SharedSeed
	locations []location.Location
}

// NewBeacon creates a Beacon with no locations and an auto-generated SharedSeed
func NewBeacon() Beacon {
	return NewBeaconFromSharedSeed(NewSharedSeed())
}

// NewBeaconFromSharedSeed creates a Beacon with no locations
// and the given SharedSeed
func NewBeaconFromSharedSeed(seed SharedSeed) Beacon {
	return Beacon{
		seed: seed,
	}
}

// Locations returns the list of locations in the Beacon
func (b Beacon) Locations() []location.Location {
	return b.locations
}

// ToObject generates a json-encodable map of the Beacon's data
func (b Beacon) ToObject() map[string]interface{} {
	m := make(map[string]interface{})
	m["shared_seed"] = b.seed.ToString()
	locDict := make([]interface{}, 0)
	for _, loc := range b.locations {
		locDict = append(locDict, loc.ToObject())
	}
	m["locations"] = locDict
	return m
}

// AddLocation appends a location to the beacon
func (b *Beacon) AddLocation(loc location.Location) {
	b.locations = append(b.locations, loc)
}

// GetSharedSeed fetches the shared seed from the beacon
func (b Beacon) GetSharedSeed() SharedSeed {
	return b.seed
}

// EncodeLocationListTlvs encodes a list of locations as tlvs
func (b Beacon) EncodeLocationListTlvs() []byte {
	encoded := make([]byte, 0)
	locationCount := uint64(len(b.locations))
	record := tlv.MakeStaticRecord(util.LocationCountTLVType, &locationCount, tlv.VarIntSize(locationCount), util.EVarInt, util.DVarInt)
	locationCountEncoded := encodeUtils.TLVRecordToBytes(record)
	encoded = append(encoded, locationCountEncoded...)

	for _, loc := range b.locations {
		encoded = append(encoded, loc.EncodedTLV()...)
	}

	return encodeUtils.TLVRecordToBytes(tlv.MakePrimitiveRecord(util.LocationListTLVType, &encoded))
}

// EncodeTLV encodes the tlv
func (b Beacon) EncodeTLV() []byte {
	// encode the shared seed
	ssEncoded := b.seed.EncodedTLV()
	llEncoded := b.EncodeLocationListTlvs()
	combinedTlvs := append(ssEncoded, llEncoded...)
	return encodeUtils.TLVRecordToBytes(tlv.MakePrimitiveRecord(util.BeaconTLVType, &combinedTlvs))
}

// ToBech32Str encodes the tlv to a bech32 string (https://en.bitcoin.it/wiki/BIP_0173#Bech32 )
func (b Beacon) ToBech32Str() string {
	encodedBytes := b.EncodeTLV()
	res, err := encodeUtils.Bech32EncodeBytes(encodedBytes, MoneysocketHrp)
	// this is theoretically possible with enough locations
	if err != nil {
		panic(err)
	}
	return res
}

// DecodeTLV decodes a TLV (type-length-value, defined in BOLT #1: https://git.io/JLCRq)
// into a Beacon. Returns an error if Beacon cannot be decoded
func DecodeTLV(b []byte) (beacon Beacon, err error) {
	beaconTlv, _, err := util.TLVPop(b)
	if err != nil {
		return beacon, err
	}
	if beaconTlv.Type() != util.BeaconTLVType {
		return beacon, fmt.Errorf("got unexpected tlv type: %d expected %d", beaconTlv.Type(), util.BeaconTLVType)
	}

	ssTlv, remainder, err := util.TLVPop(beaconTlv.Value())
	if err != nil {
		return beacon, err
	}

	if ssTlv.Type() != util.SharedSeedTLVType {
		return beacon, fmt.Errorf("got unexpected shared seed tlv type %d, expected: %d", ssTlv.Type(), util.SharedSeedTLVType)
	}

	llTlv, remainder, err := util.TLVPop(remainder)
	if err != nil {
		return beacon, err
	}

	if llTlv.Type() != util.LocationListTLVType {
		return beacon, fmt.Errorf("got unexpected location list tlv type: %d, expected: %d", llTlv.Type(), util.LocationListTLVType)
	}

	beacon.seed, err = BytesToSharedSeed(ssTlv.Value())
	if err != nil {
		return beacon, err
	}

	lcTlv, remainder, err := util.TLVPop(llTlv.Value())
	if err != nil {
		return beacon, err
	}
	if lcTlv.Type() != util.LocationCountTLVType {
		return beacon, fmt.Errorf("got unexpected location list tlv type: %d, expected: %d", lcTlv.Type(), util.LocationCountTLVType)
	}

	locationCount, _, err := bigsize.Pop(lcTlv.Value())
	if err != nil {
		return beacon, err
	}

	// TODO break this out into it's own function to reduce cyclomatic compleixty
	var locations []location.Location
	for i := 0; i < int(locationCount); i++ {
		llTlv, remainder, err = util.TLVPop(remainder)
		if err != nil {
			return beacon, err
		}
		var loc location.Location
		switch llTlv.Type() {
		case util.WebsocketLocationTLVType:
			loc, err = location.WebsocketLocationFromTLV(llTlv)
		case util.WebRTCLocationTLVType:
			loc, err = location.WebRTCLocationFromTLV(llTlv)
		case util.BluetoothLocationTLVType:
			loc, err = location.BluetoothLocationFromTLV(llTlv)
		case util.NFCLocationTLVType:
			loc, err = location.NfcLocationFromTLV(llTlv)
		default:
			panic(fmt.Errorf("location type %d not yet implemented", llTlv.Type()))
		}
		if err != nil {
			return beacon, err
		}
		locations = append(locations, loc)
	}
	beacon.locations = locations
	return beacon, err
}

// DecodeFromBech32Str decodes a Beacon from a bech32 string (https://en.bitcoin.it/wiki/BIP_0173#Bech32 )
func DecodeFromBech32Str(bech32 string) (Beacon, error) {
	hrp, decodedBytes, err := encodeUtils.Bech32DecodeBytes(bech32)
	if err != nil {
		return Beacon{}, err
	}
	_ = decodedBytes
	if hrp != MoneysocketHrp {
		return Beacon{}, fmt.Errorf("got hrp %s when decoding tlv, expected %s", hrp, MoneysocketHrp)
	}
	return DecodeTLV(decodedBytes)
}
