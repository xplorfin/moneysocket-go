package location

import (
	"fmt"

	moneysocketUtil "github.com/xplorfin/moneysocket-go/moneysocket/util"

	tlvHelper "github.com/lightningnetwork/lnd/tlv"
	beaconUtil "github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util/bigsize"
)

// WebsocketLocation is a Location type for websocket beacons
type WebsocketLocation struct {
	// tls defined wether or nopt to use tls in connections
	tls bool
	// host of the websocket
	host string
	// port to use
	port uint64
}

const (
	// WebsocketType is the type of WebsocketLocation
	WebsocketType = "WebSocket"
	// HostTlvType position when encoding the TLV
	HostTlvType = 0
	// UseTlvType position when encoding the TLV
	// bool, corollary to use_tls_enum_value https://git.io/JLgWa
	UseTlsTlvType = 1

	// UseTlsEnumValueTrue indicates tls is used
	UseTlsEnumValueTrue = 0
	// UseTlsEnumFalse indicates tls is not used
	UseTlsEnumFalse = 1

	// PortTlvType is the position of the port in the encoded tlv
	PortTlvType = 2
	// DefaultTlsPort defined the default wss (secure websocket, see: https://tools.ietf.org/html/rfc6455#page-55 ) port
	DefaultTlsPort = 443
	// Default ws port (unsecure websocket, see; https://tools.ietf.org/html/rfc6455#page-54 )
	DefaultNoTlsPort = 80
)

// statically assert that this type binds to location interface
var _ Location = WebsocketLocation{}

// NewWebsocketLocation generates a WebsocketLocation from a host/tls param
func NewWebsocketLocation(host string, useTls bool) WebsocketLocation {
	port := DefaultTlsPort
	if !useTls {
		port = DefaultNoTlsPort
	}
	return NewWebsocketLocationPort(host, useTls, port)
}

// NewWebsocketLocationPort generates a WebsocketLocation from a host/tls/port params
func NewWebsocketLocationPort(host string, useTls bool, port int) WebsocketLocation {
	return WebsocketLocation{
		tls:  useTls,
		host: host,
		port: uint64(port),
	}
}

// getProtocol gets the protocol string (see: https://tools.ietf.org/html/rfc6455#page-53 )
func (ws WebsocketLocation) getProtocol() string {
	if ws.IsTls() {
		return "wss"
	}
	return "ws"
}

// Type is the tlv type of a WebsocketLocation
func (ws WebsocketLocation) Type() tlvHelper.Type {
	return beaconUtil.WebsocketLocationTLVType
}

// ToString converts a WebsocketLocation to an address to connect ot
func (ws WebsocketLocation) ToString() string {
	return fmt.Sprintf("%s://%s:%d", ws.getProtocol(), ws.host, ws.port)
}

// IsTls determines wether or not the WebsocketLocation should use wss:// or ws://
func (ws WebsocketLocation) IsTls() bool {
	return ws.tls
}

// EncodedTLV encodes a TLV for a WebsocketLocation
func (ws WebsocketLocation) EncodedTLV() []byte {
	return moneysocketUtil.TLVRecordToBytes(ws.TLV()...)
}

// TLV gets the tlv for the WebsocketLocation
func (ws WebsocketLocation) TLV() []tlvHelper.Record {
	byteHost := []uint8(ws.host)
	record := tlvHelper.MakePrimitiveRecord(HostTlvType, &byteHost)
	encoded := moneysocketUtil.TLVRecordToBytes(record)
	if !ws.IsTls() {
		value := uint64(UseTlsEnumValueTrue)
		encoded = append(encoded, moneysocketUtil.TLVRecordToBytes(tlvHelper.MakeStaticRecord(UseTlsTlvType, &value, tlvHelper.VarIntSize(uint64(UseTlsTlvType)), beaconUtil.EVarInt, beaconUtil.DVarInt))...)
		if ws.port != DefaultNoTlsPort {
			encoded = append(encoded, moneysocketUtil.TLVRecordToBytes(tlvHelper.MakeStaticRecord(PortTlvType, &ws.port, tlvHelper.VarIntSize(ws.port), beaconUtil.EVarInt, beaconUtil.DVarInt))...)
		}
	} else {
		if ws.port != DefaultTlsPort {
			encoded = append(encoded, moneysocketUtil.TLVRecordToBytes(tlvHelper.MakeStaticRecord(PortTlvType, &ws.port, tlvHelper.VarIntSize(ws.port), beaconUtil.EVarInt, beaconUtil.DVarInt))...)
		}
	}
	return []tlvHelper.Record{tlvHelper.MakeStaticRecord(beaconUtil.WebsocketLocationTLVType, &encoded, uint64(len(encoded)), tlvHelper.EVarBytes, tlvHelper.DVarBytes)}
}

// WebsocketLocationFromTLV gets the WebsocketLocation from a tlv
// returns error if this cannot be done
func WebsocketLocationFromTLV(tlv beaconUtil.TLV) (wsl WebsocketLocation, err error) {
	if tlv.Type() != beaconUtil.WebsocketLocationTLVType {
		return wsl, fmt.Errorf("got unexpected subTlv type: %d, expected %d", tlv.Type(), beaconUtil.WebsocketLocationTLVType)
	}
	wsl.tls = true
	hasHostTlv := false
	tlvs, err := beaconUtil.NamespaceIterTLVs(tlv.Value())
	for _, subTlv := range tlvs {
		switch subTlv.Type() {
		case HostTlvType:
			wsl.host = string(subTlv.Value())
			hasHostTlv = true
		case UseTlsTlvType:
			wsl.tls = false
		case PortTlvType:
			wsl.port, _, err = bigsize.Pop(subTlv.Value())
			if err != nil {
				return wsl, fmt.Errorf("could not convert port from string (port is %s)", string(subTlv.Value()))
			}
		}
	}
	if !hasHostTlv {
		return wsl, fmt.Errorf("expected host subTlv %d, got none", HostTlvType)
	}

	// if port is not set, set to default based on tls value
	if wsl.port == 0 {
		if wsl.tls {
			wsl.port = DefaultTlsPort
		} else {
			wsl.port = DefaultNoTlsPort
		}
	}

	return wsl, err
}

// ToObject converts WebsocketLocation to a json-serializable map
func (ws WebsocketLocation) ToObject() map[string]interface{} {
	m := make(map[string]interface{})
	m["type"] = WebsocketType
	m["host"] = ws.host
	m["port"] = ws.port
	m["use_tls"] = ws.tls
	return m
}
