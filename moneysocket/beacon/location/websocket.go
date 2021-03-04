package location

import (
	"fmt"

	util2 "github.com/xplorfin/moneysocket-go/moneysocket/util"

	tlvHelper "github.com/lightningnetwork/lnd/tlv"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon/util/bigsize"
)

type WebsocketLocation struct {
	// wether or nopt to use tls
	tls bool
	// the host of the websocket
	host string
	// the port to use
	port uint64
}

const (
	WebsocketType = "WebSocket"
	// host tlvHelper position
	HostTlvType   = 0
	UseTlsTlvType = 1
	// bool, corollary to use_tls_enum_value https://git.io/JLgWa

	// indicate tls is used
	UseTlsEnumValueTrue = 0
	// indicate tls is not used
	UseTlsEnumFalse = 1

	// port tlvHelper position
	PortTlvType = 2
	// default wss port
	DefaultTlsPort = 443
	// default ws port
	DefaultNoTlsPort = 80
)

// statically assert that this type binds to location interface
var _ Location = WebsocketLocation{}

// new websocket location
func NewWebsocketLocation(host string, useTls bool) WebsocketLocation {
	port := DefaultTlsPort
	if !useTls {
		port = DefaultNoTlsPort
	}
	return NewWebsocketLocationPort(host, useTls, port)
}

func NewWebsocketLocationPort(host string, useTls bool, port int) WebsocketLocation {
	return WebsocketLocation{
		tls:  useTls,
		host: host,
		port: uint64(port),
	}
}

func (ws WebsocketLocation) getProtocol() string {
	if ws.IsTls() {
		return "wss"
	}
	return "ws"
}

func (ws WebsocketLocation) Type() tlvHelper.Type {
	return util.WebsocketLocationTlvType
}

func (ws WebsocketLocation) ToString() string {
	return fmt.Sprintf("%s://%s:%d", ws.getProtocol(), ws.host, ws.port)
}

func (ws WebsocketLocation) IsTls() bool {
	return ws.tls
}

func (ws WebsocketLocation) EncodedTlv() []byte {
	return util2.TlvRecordToBytes(ws.Tlv()...)
}

func (ws WebsocketLocation) Tlv() []tlvHelper.Record {
	byteHost := []uint8(ws.host)
	record := tlvHelper.MakePrimitiveRecord(HostTlvType, &byteHost)
	encoded := util2.TlvRecordToBytes(record)
	if !ws.IsTls() {
		value := uint64(0)
		encoded = append(encoded, util2.TlvRecordToBytes(tlvHelper.MakeStaticRecord(UseTlsTlvType, &value, tlvHelper.VarIntSize(uint64(UseTlsTlvType)), util.EVarInt, util.DVarInt))...)
		if ws.port != DefaultNoTlsPort {
			encoded = append(encoded, util2.TlvRecordToBytes(tlvHelper.MakeStaticRecord(PortTlvType, &ws.port, tlvHelper.VarIntSize(ws.port), util.EVarInt, util.DVarInt))...)
		}
	} else {
		if ws.port != DefaultTlsPort {
			encoded = append(encoded, util2.TlvRecordToBytes(tlvHelper.MakeStaticRecord(PortTlvType, &ws.port, tlvHelper.VarIntSize(ws.port), util.EVarInt, util.DVarInt))...)
		}
	}
	return []tlvHelper.Record{tlvHelper.MakeStaticRecord(util.WebsocketLocationTlvType, &encoded, uint64(len(encoded)), tlvHelper.EVarBytes, tlvHelper.DVarBytes)}
}

func WebsocketLocationFromTlv(tlv util.Tlv) (wsl WebsocketLocation, err error) {
	if tlv.Type() != util.WebsocketLocationTlvType {
		return wsl, fmt.Errorf("got unexpected subTlv type: %d, expected %d", tlv.Type(), util.WebsocketLocationTlvType)
	}
	wsl.tls = true
	hasHostTlv := false
	tlvs, err := util.NamespaceIterTlvs(tlv.Value())
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

func (ws WebsocketLocation) ToObject() map[string]interface{} {
	m := make(map[string]interface{})
	m["type"] = WebsocketType
	m["host"] = ws.host
	m["port"] = ws.port
	m["use_tls"] = ws.tls
	return m
}
