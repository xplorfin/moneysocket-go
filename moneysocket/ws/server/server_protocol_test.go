package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/posener/wstest"
	. "github.com/stretchr/testify/assert"
)

type TestIncomingMessages struct {
	// payload
	payload []byte
	// wether or not message is a binary
	isBinary bool
}

type TestWebsocketServiceProtocol struct {
	BaseWebsocketService
	onConnectingCalls int
	onConnectCalls    int
	onOpenCalls       int
	onMessageCalls    int
	onCloseCalls      int
	Messages          []TestIncomingMessages
}

func (t *TestWebsocketServiceProtocol) OnConnecting(r *http.Request) {
	t.onConnectingCalls++
}

func (t *TestWebsocketServiceProtocol) OnConnect(r *http.Request) {
	t.onConnectCalls++
}

func (t *TestWebsocketServiceProtocol) OnOpen() {
	t.onOpenCalls++
}

func (t *TestWebsocketServiceProtocol) OnWsMessage(payload []byte, isBinary bool) {
	t.onMessageCalls++
	t.Messages = append(t.Messages, TestIncomingMessages{
		payload:  payload,
		isBinary: isBinary,
	})
}

func (t *TestWebsocketServiceProtocol) OnClose(wasClean bool, code int, reason string) {
	t.onCloseCalls++
}

func (t *TestWebsocketServiceProtocol) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHTTP(t, w, r)
}

func NewTestWebsocketServiceProtocol() TestWebsocketServiceProtocol {
	return TestWebsocketServiceProtocol{
		BaseWebsocketService: NewBaseWebsocketService(),
		onConnectingCalls:    0,
		onConnectCalls:       0,
		onOpenCalls:          0,
		onMessageCalls:       0,
		onCloseCalls:         0,
	}
}

func TestWebsocketEventHandlers(t *testing.T) {
	var (
		testBinaryMessage = []byte("test binary message")
		testTextMessage   = []byte("test text message")
	)
	// start the server
	testWebsocketServiceProtocol := NewTestWebsocketServiceProtocol()
	d := wstest.NewDialer(&testWebsocketServiceProtocol)

	c, _, err := d.Dial("ws://example.org/ws", nil)
	if err != nil {
		t.Error(err)
	}

	err = c.WriteMessage(websocket.BinaryMessage, testBinaryMessage)
	if err != nil {
		t.Error(err)
	}
	err = c.WriteMessage(websocket.TextMessage, testTextMessage)
	if err != nil {
		t.Error(err)
	}

	// process is async
	const timeoutInterval = 2
	timeout := time.Now().Add(timeoutInterval * time.Second)
	for {
		if len(testWebsocketServiceProtocol.Messages) == 2 {
			Equal(t, testWebsocketServiceProtocol.onConnectingCalls, 1)
			Equal(t, testWebsocketServiceProtocol.onConnectCalls, 1)
			Equal(t, testWebsocketServiceProtocol.onOpenCalls, 1)
			// verify messages
			Equal(t, testWebsocketServiceProtocol.Messages[0].isBinary, true)
			Equal(t, testWebsocketServiceProtocol.Messages[0].payload, testBinaryMessage)

			Equal(t, testWebsocketServiceProtocol.Messages[1].isBinary, false)
			Equal(t, testWebsocketServiceProtocol.Messages[1].payload, testTextMessage)
			break
		}
		if time.Now().Unix() > timeout.Unix() {
			t.Errorf("process timed out after %d seconds without receiving messages", timeoutInterval)
			break
		}
	}

	// TODO test close
}
