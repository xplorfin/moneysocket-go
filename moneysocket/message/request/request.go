package request

import (
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// MoneysocketRequest is an interface for requests
type MoneysocketRequest interface {
	base.MoneysocketMessage
	// UUID is the uuid of the of the request
	UUID() string
	// RequestName gets the request name
	RequestName() string
	// MessageType is the type of the message
	MessageType() base.MessageType
}

// BaseMoneySocketRequest is a moneysocket request
type BaseMoneySocketRequest struct {
	base.MoneysocketMessage
	// BaseUUID is the uuid of the request
	BaseUUID string
	// RequestType is the request type
	RequestType base.MessageType
}

// MessageClass is the base.MessageClass. This is always base.Request
func (b BaseMoneySocketRequest) MessageClass() base.MessageClass {
	return base.Request
}

// UUID is the uuid of the BaseMoneySocketRequest
func (b BaseMoneySocketRequest) UUID() string {
	return b.BaseUUID
}

// RequestName gets the request name from the MessageClass
func (b BaseMoneySocketRequest) RequestName() string {
	return b.MessageType().ToString()
}

// MessageType gets the message type
func (b BaseMoneySocketRequest) MessageType() base.MessageType {
	return b.RequestType
}

const (
	// UUIDKey is the key used when encoding json
	UUIDKey = "request_uuid"
	// NameKey is the key used for encoding json
	NameKey = "request_name"
)

// EncodeMoneysocketRequest creates a moneysocket request
func EncodeMoneysocketRequest(msg MoneysocketRequest, toEncode map[string]interface{}) error {
	err := base.EncodeMoneysocketMessage(msg, toEncode)
	if err != nil {
		return err
	}
	toEncode[UUIDKey] = msg.UUID()
	toEncode[NameKey] = msg.RequestName()
	return nil
}

// NewBaseMoneySocketRequest generates a new base moneysocket request. Should only be used by other message classes
func NewBaseMoneySocketRequest(requestType base.MessageType) BaseMoneySocketRequest {
	return BaseMoneySocketRequest{
		base.NewBaseBaseMoneysocketMessage(base.Request),
		uuid.NewV4().String(),
		requestType,
	}
}

// DecodeRequest decodes a base moneysocket request
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
