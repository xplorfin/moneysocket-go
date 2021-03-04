package ws_server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

const (
	SecurePrefix   = "wss://"
	UnsecurePrefix = "ws://"
)

// listener for websocket
type WebsocketListener struct{}

// TODO
type TlsInfo struct{}

// helper for serving http and http servers
func Listen(rawWsUrl string, tlsInfo *TlsInfo, handler http.HandlerFunc) error {
	wsUrl, err := url.Parse(rawWsUrl)
	if err != nil {
		return err
	}
	if tlsInfo != nil && wsUrl.Scheme == SecurePrefix {
		return fmt.Errorf("must specify tlsInfo to listen with TLS, change the '%s' prefix to '%s' in '%s' or pass in a tls config",
			SecurePrefix,
			UnsecurePrefix,
			wsUrl)
	}

	if tlsInfo != nil {
		// TODO
	} else {
		// start insecurely
		log.Println(fmt.Sprintf("starting without TLS on %s", wsUrl.Host))
		err = http.ListenAndServe(wsUrl.Host, handler)
		if err != nil {
			return err
		}
	}
	return nil
}

func ServeHTTP(p WebSocketServerProtocol, w http.ResponseWriter, r *http.Request) {
	g, _ := errgroup.WithContext(p.Context())

	upgrader := p.Upgrader()
	if !websocket.IsWebSocketUpgrade(r) {
		WebsocketClientRoute(UnsecurePrefix+r.Host+"/", w, r)
		return
	}

	p.OnConnecting(r)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade: ", err)
		return
	}

	client := &Client{
		hub:      p.Hub(),
		conn:     conn,
		send:     make(chan []byte, 256),
		protocol: &p,
	}

	client.hub.register <- client

	p.OnConnect(r)
	p.OnOpen()

	g.Go(func() error {
		return client.writePump()
	})
	g.Go(func() error {
		return client.readPump()
	})
	_ = g.Wait()
}
