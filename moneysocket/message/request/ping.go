package request

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type PingRequest struct {
	BaseMoneySocketRequest
}

func NewPingRequest() PingRequest {
	return PingRequest{
		NewBaseMoneySocketRequest(base.PingRequest),
	}
}

func (p PingRequest) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketRequest(p, m)
	if err != nil {
		panic(err)
	}
	encodedRequest, err := json.Marshal(m)
	return encodedRequest, err
}

func DecodePing(payload []byte) (PingRequest, error) {
	request, err := DecodeRequest(payload)
	if err != nil {
		return PingRequest{}, err
	}
	return PingRequest{request}, err
}

var _ MoneysocketRequest = &PingRequest{}
