package request

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// PingRequest is a message type used for pinging
type PingRequest struct {
	BaseMoneySocketRequest
}

// NewPingRequest creates a PingRequest
func NewPingRequest() PingRequest {
	return PingRequest{
		NewBaseMoneySocketRequest(base.PingRequest),
	}
}

// ToJSON encodes a PingRequest to json
func (p PingRequest) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketRequest(p, m)
	if err != nil {
		panic(err)
	}
	encodedRequest, err := json.Marshal(m)
	return encodedRequest, err
}

// DecodePing gets a PingRequest from json
func DecodePing(payload []byte) (PingRequest, error) {
	request, err := DecodeRequest(payload)
	if err != nil {
		return PingRequest{}, err
	}
	return PingRequest{request}, err
}

var _ MoneysocketRequest = &PingRequest{}
