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

type DefaultMessage string

func (api *API) handleDefaultMessage(userId string, message DefaultMessage) {
	//fetch user from id
	//do something with user

	//1. either send response to user
	//2. or send something to some other user
	api.broadcast(userId, []byte(userId+" said "+string(message)))
}

func (api *API) handleMessage(userId string, message []byte) {
	messageType := ""
	switch messageType {
	default:
		api.handleDefaultMessage(userId, DefaultMessage(message))
	}
}

func (api *API) broadcast(senderId string, message []byte) {
	connections, err := api.connectionRepo.Connections()
	if err != nil {
		return
	}

	for _, conn := range connections {
		if conn.Id() == senderId {
			return
		}
		if err := conn.Conn().WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println(err, message)
			continue
		}
	}

}

func (api *API) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := api.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c, err := api.connectionRepo.CreateConnection(conn)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("new connection!", c.Id())
	u, err := api.Application.CreateUser()
	fmt.Println("new user!", u.Id())
	if err != nil {
		log.Println(err)
		return
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			return
		}

		api.handleMessage(c.Id(), message)
	}
}
