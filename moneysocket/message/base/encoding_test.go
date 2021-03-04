package base

import (
	"encoding/json"
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	testBaseMessage := NewBaseMoneysocketMessage(Notification)
	testMap := make(map[string]interface{})
	err := EncodeMoneysocketMessage(testBaseMessage, testMap)
	if err != nil {
		t.Error(err)
	}

	// since the timestamp in the map is a decimal, will always be greater or equal
	// nil would show up as zero
	GreaterOrEqual(t, testMap["timestamp"], float64(testBaseMessage.Time.Unix()))
	Equal(t, testMap["protocol"], testBaseMessage.Protocol())
	Equal(t, testMap["protocol_version"], testBaseMessage.ProtocolVersion())
	Equal(t, testMap["message_class"], testBaseMessage.MessageClass().ToString())

	res, err := json.Marshal(testMap)
	if err != nil {
		t.Error(err)
	}

	b, err := DecodeBaseMoneysocketMessage(res)
	if err != nil {
		t.Error(err)
	}

	GreaterOrEqual(t, b.Time.Unix(), testBaseMessage.Time.Unix())
	Equal(t, b.Protocol(), testBaseMessage.Protocol())
	Equal(t, b.ProtocolVersion(), testBaseMessage.ProtocolVersion())
	Equal(t, b.MessageClass(), testBaseMessage.MessageClass())
}
