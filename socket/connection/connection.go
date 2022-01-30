package connection

import (
	"github.com/gorilla/websocket"
	"github.com/lithammer/shortuuid"
	"github.com/scarlettmiss/engine-w/socket/hub"
	"sync"
)

type Connection struct {
	id      string
	connMux sync.Mutex
	conn    *websocket.Conn
}

func New(conn *websocket.Conn) (*Connection, error) {
	if conn == nil {
		return nil, hub.ErrInvalidConnection
	}
	return &Connection{
		id:   shortuuid.New(),
		conn: conn,
	}, nil
}

func (c *Connection) Id() string {
	return c.id
}

func (c *Connection) Conn() *websocket.Conn {
	return c.conn
}
