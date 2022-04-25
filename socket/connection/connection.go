package connection

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/lithammer/shortuuid"
	"sync"
)

var (
	//ErrInvalidConnection invalid connection
	ErrInvalidConnection = errors.New("connection not provided")
)

type Connection struct {
	id      string
	connMux sync.Mutex
	conn    *websocket.Conn
}

func New(conn *websocket.Conn) (*Connection, error) {
	if conn == nil {
		return nil, ErrInvalidConnection
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
