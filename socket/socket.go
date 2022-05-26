package socket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/scarlettmiss/engine-w/application"
	"github.com/scarlettmiss/engine-w/application/domain/user"
	"github.com/scarlettmiss/engine-w/socket/hub"
	"log"
	"net/http"
)

type API struct {
	*application.Application
	upgrader       websocket.Upgrader
	connectionRepo *hub.Repository
}

type MessageType string

const (
	CreateUser             MessageType = "createUser"
	UpdateUser             MessageType = "updateUser"
	CreateRoom             MessageType = "createRoom"
	UserJoinRoom           MessageType = "userJoinRoom"
	UserAddFriend          MessageType = "userAddFriend"
	UserRequestMatchMaking MessageType = "userRequestMatchMaking"
)

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

type UserInfoMessage []byte

type UserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateMessage string

type ErrorMessage struct {
	Error string
}

type UserCreatedEvent struct {
	Type MessageType
	User *user.User
}

func (api *API) handleDefaultMessage(connectionId string) {
	api.errorMessage(connectionId, "not a valid action")
}

func (api *API) errorMessage(connectionId string, errorMessage string) {
	fmt.Println(errorMessage)
	m, err := json.Marshal(ErrorMessage{Error: errorMessage})
	if err != nil {
		fmt.Println("Could not marshal error message", err)
		return
	}

	api.sendSelf(connectionId, m)
}

func (api *API) handleUserCreation(connectionId string, message UserInfoMessage) {
	var userInfo UserInfo
	err := json.Unmarshal(message, &userInfo)
	if err != nil {
		api.errorMessage(connectionId, "Could not Create User")
		return
	}
	u, err := api.Application.CreateUser(userInfo.Username, userInfo.Email, userInfo.Password)
	if err != nil {
		api.errorMessage(connectionId, "Could not Create User")
		return
	}

	userCreated, err := json.Marshal(UserCreatedEvent{Type: "userCreated", User: u})
	if err != nil {
		fmt.Println("Could not marshal error message", err)
		return
	}
	api.sendSelf(connectionId, userCreated)

}

func (api *API) handleUserUpdate(message UserInfoMessage) {
	//fetch user from id
	//do something with user

	//1. either send response to user
	//2. or send something to some other user
	//api.broadcast(userId, []byte(userId+" said "+string(message)))
}

func (api *API) handleMessage(connectionId string, message []byte) {
	var messageType MessageType
	err := json.Unmarshal(message, &messageType)
	if err != nil {
		log.Println(err, message)
		return
	}

	switch messageType {
	case CreateUser:
		api.handleUserCreation(connectionId, message)
	case UpdateUser:
		api.handleUserCreation(connectionId, message)
	default:
		api.handleDefaultMessage(connectionId)
	}
}

func (api *API) broadcast(connectionId string, message []byte) {
	connections, err := api.connectionRepo.Connections()
	if err != nil {
		return
	}

	for _, conn := range connections {
		if conn.Id() == connectionId {
			return
		}
		if err := conn.Conn().WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println(err, message)
			continue
		}
	}
}

func (api *API) sendSelf(connectionId string, message []byte) {
	conn, err := api.connectionRepo.Connection(connectionId)
	if err != nil {
		log.Fatal("could not find interface with id: ", connectionId, err)
		return
	}

	err = conn.Conn().WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Message could not be send to self", string(message), err)
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
