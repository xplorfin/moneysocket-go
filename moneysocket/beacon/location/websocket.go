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
	// UseTLSTLVType position when encoding the TLV
	// bool, corollary to use_tls_enum_value https://git.io/JLgWa
	UseTLSTLVType = 1

	// UseTLSEnumValueTrue indicates tls is used
	UseTLSEnumValueTrue = 0
	// UseTLSEnumFalse indicates tls is not used
	UseTLSEnumFalse = 1

	// PortTLVType is the position of the port in the encoded tlv
	PortTLVType = 2
	// DefaultTLSPort defined the default wss (secure websocket, see: https://tools.ietf.org/html/rfc6455#page-55 ) port
	DefaultTLSPort = 443
	// DefaultNoTLSPort is the default ws port (unsecure websocket, see; https://tools.ietf.org/html/rfc6455#page-54 )
	DefaultNoTLSPort = 80
)

// statically assert that this type binds to location interface
var _ Location = WebsocketLocation{}

// NewWebsocketLocation generates a WebsocketLocation from a host/tls param
func NewWebsocketLocation(host string, useTLS bool) WebsocketLocation {
	port := DefaultTLSPort
	if !useTLS {
		port = DefaultNoTLSPort
	}
	return NewWebsocketLocationPort(host, useTLS, port)
}

// NewWebsocketLocationPort generates a WebsocketLocation from a host/tls/port params
func NewWebsocketLocationPort(host string, useTLS bool, port int) WebsocketLocation {
	return WebsocketLocation{
		tls:  useTLS,
		host: host,
		port: uint64(port),
	}
}

// getProtocol gets the protocol string (see: https://tools.ietf.org/html/rfc6455#page-53 )
func (ws WebsocketLocation) getProtocol() string {
	if ws.IsTLS() {
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

// IsTLS determines wether or not the WebsocketLocation should use wss:// or ws://
func (ws WebsocketLocation) IsTLS() bool {
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
	if !ws.IsTLS() {
		value := uint64(UseTLSEnumValueTrue)
		encoded = append(encoded, moneysocketUtil.TLVRecordToBytes(tlvHelper.MakeStaticRecord(UseTLSTLVType, &value, tlvHelper.VarIntSize(uint64(UseTLSTLVType)), beaconUtil.EVarInt, beaconUtil.DVarInt))...)
		if ws.port != DefaultNoTLSPort {
			encoded = append(encoded, moneysocketUtil.TLVRecordToBytes(tlvHelper.MakeStaticRecord(PortTLVType, &ws.port, tlvHelper.VarIntSize(ws.port), beaconUtil.EVarInt, beaconUtil.DVarInt))...)
		}
	} else {
		if ws.port != DefaultTLSPort {
			encoded = append(encoded, moneysocketUtil.TLVRecordToBytes(tlvHelper.MakeStaticRecord(PortTLVType, &ws.port, tlvHelper.VarIntSize(ws.port), beaconUtil.EVarInt, beaconUtil.DVarInt))...)
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
		case UseTLSTLVType:
			wsl.tls = false
		case PortTLVType:
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
			wsl.port = DefaultTLSPort
		} else {
			wsl.port = DefaultNoTLSPort
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
