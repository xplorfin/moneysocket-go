package request

import (
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type MoneysocketRequest interface {
	base.MoneysocketMessage
	Uuid() string
	RequestName() string
	MessageType() base.MessageType
}

type BaseMoneySocketRequest struct {
	base.BaseMoneysocketMessage
	BaseUuid    string
	RequestType base.MessageType
}

func (b BaseMoneySocketRequest) MessageClass() base.MessageClass {
	return base.Request
}

func (b BaseMoneySocketRequest) Uuid() string {
	return b.BaseUuid
}

func (b BaseMoneySocketRequest) RequestName() string {
	return b.MessageType().ToString()
}

// get the message type
func (b BaseMoneySocketRequest) MessageType() base.MessageType {
	return b.RequestType
}

const (
	RequestUuidKey = "request_uuid"
	RequestNameKey = "request_name"
)

// create a moneysocket request
func EncodeMoneysocketRequest(msg MoneysocketRequest, toEncode map[string]interface{}) error {
	err := base.EncodeMoneysocketMessage(msg, toEncode)
	if err != nil {
		return err
	}
	toEncode[RequestUuidKey] = msg.Uuid()
	toEncode[RequestNameKey] = msg.RequestName()
	return nil
}

// generate a new base moneysocket request. Should only be used by other message classes
func NewBaseMoneySocketRequest(requestType base.MessageType) BaseMoneySocketRequest {
	return BaseMoneySocketRequest{
		base.NewBaseMoneysocketMessage(base.Request),
		uuid.NewV4().String(),
		requestType,
	}
}

func DecodeRequest(request []byte) (b BaseMoneySocketRequest, err error) {
	baseMessage, err := base.DecodeBaseMoneysocketMessage(request)
	if err != nil {
		return b, err
	}
	reqUuid, err := jsonparser.GetString(request, RequestUuidKey)
	if err != nil {
		return b, err
	}
	reqType, err := jsonparser.GetString(request, RequestNameKey)
	if err != nil {
		return b, err
	}
	return BaseMoneySocketRequest{
		BaseMoneysocketMessage: baseMessage,
		BaseUuid:               reqUuid,
		RequestType:            base.MessageTypeFromString(reqType),
	}, nil
}
