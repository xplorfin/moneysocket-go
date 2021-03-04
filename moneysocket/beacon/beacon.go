package beacon

import (
	"fmt"

	"github.com/lightningnetwork/lnd/tlv"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/location"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util/bigsize"
	encodeUtils "github.com/xplorfin/moneysocket-go/moneysocket/util"
)

const MoneysocketHrp = "moneysocket"

type Beacon struct {
	seed      SharedSeed
	locations []location.Location
}

func NewBeacon() Beacon {
	return NewBeaconFromSeed(NewSharedSeed())
}

func NewBeaconFromSeed(seed SharedSeed) Beacon {
	return Beacon{
		seed: seed,
	}
}

func (b Beacon) Locations() []location.Location {
	return b.locations
}

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

// apppend a location to the beacon
func (b *Beacon) AddLocation(loc location.Location) {
	b.locations = append(b.locations, loc)
}

// fetch the shared seed from the beacon
func (b Beacon) GetSharedSeed() SharedSeed {
	return b.seed
}

// encode a list of locations as tlvs
func (b Beacon) EncodeLocationListTlvs() []byte {
	encoded := make([]byte, 0)
	locationCount := uint64(len(b.locations))
	locationCountEncoded := encodeUtils.TlvRecordToBytes(tlv.MakeStaticRecord(util.LocationCountTlvType, &locationCount, tlv.VarIntSize(locationCount), util.EVarInt, util.DVarInt))
	encoded = append(encoded, locationCountEncoded...)
	for _, loc := range b.locations {
		encoded = append(encoded, loc.EncodedTlv()...)
	}

	return encodeUtils.TlvRecordToBytes(tlv.MakePrimitiveRecord(util.LocationListTlvType, &encoded))
}

// encode beacon tlvs
func (b Beacon) EncodeTlvs() []byte {
	// encode the shared seed
	ssEncoded := b.seed.EncodedTLV()
	llEncoded := b.EncodeLocationListTlvs()
	combinedTlvs := append(ssEncoded, llEncoded...)
	return encodeUtils.TlvRecordToBytes(tlv.MakePrimitiveRecord(util.BeaconTlvType, &combinedTlvs))
}

func (b Beacon) ToBech32Str() string {
	encodedBytes := b.EncodeTlvs()
	res, err := encodeUtils.Bech32EncodeBytes(encodedBytes, MoneysocketHrp)
	// this is theoretically possible with enough locations
	if err != nil {
		panic(err)
	}
	return res
}

func DecodeTlvs(b []byte) (beacon Beacon, err error) {
	beaconTlv, _, err := util.TlvPop(b)
	if err != nil {
		return beacon, err
	}
	if beaconTlv.Type() != util.BeaconTlvType {
		return beacon, fmt.Errorf("got unexpected tlv type: %d expected %d", beaconTlv.Type(), util.BeaconTlvType)
	}

	ssTlv, remainder, err := util.TlvPop(beaconTlv.Value())
	if err != nil {
		return beacon, err
	}

	if ssTlv.Type() != util.SharedSeedTlvType {
		return beacon, fmt.Errorf("got unexpected shared seed tlv type %d, expected: %d", ssTlv.Type(), util.SharedSeedTlvType)
	}

	llTlv, remainder, err := util.TlvPop(remainder)
	if err != nil {
		return beacon, err
	}

	if llTlv.Type() != util.LocationListTlvType {
		return beacon, fmt.Errorf("got unexpected location list tlv type: %d, expected: %d", llTlv.Type(), util.LocationListTlvType)
	}

	beacon.seed, err = BytesToSharedSeed(ssTlv.Value())
	if err != nil {
		return beacon, err
	}

	lcTlv, remainder, err := util.TlvPop(llTlv.Value())
	if err != nil {
		return beacon, err
	}
	if lcTlv.Type() != util.LocationCountTlvType {
		return beacon, fmt.Errorf("got unexpected location list tlv type: %d, expected: %d", lcTlv.Type(), util.LocationCountTlvType)
	}

	locationCount, _, err := bigsize.Pop(lcTlv.Value())
	if err != nil {
		return beacon, err
	}

	var locations []location.Location
	for i := 0; i < int(locationCount); i++ {
		llTlv, remainder, err = util.TlvPop(remainder)
		if err != nil {
			return beacon, err
		}
		var loc location.Location
		switch llTlv.Type() {
		case util.WebsocketLocationTlvType:
			loc, err = location.WebsocketLocationFromTlv(llTlv)
		case util.WebrtcLocationTlvType:
			loc, err = location.WebRtcLocationFromTlv(llTlv)
		case util.BluetoothLocationTlvType:
			loc, err = location.BluetoothLocationFromTlv(llTlv)
		case util.NfcLocationTlvType:
			loc, err = location.NfcLocationFromTlv(llTlv)
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

func DecodeFromBech32Str(tlvBytes string) (Beacon, error) {
	hrp, decodedBytes, err := encodeUtils.Bech32DecodeBytes(tlvBytes)
	if err != nil {
		return Beacon{}, err
	}
	_ = decodedBytes
	if hrp != MoneysocketHrp {
		return Beacon{}, fmt.Errorf("got hrp %s when decoding tlv, expected %s", hrp, MoneysocketHrp)
	}
	return DecodeTlvs(decodedBytes)
}
