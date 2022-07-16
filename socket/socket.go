package socket

import (
	"encoding/json"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/scarlettmiss/engine-w/application"
	"github.com/scarlettmiss/engine-w/application/domain/session"
	"log"
)

type API struct {
	*application.Application
	*socketio.Server
}

const (
	CreateUser             string = "createUser"
	UpdateUser             string = "updateUser"
	CreateRoom             string = "createRoom"
	UserJoinRoom           string = "userJoinRoom"
	UserAddFriend          string = "userAddFriend"
	UserRequestMatchMaking string = "userRequestMatchMaking"
)

func New(application *application.Application) (*API, error) {
	server := socketio.NewServer(nil)
	api := &API{
		Application: application,
		Server:      server,
	}
	return api, nil
}

type UserCreationInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateInfo struct {
	UserId   string  `json:"user_id"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type SessionInfo struct {
	UserId     string             `json:"userId"`
	Capacity   int                `json:"capacity"`
	Rating     int                `json:"rating"`
	Constraint session.Constraint `json:"constraint"`
}

func (api *API) handleDefaultMessage(c socketio.Conn) {
	api.errorMessage(c, "not a valid action")
}

func (api *API) errorMessage(c socketio.Conn, errorMessage string) {
	fmt.Println(errorMessage)
	c.Emit("error", errorMessage)
}

func (api *API) handleUserCreation(c socketio.Conn, message string) {
	log.Println("create user:", message)

	var userInfo UserCreationInfo
	err := json.Unmarshal([]byte(message), &userInfo)
	if err != nil {
		api.errorMessage(c, "Could not Create User")
		return
	}
	u, err := api.Application.CreateUser(userInfo.Username, userInfo.Email, userInfo.Password)
	if err != nil {
		api.errorMessage(c, "Could not Create User")
		return
	}
	c.Emit("userCreated", u)
}

func (api *API) handleUserRoomCreation(c socketio.Conn, message string) {
	log.Println("create session:", message)

	var sessionInfo SessionInfo
	err := json.Unmarshal([]byte(message), &sessionInfo)
	if err != nil {
		api.errorMessage(c, "Could not Create room")
		return
	}
	s, err := api.CreateSession(sessionInfo.UserId, sessionInfo.Capacity, sessionInfo.Rating, sessionInfo.Constraint)
	if err != nil {
		api.errorMessage(c, "Could not Create room")
		return
	}
	c.Join(s.Id())
	log.Println("joined session:", s.Id())
	c.Emit("jointRoom", s.Id())
}

func (api *API) handleUserUpdate(c socketio.Conn, message string) {
	var userInfo UserUpdateInfo
	err := json.Unmarshal([]byte(message), &userInfo)
	if err != nil {
		api.errorMessage(c, "Could not update User")
		return
	}
	u, err := api.Application.UpdateUser(userInfo.UserId, userInfo.Username, userInfo.Email, userInfo.Password)
	if err != nil {
		api.errorMessage(c, "Could not update User")
		return
	}
	c.Emit("userUpdated", u)

}

func (api *API) CreateHandlers() {

	api.Server.OnConnect("/", func(c socketio.Conn) error {
		c.SetContext("")
		log.Println("connected:", c.ID())
		return nil
	})

	api.Server.OnEvent("/", "notice", func(c socketio.Conn, msg string) {
		log.Println("notice:", msg)
		c.Emit("reply", "have "+msg)
	})

	api.Server.OnEvent("/", CreateUser, func(c socketio.Conn, msg string) {
		api.handleUserCreation(c, msg)
	})

	api.Server.OnEvent("/", UpdateUser, func(c socketio.Conn, msg string) {
		api.handleUserUpdate(c, msg)
	})

	api.Server.OnEvent("/", CreateRoom, func(c socketio.Conn, msg string) {
		api.handleUserRoomCreation(c, msg)
	})

	api.Server.OnEvent("/chat", "msg", func(c socketio.Conn, msg string) string {
		c.SetContext(msg)
		return "recv " + msg
	})

	api.Server.OnEvent("/", "bye", func(c socketio.Conn) string {
		last := c.Context().(string)
		c.Emit("bye", last)
		c.Close()
		return last
	})

	go func() {
		if err := api.Server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()

}
