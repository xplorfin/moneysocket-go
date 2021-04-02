package base

import (
	"fmt"
	"time"
)

// MoneysocketMessage is probably one of the worst ways to handle json in go
// but I wanted to make sure the first vresion was as close to moneysocket as possible
// architecturally so I could safely refactor with test cases in place
type MoneysocketMessage interface {
	// get timestamp message was created
	Timestamp() time.Time
	// get protocol
	Protocol() string
	// get protocol version
	ProtocolVersion() string
	// get message class (request or message)
	MessageClass() MessageClass
	// convert a message to json
	ToJSON() ([]byte, error)
	// wetrher or not the message is valid
	IsValid() (bool, error)
	// wether or not a message must be clear text
	MustBeClearText() bool
	// encryption level
	CryptLevel() string
}

// json keys
const (
	timestampKey    = "timestamp"
	protocolKey     = "protocol"
	protocolVersion = "protocol_version"
	MessageClassKey = "message_class"
)

// EncodeMoneysocketMessage maps are passed by reference by default https://bit.ly/35KrDps
func EncodeMoneysocketMessage(msg MoneysocketMessage, toEncode map[string]interface{}) error {
	if toEncode == nil {
		return fmt.Errorf("map must be initialized")
	}
	// python's time format is [unix].[nano]
	toEncode[timestampKey] = float64(msg.Timestamp().UnixNano()) / float64(time.Second)
	toEncode[protocolKey] = msg.Protocol()
	toEncode[protocolVersion] = msg.ProtocolVersion()
	toEncode[MessageClassKey] = msg.MessageClass().ToString()
	return nil
}
