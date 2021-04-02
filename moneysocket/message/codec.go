package message

import (
	"encoding/json"
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/xplorfin/moneysocket-go/moneysocket/beacon"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/notification"
)

// IsCypherText checks if string is cipher (in actuality, this method checks whether a string is not json)
func IsCypherText(text []byte) bool {
	var js map[string]interface{}
	return json.Unmarshal(text, &js) != nil
}

// WireEncode encodes a message, encrypt if shared seed is present
func WireEncode(msg base.MoneysocketMessage, ss *beacon.SharedSeed) ([]byte, error) {
	encoded, err := msg.ToJSON()
	if err != nil {
		return nil, err
	}
	if msg.MustBeClearText() || ss == nil {
		return encoded, nil
	}
	// encrypt
	encryptedString, err := Encrypt(encoded, ss.DeriveAES256Key())
	return encryptedString, err
}

// WireDecode decodes a message from a moneysocket message
func WireDecode(msgBytes []byte, sharedSeed *beacon.SharedSeed) (msg base.MoneysocketMessage, msgType base.MessageType, err error) {
	isCypherText := IsCypherText(msgBytes)

	if isCypherText && sharedSeed == nil {
		return msg, msgType, fmt.Errorf("no seed to decrypt cypher text")
	}

	rawMsg := msgBytes

	if isCypherText && sharedSeed != nil {
		// nolint: todo fix, this has a high bug surface
		raw, err := Decrypt(msgBytes, (*sharedSeed).DeriveAES256Key())
		if err != nil {
			return msg, msgType, err
		}
		rawMsg = []byte(raw)
	}

	rawMsgClass, err := jsonparser.GetString(rawMsg, base.MessageClassKey)
	if err != nil {
		return msg, msgType, err
	}

	msgClass := base.MessageClassFromString(rawMsgClass)

	switch msgClass {
	case base.Request:
		return FromText(rawMsg)
	case base.Notification:
		return notification.FromText(rawMsg)
	default:
		panic(fmt.Sprintf("unhandled message type %d", msgClass))
	}
}

// LocalEncode encodes a message
func LocalEncode(msg base.MoneysocketMessage, sharedSeed *beacon.SharedSeed) (encoded bool, msgBytes []byte) {
	msgBytes, err := msg.ToJSON()
	if err != nil {
		panic(err)
	}
	if msg.MustBeClearText() {
		return false, msgBytes
	}
	if sharedSeed == nil {
		panic("shared sed can't be null")
	}
	enc, err := Encrypt(msgBytes, sharedSeed.DeriveAES256Key())
	if err != nil {
		panic(err)
	}
	return true, enc
}
