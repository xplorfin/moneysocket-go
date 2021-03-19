package util

import "github.com/lightningnetwork/lnd/tlv"

// maintains parity with beacons/__init__.py https://git.io/JLCRc
// note: if you add a new type, add it to the test as well.
// This asserts types are unique and follow the spec (other than tlv_type_start, of course)

const (
	// TLVMinimum is the minimum tlv value
	// types cannot be less than 2^16 according to bolt-1 https://git.io/JLCRq
	TLVMinimum = 65536

	// TLVTypeStart is the starting value is 2^16+443 and
	// will increment by 2 for new types https://git.io/JLCRc
	TLVTypeStart = TLVMinimum + 443

	// BeaconTLVType is the tlv type for a beacon
	BeaconTLVType tlv.Type = TLVTypeStart

	// SharedSeedTLVType is the tlv type for a shared type
	SharedSeedTLVType tlv.Type = TLVTypeStart + 2

	// LocationCountTLVType is the tlv type for a location count
	LocationCountTLVType tlv.Type = TLVTypeStart + 4

	// LocationListTLVType is the tkv type that prefixes a list of tlv locations
	LocationListTLVType tlv.Type = TLVTypeStart + 6

	// WebsocketLocationTLVType is the tlv type used for websockets
	WebsocketLocationTLVType tlv.Type = TLVTypeStart + 8

	// WebRTCLocationTLVType is the tlv type of a web rtc location
	// TODO beacons/__init.py https://git.io/JLC0J
	WebRTCLocationTLVType = TLVTypeStart + 10

	// BluetoothLocationTLVType is the bluetooth location tlv type
	// TODO beacons/__init__.py https://git.io/JLC0I
	BluetoothLocationTLVType tlv.Type = TLVTypeStart + 12

	// NFCLocationTLVType is the tlv type of an nfc loation tlv
	// TODO beacons/__init__.py https://git.io/JLC0g
	NFCLocationTLVType tlv.Type = TLVTypeStart + 14
)

// TLVTypes is a list of all custom implemented tlv types in the package
var TLVTypes = []tlv.Type{BeaconTLVType, SharedSeedTLVType, LocationCountTLVType, LocationListTLVType, WebsocketLocationTLVType, WebRTCLocationTLVType, BluetoothLocationTLVType, NFCLocationTLVType}
