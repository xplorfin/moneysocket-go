package request

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type RequestProvider struct {
	BaseMoneySocketRequest
}

func NewRequestProvider() RequestProvider {
	return RequestProvider{
		NewBaseMoneySocketRequest(base.ProviderRequest),
	}
}

// encode a request provider to json
func (rp RequestProvider) ToJson() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketRequest(rp, m)
	if err != nil {
		panic(err)
	}
	encodedRequest, err := json.Marshal(m)
	return encodedRequest, err
}

// decode a request provider from json
func DecodeRequestProvider(payload []byte) (RequestProvider, error) {
	request, err := DecodeRequest(payload)
	if err != nil {
		return RequestProvider{}, err
	}
	return RequestProvider{request}, err
}

var _ MoneysocketRequest = &RequestProvider{}
