package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"

	"github.com/zq-xu/gotools/logx"
)

type client struct {
	Id      string
	Group   string
	lock    *sync.Mutex
	Socket  *websocket.Conn
	Message chan []byte
	ms      []Manager
}

func NewClient(grp string, conn *websocket.Conn, ms ...Manager) *client {
	c := &client{
		Id:      uuid.NewV4().String(),
		Group:   grp,
		lock:    &sync.Mutex{},
		Socket:  conn,
		Message: make(chan []byte, 1024),
		ms:      ms,
	}

	for _, m := range c.ms {
		m.RegisterClient(c)
	}

	return c
}

func (c *client) Read(fn func(p []byte)) {
	defer c.Close()

	for {
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			break
		}

		fn(message)
	}
}

func (c *client) WriteMessage(msg []byte) {
	c.write(websocket.BinaryMessage, msg)
}

func (c *client) WriteCloseMessage() {
	c.write(websocket.CloseMessage, []byte{})
}

func (c *client) write(messageType int, data []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.Socket == nil {
		logx.Logger.Warningf("socket %s/%s is closed.", c.Group, c.Id)
		return
	}

	err := c.Socket.WriteMessage(messageType, data)
	if err != nil {
		logx.Logger.Errorf("socket %s/%s write message err: %s", c.Group, c.Id, err)
	}
}

func (c *client) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	err := c.Socket.Close()
	if err != nil {
		logx.Logger.Errorf("client [%s] disconnect err: %s", c.Id, err)
	}

	c.Socket = nil

	for _, m := range c.ms {
		m.UnregisterClient(c)
	}
}
