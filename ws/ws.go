package ws

import (
	"crypto/tls"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	TLSDialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}
)

type FastWsClient struct {
	tag            int
	conn           *websocket.Conn
	onMessage      func([]byte)
	onConnected    func()
	onDisConnected func()
}

func (fc *FastWsClient) Create(scheme string, host, path string) error {
	u := url.URL{Scheme: scheme, Host: host, Path: path}
	var er error
	if scheme == "wss" {
		fc.conn, _, er = TLSDialer.Dial(u.String(), nil)
	} else if scheme == "ws" {
		fc.conn, _, er = websocket.DefaultDialer.Dial(u.String(), nil)
	} else {
		return errors.New("未知Scheme")
	}

	if er != nil {
		return er
	}

	fc.onConnected()

	return nil
}

func (fc *FastWsClient) HandleConnected(f func()) {
	fc.onConnected = f
}

func (fc *FastWsClient) HandleDisConnected(f func()) {
	fc.onDisConnected = f
}

func (fc *FastWsClient) HandleMessage(f func(data []byte)) {
	fc.onMessage = f
}

func (fc *FastWsClient) SendMessage(data []byte) error {
	er := fc.conn.WriteMessage(websocket.TextMessage, data)
	if er != nil {
		log.Println(er)
		return er
	}
	return nil
}

func (fc *FastWsClient) Ping() error {
	er := fc.conn.WriteMessage(websocket.PingMessage, []byte(""))
	if er != nil {
		log.Println(er)
		return er
	}
	return nil
}

func (fc *FastWsClient) Run() {
	for {
		_, message, err := fc.conn.ReadMessage()
		if err != nil {
			log.Println(err)

			fc.onDisConnected()

			return
		}
		fc.onMessage(message)
	}
}
