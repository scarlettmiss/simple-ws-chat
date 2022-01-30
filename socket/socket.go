package socket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/scarlettmiss/engine-w/application"
	"github.com/scarlettmiss/engine-w/socket/hub"
	"log"
	"net/http"
)

type API struct {
	*application.Application
	upgrader       websocket.Upgrader
	connectionRepo *hub.Repository
}

func New(application *application.Application) (*API, error) {
	api := &API{
		Application: application,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		connectionRepo: hub.New(),
	}
	return api, nil
}

func (api *API) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := api.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("new client")
	_, err = api.connectionRepo.CreateConnection(conn)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = api.Application.CreateUser()
	if err != nil {
		log.Println(err)
		return
	}
}
