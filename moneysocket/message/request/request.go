package request

import (
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type MoneysocketRequest interface {
	base.MoneysocketMessage
	UUID() string
	RequestName() string
	MessageType() base.MessageType
}

type BaseMoneySocketRequest struct {
	base.MoneysocketMessage
	BaseUUID    string
	RequestType base.MessageType
}

func (b BaseMoneySocketRequest) MessageClass() base.MessageClass {
	return base.Request
}

func (b BaseMoneySocketRequest) UUID() string {
	return b.BaseUUID
}

func (b BaseMoneySocketRequest) RequestName() string {
	return b.MessageType().ToString()
}

// get the message type
func (b BaseMoneySocketRequest) MessageType() base.MessageType {
	return b.RequestType
}

const (
	UUIDKey = "request_uuid"
	NameKey = "request_name"
)

// create a moneysocket request
func EncodeMoneysocketRequest(msg MoneysocketRequest, toEncode map[string]interface{}) error {
	err := base.EncodeMoneysocketMessage(msg, toEncode)
	if err != nil {
		return err
	}
	toEncode[UUIDKey] = msg.UUID()
	toEncode[NameKey] = msg.RequestName()
	return nil
}

// generate a new base moneysocket request. Should only be used by other message classes
func NewBaseMoneySocketRequest(requestType base.MessageType) BaseMoneySocketRequest {
	return BaseMoneySocketRequest{
		base.NewBaseBaseMoneysocketMessage(base.Request),
		uuid.NewV4().String(),
		requestType,
	}
}

func DecodeRequest(request []byte) (b BaseMoneySocketRequest, err error) {
	baseMessage, err := base.DecodeBaseBaseMoneysocketMessage(request)
	if err != nil {
		return b, err
	}
	reqUUID, err := jsonparser.GetString(request, UUIDKey)
	if err != nil {
		return b, err
	}
	reqType, err := jsonparser.GetString(request, NameKey)
	if err != nil {
		return b, err
	}
	return BaseMoneySocketRequest{
		MoneysocketMessage: baseMessage,
		BaseUUID:           reqUUID,
		RequestType:        base.MessageTypeFromString(reqType),
	}, nil
}
