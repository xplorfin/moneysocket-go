package notification

import (
	"github.com/buger/jsonparser"
	uuid "github.com/satori/go.uuid"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
)

// MoneysocketNotification is a notification for a message.
type MoneysocketNotification interface {
	base.MoneysocketMessage
	// NotificationUUID gets the notification uuid
	NotificationUUID() string
	// RequestReferenceUUID gets the request reference uuid
	RequestReferenceUUID() string
	// NotificationName gets the notification name
	NotificationName() string
	// RequestType gets the request type
	RequestType() base.MessageType
}

// BaseMoneySocketNotification is the notification type.
type BaseMoneySocketNotification struct {
	base.MoneysocketMessage
	// BaseNotificationUUID is the uuid for this message
	BaseNotificationUUID string
	// BaseRequestReferenceUUID is the request reference id
	BaseRequestReferenceUUID string
	// requestType is the base.MessageType
	requestType base.MessageType
}

// NewBaseMoneySocketNotification creates a.
func NewBaseMoneySocketNotification(notificationType base.MessageType, requestUUID string) BaseMoneySocketNotification {
	return BaseMoneySocketNotification{
		MoneysocketMessage:       base.NewBaseBaseMoneysocketMessage(base.Notification),
		BaseNotificationUUID:     uuid.NewV4().String(),
		requestType:              notificationType,
		BaseRequestReferenceUUID: requestUUID,
	}
}

// RequestType is the base.MessageType of the Notification.
func (b BaseMoneySocketNotification) RequestType() base.MessageType {
	return b.requestType
}

// MessageClass is the base.MessageClass of the notification. This is always notification.
func (b BaseMoneySocketNotification) MessageClass() base.MessageClass {
	return base.Notification
}

// NotificationUUID returns the uuid of the notification.
func (b BaseMoneySocketNotification) NotificationUUID() string {
	return b.BaseNotificationUUID
}

// RequestReferenceUUID gets the uuid of the request.
func (b BaseMoneySocketNotification) RequestReferenceUUID() string {
	return b.BaseRequestReferenceUUID
}

// NotificationName gets the name of the notification (from the BaseMoneySocketNotification.RequestType).
func (b BaseMoneySocketNotification) NotificationName() string {
	return b.RequestType().ToString()
}

const (
	// UUIDKey is the notification uuid in json.
	UUIDKey = "notification_uuid"
	// RequestReferenceUUIDKey is the json key for encoding the notification messages.
	RequestReferenceUUIDKey = "request_reference_uuid"
	// NameKey is the notification name key.
	NameKey = "notification_name"
)

// EncodeMoneySocketNotification encodes a MoneysocketNotification to json.
// This is used by sub-structs and should not be called directly.
func EncodeMoneySocketNotification(msg MoneysocketNotification, toEncode map[string]interface{}) error {
	err := base.EncodeMoneysocketMessage(msg, toEncode)
	if err != nil {
		return err
	}
	toEncode[UUIDKey] = msg.NotificationUUID()
	toEncode[RequestReferenceUUIDKey] = msg.RequestReferenceUUID()
	toEncode[NameKey] = msg.NotificationName()
	return nil
}

// DecodeRequest decodes a BaseMoneySocketNotification from json.
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
