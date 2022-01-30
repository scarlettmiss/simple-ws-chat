package hub

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/scarlettmiss/engine-w/application/domain/session"
	"github.com/scarlettmiss/engine-w/socket/connection"
	"sync"
)

var (
	//ErrInvalidConnection invalid connection
	ErrInvalidConnection = errors.New("connection not provided")
)

type Repository struct {
	mux     sync.Mutex
	clients map[string]*connection.Connection
}

func New() *Repository {
	return &Repository{
		clients: map[string]*connection.Connection{},
	}
}

func (r *Repository) CreateConnection(conn *websocket.Conn) (*connection.Connection, error) {
	r.mux.Lock()
	defer r.mux.Unlock()
	c, err := connection.New(conn)
	if err != nil {
		return nil, err
	}

	r.clients[c.Id()] = c

	return c, nil
}

func (r *Repository) Connection(id string) (*connection.Connection, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	c, ok := r.clients[id]
	if !ok {
		return nil, session.ErrNotFound
	}

	return c, nil
}

func (r *Repository) Connections() (map[string]*connection.Connection, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.clients, nil
}

func (r *Repository) Disconnect(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.clients[id]
	if !ok {
		return session.ErrNotFound
	}

	delete(r.clients, id)

	return nil
}
