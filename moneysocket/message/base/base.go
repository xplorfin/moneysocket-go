package base

import (
	"time"

	"github.com/buger/jsonparser"
)

const (
	PROTOCOL_NAME = "Moneysocket"
	VERSION       = "0.0.0"
)

type BaseMoneysocketMessage struct {
	Time                time.Time
	BaseProtocol        string
	BaseProtocolVersion string
	BaseMessageClass    MessageClass
}

func (m BaseMoneysocketMessage) CryptLevel() string {
	return "AES"
}

func (m BaseMoneysocketMessage) ProtocolVersion() string {
	return m.BaseProtocolVersion
}

func (m BaseMoneysocketMessage) MessageClass() MessageClass {
	return m.BaseMessageClass
}

func (m BaseMoneysocketMessage) MustBeClearText() bool {
	return false
}

func (m BaseMoneysocketMessage) ToJson() ([]byte, error) {
	panic("must be implemented in children classes. You can use EncodeMoneysocketMessage as a helper method")
}

func (m BaseMoneysocketMessage) IsValid() (bool, error) {
	panic("must be implemented in children classes")
}

func NewBaseMoneysocketMessage(messageType MessageClass) BaseMoneysocketMessage {
	return BaseMoneysocketMessage{
		Time:                time.Now(),
		BaseProtocol:        PROTOCOL_NAME,
		BaseProtocolVersion: VERSION,
		BaseMessageClass:    messageType,
	}
}

func (m BaseMoneysocketMessage) Timestamp() time.Time {
	return m.Time
}

func (m BaseMoneysocketMessage) Protocol() string {
	return m.BaseProtocol
}

// decode a moneysocket message from json
func DecodeBaseMoneysocketMessage(payload []byte) (b BaseMoneysocketMessage, err error) {
	// TODO get float
	parsedTime, err := jsonparser.GetFloat(payload, timestampKey)
	if err != nil {
		return BaseMoneysocketMessage{}, err
	}

	o := int64(float64(time.Millisecond) * (parsedTime - float64(int64(parsedTime))))
	_ = o
	b.Time = time.Unix(int64(parsedTime), int64(float64(time.Second)*(parsedTime-float64(int64(parsedTime)))))
	msgClass, err := jsonparser.GetString(payload, MessageClassKey)
	if err != nil {
		return BaseMoneysocketMessage{}, err
	}
	b.BaseMessageClass = MessageClassFromString(msgClass)
	b.BaseProtocolVersion, err = jsonparser.GetString(payload, protocolVersion)
	if err != nil {
		return BaseMoneysocketMessage{}, err
	}
	b.BaseProtocol, err = jsonparser.GetString(payload, protocolKey)
	if err != nil {
		return BaseMoneysocketMessage{}, err
	}
	return b, err

}

var _ MoneysocketMessage = &BaseMoneysocketMessage{}
