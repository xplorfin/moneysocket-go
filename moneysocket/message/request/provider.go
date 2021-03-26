package request

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type Provider struct {
	BaseMoneySocketRequest
}

func NewRequestProvider() Provider {
	return Provider{
		NewBaseMoneySocketRequest(base.ProviderRequest),
	}
}

// encode a request provider to json
func (rp Provider) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketRequest(rp, m)
	if err != nil {
		panic(err)
	}
	encodedRequest, err := json.Marshal(m)
	return encodedRequest, err
}

// decode a request provider from json
func DecodeRequestProvider(payload []byte) (Provider, error) {
	request, err := DecodeRequest(payload)
	if err != nil {
		return Provider{}, err
	}
	return Provider{request}, err
}

var _ MoneysocketRequest = &Provider{}
