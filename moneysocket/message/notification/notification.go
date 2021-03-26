package notification

import (
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type MoneysocketNotification interface {
	base.MoneysocketMessage
	// get the notification uuid
	NotificationUUID() string
	RequestReferenceUUID() string
	NotificationName() string
	RequestType() base.MessageType
}

// BaseMoneySocketNotification is the notification type
type BaseMoneySocketNotification struct {
	base.MoneysocketMessage
	// BaseNotificationUUID is the uuid for this message
	BaseNotificationUUID string
	// BaseRequestReferenceUUID is the request reference id
	BaseRequestReferenceUUID string
	// requestType is the base.MessageType
	requestType base.MessageType
}

func NewBaseMoneySocketNotification(notificationType base.MessageType, requestUUID string) BaseMoneySocketNotification {
	return BaseMoneySocketNotification{
		MoneysocketMessage:       base.NewBaseBaseMoneysocketMessage(base.Notification),
		BaseNotificationUUID:     uuid.NewV4().String(),
		requestType:              notificationType,
		BaseRequestReferenceUUID: requestUUID,
	}
}

func (b BaseMoneySocketNotification) RequestType() base.MessageType {
	return b.requestType
}

func (b BaseMoneySocketNotification) MessageClass() base.MessageClass {
	return base.Notification
}

func (b BaseMoneySocketNotification) NotificationUUID() string {
	return b.BaseNotificationUUID
}

func (b BaseMoneySocketNotification) RequestReferenceUUID() string {
	return b.BaseRequestReferenceUUID
}

func (b BaseMoneySocketNotification) NotificationName() string {
	return b.RequestType().ToString()
}

const (
	UUIDKey                 = "notification_uuid"
	RequestReferenceUUIDKey = "request_reference_uuid"
	NameKey                 = "notification_name"
)

func EncodeMoneysocketNotification(msg MoneysocketNotification, toEncode map[string]interface{}) error {
	err := base.EncodeMoneysocketMessage(msg, toEncode)
	if err != nil {
		return err
	}
	toEncode[UUIDKey] = msg.NotificationUUID()
	toEncode[RequestReferenceUUIDKey] = msg.RequestReferenceUUID()
	toEncode[NameKey] = msg.NotificationName()
	return nil
}

func DecodeRequest(request []byte) (b BaseMoneySocketNotification, err error) {
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
	refUUID, err := jsonparser.GetString(request, RequestReferenceUUIDKey)
	if err != nil {
		return b, err
	}
	return BaseMoneySocketNotification{
		MoneysocketMessage:       baseMessage,
		BaseNotificationUUID:     reqUUID,
		BaseRequestReferenceUUID: refUUID,
		requestType:              base.MessageTypeFromString(reqType),
	}, nil
}
