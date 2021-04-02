package request

import (
	"encoding/json"

	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// Provider is a provider message.
type Provider struct {
	BaseMoneySocketRequest
}

// NewRequestProvider creates a Provider for messages of base.ProviderRequest type.
func NewRequestProvider() Provider {
	return Provider{
		NewBaseMoneySocketRequest(base.ProviderRequest),
	}
}

// ToJSON encodes a request provider to json.
func (rp Provider) ToJSON() ([]byte, error) {
	m := make(map[string]interface{})
	err := EncodeMoneysocketRequest(rp, m)
	if err != nil {
		panic(err)
	}
	encodedRequest, err := json.Marshal(m)
	return encodedRequest, err
}

// DecodeRequestProvider decode a request provider from json.
func DecodeRequestProvider(payload []byte) (Provider, error) {
	request, err := DecodeRequest(payload)
	if err != nil {
		return Provider{}, err
	}
	return Provider{request}, err
}

var _ MoneysocketRequest = &Provider{}
