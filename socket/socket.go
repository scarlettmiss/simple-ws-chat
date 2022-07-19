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
	UserAuthentication     string = "userAuthentication"
	UpdateUser             string = "updateUser"
	CreateRoom             string = "createRoom"
	UserJoinRoom           string = "userJoinRoom"
	UserLeaveRoom          string = "userLeaveRoom"
	UserMessage            string = "userMessage"
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
	Password string `json:"password"`
}

type UserAuthenticationInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserUpdateInfo struct {
	UserId   string  `json:"user_id"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type SessionCreationInfo struct {
	UserId     string             `json:"user_id"`
	Capacity   int                `json:"capacity"`
	Rating     int                `json:"rating"`
	Constraint session.Constraint `json:"constraint"`
}

type UserSessionInfo struct {
	UserId    string `json:"user_id"`
	SessionId string `json:"session_id"`
}

type UserChatMessage struct {
	UserId  string `json:"user_id"`
	Message string `json:"message"`
}

func (api *API) handleDefaultMessage(c socketio.Conn) {
	api.errorMessage(c, "not a valid action")
}

func (api *API) errorMessage(c socketio.Conn, errorMessage string) {
	fmt.Println(errorMessage)
	c.Emit("error", errorMessage)
}

func (api *API) handleUserCreation(message string) string {
	log.Println("create user:", message)

	var userInfo UserCreationInfo
	err := json.Unmarshal([]byte(message), &userInfo)
	if err != nil {
		return err.Error()
	}
	u, err := api.Application.CreateUser(userInfo.Username, userInfo.Password)
	if err != nil {
		return err.Error()
	}
	b, err := json.Marshal(u)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (api *API) handleUserAuthentication(message string) string {
	log.Println("Authenticate user:", message)

	var userInfo UserAuthenticationInfo
	err := json.Unmarshal([]byte(message), &userInfo)
	if err != nil {
		return err.Error()
	}
	err = api.Application.Authenticate(userInfo.Username, userInfo.Password)
	if err != nil {
		return err.Error()
	}
	return "logged in"
}

func (api *API) handleUserRoomCreation(c socketio.Conn, message string) {
	log.Println("create session:", message)

	var sessionInfo SessionCreationInfo
	err := json.Unmarshal([]byte(message), &sessionInfo)
	if err != nil {
		api.errorMessage(c, "Could not Create session")
		return
	}
	s, err := api.CreateSession(sessionInfo.UserId, sessionInfo.Capacity, sessionInfo.Rating, sessionInfo.Constraint)
	if err != nil {
		api.errorMessage(c, "Could not Create session")
		return
	}
	c.Emit("CreatedSession", s.Id())
	c.Join(s.Id())
	c.Emit("joinedSession", s.Id())

	log.Println("Created and joined session:", s.Id())
}

func (api *API) handleUserJoinedRoom(c socketio.Conn, message string) {
	log.Println("join session:", message)

	var sessionInfo UserSessionInfo
	err := json.Unmarshal([]byte(message), &sessionInfo)
	if err != nil {
		api.errorMessage(c, "Could not Join session")
		return
	}
	err = api.JoinSession(sessionInfo.SessionId, sessionInfo.UserId)
	if err != nil {
		api.errorMessage(c, "Could not Join session")
		return
	}
	c.Join(sessionInfo.SessionId)
	log.Println("joined session:", sessionInfo.SessionId)
	c.Emit("joinedSession", sessionInfo.SessionId)
}

func (api *API) handleUserLeavesRoom(c socketio.Conn, message string) {
	log.Println("leave session:", message)

	var sessionInfo UserSessionInfo
	err := json.Unmarshal([]byte(message), &sessionInfo)
	if err != nil {
		api.errorMessage(c, "Could not leave session")
		return
	}
	err = api.LeaveSession(sessionInfo.SessionId, sessionInfo.UserId)
	if err != nil {
		api.errorMessage(c, "Could not leave session")
		return
	}
	c.Leave(sessionInfo.SessionId)
	log.Println("left session:", sessionInfo.SessionId)
	c.Emit("LeftSession", sessionInfo.SessionId)
}

func (api *API) handleChatMessage(c socketio.Conn, message string) {
	log.Println("user message session:", message)

	var chatMessageInfo UserChatMessage
	err := json.Unmarshal([]byte(message), &chatMessageInfo)
	if err != nil {
		api.errorMessage(c, "handleChatMessage: Could not send message")
		return
	}
	s, err := api.UserSession(chatMessageInfo.UserId)
	if err != nil {
		api.errorMessage(c, err.Error())
		return
	}

	api.Server.BroadcastToRoom("/chat", s.Id(), "reply", chatMessageInfo.Message)
	c.Emit("Ack", chatMessageInfo.Message)
}

func (api *API) handleUserMessageRoom(c socketio.Conn, message string) {
	log.Println("message session:", message)

	var sessionInfo UserSessionInfo
	err := json.Unmarshal([]byte(message), &sessionInfo)
	if err != nil {
		api.errorMessage(c, "Could not leave session")
		return
	}
	err = api.LeaveSession(sessionInfo.SessionId, sessionInfo.UserId)
	if err != nil {
		api.errorMessage(c, "Could not leave session")
		return
	}
	c.Leave(sessionInfo.SessionId)
	log.Println("left session:", sessionInfo.SessionId)
	c.Emit("LeftSession", sessionInfo.SessionId)
}

func (api *API) handleUserUpdate(c socketio.Conn, message string) {
	var userInfo UserUpdateInfo
	err := json.Unmarshal([]byte(message), &userInfo)
	if err != nil {
		api.errorMessage(c, "Could not update User")
		return
	}
	u, err := api.Application.UpdateUser(userInfo.UserId, userInfo.Username, userInfo.Password)
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

	api.Server.OnEvent("/", CreateUser, func(c socketio.Conn, msg string) string {
		result := api.handleUserCreation(msg)
		return result
	})

	api.Server.OnEvent("/", UserAuthentication, func(c socketio.Conn, msg string) string {
		result := api.handleUserAuthentication(msg)
		return result
	})

	api.Server.OnEvent("/", UpdateUser, func(c socketio.Conn, msg string) {
		api.handleUserUpdate(c, msg)
	})

	api.Server.OnEvent("/", CreateRoom, func(c socketio.Conn, msg string) {
		api.handleUserRoomCreation(c, msg)
	})

	api.Server.OnEvent("/", UserJoinRoom, func(c socketio.Conn, msg string) {
		api.handleUserJoinedRoom(c, msg)
	})

	api.Server.OnEvent("/", UserLeaveRoom, func(c socketio.Conn, msg string) {
		api.handleUserLeavesRoom(c, msg)
	})

	api.Server.OnEvent("/chat", UserMessage, func(c socketio.Conn, msg string) {
		c.SetContext(msg)
		api.handleChatMessage(c, msg)
		//return "recv " + msg
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
