package base

import (
	"bytes"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	moneysocket_message "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/request"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
)

func EmptyMessageHandler(belowNexus nexus.Nexus, msg moneysocket_message.MoneysocketMessage) {
	// do nothing
}
func EmptyMessageBinHandler(belowNexus nexus.Nexus, msg []byte) {
	// do nothing
}

// nexus message handler
type MessageHandler func(belowNexus nexus.Nexus, msg moneysocket_message.MoneysocketMessage)
type BinMessageHandler func(belowNexus nexus.Nexus, msg []byte)

type BaseNexusTestHarness struct {
	*BaseNexus
	OnMsg    MessageHandler
	OnBinMsg BinMessageHandler
}

var _ nexus.Nexus = BaseNexusTestHarness{}

const BaseNexusTestHarnessName = "BaseNexusTestHarness"

func (b BaseNexusTestHarness) OnMessage(nexus nexus.Nexus, msg moneysocket_message.MoneysocketMessage) {
	b.OnMsg(nexus, msg)
}

func (b BaseNexusTestHarness) OnBinMessage(nexus nexus.Nexus, msg []byte) {
	b.OnBinMsg(nexus, msg)
}

func NewBaseNexusTestHarness(onMsg MessageHandler, onMsgBin BinMessageHandler) BaseNexusTestHarness {
	testHarness := BaseNexusTestHarness{}
	testHarness.BaseNexus = NewBaseNexus(BaseNexusTestHarnessName)
	testHarness.OnMsg = onMsg
	testHarness.OnBinMsg = onMsgBin
	return testHarness
}

func TestBaseNexusMsgUuidOperations(t *testing.T) {
	var n1 nexus.Nexus = NewBaseNexusTestHarness(EmptyMessageHandler, EmptyMessageBinHandler)
	var n2 nexus.Nexus = NewBaseNexusTestHarness(EmptyMessageHandler, EmptyMessageBinHandler)

	// make sure uuid operatons work as expected
	AssertUuidV4(n1, t)
	AssertUuidV4(n2, t)

	if n1.IsEqual(n2) {
		t.Errorf("expected nexus n1 with uuid %s to be different from n2", n1.Uuid())
	}

}

func TestBaseNexus(t *testing.T) {
	// generate a test message
	testMsg := []byte(gofakeit.AchAccount())
	// initialize message handler hits to false
	msgHandlerHit := false
	binMsgHandlerHit := false

	// generate a message checkers for both message
	// and bin message handlers
	msgHandler := func(belowNexus nexus.Nexus, msg moneysocket_message.MoneysocketMessage) {
		msgHandlerHit = true
	}
	binMsgHandler := func(belowNexus nexus.Nexus, msg []byte) {
		binMsgHandlerHit = true
		if !bytes.Equal(testMsg, msg) {
			t.Errorf("expected bin msg %b to equal test message %b", msg, testMsg)
		}
	}
	_ = binMsgHandler
	_ = msgHandler

	// initialize base nexuses
	var n1 nexus.Nexus = NewBaseNexusTestHarness(msgHandler, binMsgHandler)
	var n2 nexus.Nexus = NewBaseNexusTestHarness(msgHandler, binMsgHandler)

	// check message handlers
	n1.OnMessage(n2, request.NewPingRequest())
	n2.OnBinMessage(n1, testMsg)

	// make sure msgHandler function is hit
	if msgHandlerHit == false {
		t.Error("expected message msgHandler to be hit")
	}
	if binMsgHandlerHit == false {
		t.Error("expected message binMsgHandler to be hit")
	}

}

func AssertUuidV4(nexus nexus.Nexus, t *testing.T) {
	if nexus.Uuid().Version() != uuid.V4 {
		t.Errorf("expected uuid version %b to equal %b", nexus.Uuid().Version(), uuid.V4)
	}
}
