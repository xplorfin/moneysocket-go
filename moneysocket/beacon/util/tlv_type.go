// in a seperate package because this is used by beacon/locations and beacon/root
// this should probably be moved to internal if moneysocket (socket) is ever moved to its
// own package
package util

import "github.com/lightningnetwork/lnd/tlv"

// maintains parity with beacons/__init__.py https://git.io/JLCRc
// note: if you add a new type, add it to the test as well.
// This asserts types are unique and follow the spec (other than tlv_type_start, of course)

const (
	// types cannot be less than 2^16 according to bolt-1 https://git.io/JLCRq
	TlvMinimum = 65536
	// starting value is 2^16+443 and will increment by 2 for new types https://git.io/JLCRc
	TlvTypeStart = TlvMinimum + 443
	// tlv type for a beacon
	BeaconTlvType tlv.Type = TlvTypeStart
	// tlv type for a shared type
	SharedSeedTlvType tlv.Type = TlvTypeStart + 2
	// tlv type for a location count
	LocationCountTlvType tlv.Type = TlvTypeStart + 4
	// tlv type for a list
	LocationListTlvType tlv.Type = TlvTypeStart + 6
	// location websockets
	WebsocketLocationTlvType tlv.Type = TlvTypeStart + 8
	// TODO beacons/__init.py https://git.io/JLC0J
	WebrtcLocationTlvType = TlvTypeStart + 10
	// TODO beacons/__init__.py https://git.io/JLC0I
	BluetoothLocationTlvType tlv.Type = TlvTypeStart + 12
	// TODO beacons/__init__.py https://git.io/JLC0g
	NfcLocationTlvType tlv.Type = TlvTypeStart + 14
)

// list of all custom implemented tlv types in the package
var TlvTypes = []tlv.Type{BeaconTlvType, SharedSeedTlvType, LocationCountTlvType, LocationListTlvType, WebsocketLocationTlvType, WebrtcLocationTlvType, BluetoothLocationTlvType, NfcLocationTlvType}
