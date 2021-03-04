package notification

import (
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

type MoneysocketNotification interface {
	base.MoneysocketMessage
	// get the notification uuid
	NotificationUuid() string
	RequestReferenceUuid() string
	NotificationName() string
	RequestType() base.MessageType
}

type BaseMoneySocketNotification struct {
	base.BaseMoneysocketMessage
	BaseNotificationUuid     string
	BaseRequestReferenceUuid string
	requestType              base.MessageType
}

func NewBaseMoneySocketNotification(notificationType base.MessageType, requestUuid string) BaseMoneySocketNotification {
	return BaseMoneySocketNotification{
		BaseMoneysocketMessage:   base.NewBaseMoneysocketMessage(base.Notification),
		BaseNotificationUuid:     uuid.NewV4().String(),
		requestType:              notificationType,
		BaseRequestReferenceUuid: requestUuid,
	}
}

func (b BaseMoneySocketNotification) RequestType() base.MessageType {
	return b.requestType
}

func (b BaseMoneySocketNotification) MessageClass() base.MessageClass {
	return base.Notification
}

func (b BaseMoneySocketNotification) NotificationUuid() string {
	return b.BaseNotificationUuid
}

func (b BaseMoneySocketNotification) RequestReferenceUuid() string {
	return b.BaseRequestReferenceUuid
}

func (b BaseMoneySocketNotification) NotificationName() string {
	return b.RequestType().ToString()
}

const (
	NotificationUuidKey     = "notification_uuid"
	RequestReferenceUuidKey = "request_reference_uuid"
	NotificationNameKey     = "notification_name"
)

func EncodeMoneysocketNotification(msg MoneysocketNotification, toEncode map[string]interface{}) error {
	err := base.EncodeMoneysocketMessage(msg, toEncode)
	if err != nil {
		return err
	}
	toEncode[NotificationUuidKey] = msg.NotificationUuid()
	toEncode[RequestReferenceUuidKey] = msg.RequestReferenceUuid()
	toEncode[NotificationNameKey] = msg.NotificationName()
	return nil
}

func DecodeRequest(request []byte) (b BaseMoneySocketNotification, err error) {
	baseMessage, err := base.DecodeBaseMoneysocketMessage(request)
	if err != nil {
		return b, err
	}
	reqUuid, err := jsonparser.GetString(request, NotificationUuidKey)
	if err != nil {
		return b, err
	}
	reqType, err := jsonparser.GetString(request, NotificationNameKey)
	if err != nil {
		return b, err
	}
	refUuid, err := jsonparser.GetString(request, RequestReferenceUuidKey)
	if err != nil {
		return b, err
	}
	return BaseMoneySocketNotification{
		BaseMoneysocketMessage:   baseMessage,
		BaseNotificationUuid:     reqUuid,
		BaseRequestReferenceUuid: refUuid,
		requestType:              base.MessageTypeFromString(reqType),
	}, nil
}
