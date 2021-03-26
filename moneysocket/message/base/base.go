package base

import (
	"time"

	"github.com/buger/jsonparser"
)

const (
	ProtocolName = "Moneysocket"
	Version      = "0.0.0"
)

// MessageBase is the message
type MessageBase struct {
	// Time is the time
	Time time.Time
	// BaseProtocol is the base protocol name
	BaseProtocol string
	// BaseProtocolVersion is the protocol version
	BaseProtocolVersion string
	// BaseMessageClass is the message class
	BaseMessageClass MessageClass
}

func (m MessageBase) CryptLevel() string {
	return "AES"
}

func (m MessageBase) ProtocolVersion() string {
	return m.BaseProtocolVersion
}

func (m MessageBase) MessageClass() MessageClass {
	return m.BaseMessageClass
}

func (m MessageBase) MustBeClearText() bool {
	return false
}

func (m MessageBase) ToJSON() ([]byte, error) {
	panic("must be implemented in children classes. You can use EncodeBaseMoneysocketMessage as a helper method")
}

func (m MessageBase) IsValid() (bool, error) {
	panic("must be implemented in children classes")
}

func NewBaseBaseMoneysocketMessage(messageType MessageClass) MessageBase {
	return MessageBase{
		Time:                time.Now(),
		BaseProtocol:        ProtocolName,
		BaseProtocolVersion: Version,
		BaseMessageClass:    messageType,
	}
}

func (m MessageBase) Timestamp() time.Time {
	return m.Time
}

func (m MessageBase) Protocol() string {
	return m.BaseProtocol
}

// decode a moneysocket message from json
func DecodeBaseBaseMoneysocketMessage(payload []byte) (b MessageBase, err error) {
	// TODO get float
	parsedTime, err := jsonparser.GetFloat(payload, timestampKey)
	if err != nil {
		return MessageBase{}, err
	}

	o := int64(float64(time.Millisecond) * (parsedTime - float64(int64(parsedTime))))
	_ = o
	b.Time = time.Unix(int64(parsedTime), int64(float64(time.Second)*(parsedTime-float64(int64(parsedTime)))))
	msgClass, err := jsonparser.GetString(payload, MessageClassKey)
	if err != nil {
		return MessageBase{}, err
	}
	b.BaseMessageClass = MessageClassFromString(msgClass)
	b.BaseProtocolVersion, err = jsonparser.GetString(payload, protocolVersion)
	if err != nil {
		return MessageBase{}, err
	}
	b.BaseProtocol, err = jsonparser.GetString(payload, protocolKey)
	if err != nil {
		return MessageBase{}, err
	}
	return b, err

}

var _ MoneysocketMessage = &MessageBase{}
