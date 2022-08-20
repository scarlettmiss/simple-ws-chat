package socket

import (
	"encoding/json"
	socketio "github.com/googollee/go-socket.io"
	"github.com/scarlettmiss/engine-w/application"
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
	getRooms               string = "getRooms"
	UserJoinRoom           string = "userJoinRoom"
	UserJoinedRoom         string = "userJoinedRoom"
	UserLeaveRoom          string = "userLeaveRoom"
	UserLeftRoom           string = "userLeftRoom"
	UserMessage            string = "userMessage"
	ChatMessage            string = "chatMessage"
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

type ErrorMessage struct {
	Error string `json:"error"`
}

type UserInfoResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type UserCreateSessionResponse struct {
	Id string `json:"id"`
}

type UserJoinSessionResponse struct {
	Id string `json:"id"`
}

type UserJoinedSessionResponse struct {
	UserName string `json:"username"`
}

type UserLeaveSessionResponse struct {
	Id string `json:"id"`
}

type UserLeftSessionResponse struct {
	UserName string `json:"username"`
}

type UserChatResponse struct {
	Message  string `json:"message"`
	UserName string `json:"username"`
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
	UserId     string `json:"user_id"`
	Capacity   int    `json:"capacity"`
	Rating     int    `json:"rating"`
	Constraint string `json:"constraint"`
}

type UserSessionInfo struct {
	UserId    string `json:"user_id"`
	SessionId string `json:"session_id"`
}

type UserRoomInfo struct {
	UserId string `json:"user_id"`
}

type UserChatMessage struct {
	UserId  string `json:"user_id"`
	Message string `json:"message"`
}

type SessionsInfo struct {
	Id         string `json:"id"`
	UsersCount int    `json:"usersCount"`
	Capacity   int    `json:"capacity"`
	MinRating  int    `json:"minRating"`
	Constraint string `json:"constraint"`
	Owner      string `json:"owner"`
}

func (api *API) handleDefaultMessage() string {
	return handleError("Error")
}

func (api *API) handleUserCreation(message string) string {
	log.Println("create user:", message)

	var userInfo UserCreationInfo
	err := json.Unmarshal([]byte(message), &userInfo)
	if err != nil {
		return handleError("Could not create User")
	}
	u, err := api.Application.CreateUser(userInfo.Username, userInfo.Password)
	if err != nil {
		return handleError(err.Error())
	}
	return handleResponse(UserInfoResponse{u.Id, u.Username})
}

func (api *API) handleUserAuthentication(message string) string {
	log.Println("Authenticate user:", message)

	var userInfo UserAuthenticationInfo
	err := json.Unmarshal([]byte(message), &userInfo)
	if err != nil {
		return handleError("Something went wrong! Please try again later")
	}
	u, err := api.Application.Authenticate(userInfo.Username, userInfo.Password)
	if err != nil {
		return handleError(err.Error())
	}
	return handleResponse(UserInfoResponse{Id: u.Id, Username: u.Username})
}

func (api *API) handleUserRoomCreation(c socketio.Conn, message string) string {
	log.Println("create session:", message)

	var sessionInfo SessionCreationInfo
	err := json.Unmarshal([]byte(message), &sessionInfo)
	if err != nil {
		return handleError(err.Error())
	}
	u, err := api.User(sessionInfo.UserId)
	if err != nil {
		return handleError(err.Error())
	}
	s, err := api.CreateSession(u.Id, sessionInfo.Capacity, sessionInfo.Rating, sessionInfo.Constraint)
	if err != nil {
		return handleError(err.Error())
	}
	api.Server.JoinRoom(c.Namespace(), s.Id(), c)
	return handleResponse(UserCreateSessionResponse{Id: s.Id()})
}

func (api *API) handleUserJoinedRoom(c socketio.Conn, message string) string {
	log.Println("join session:", message)

	var sessionInfo UserSessionInfo
	err := json.Unmarshal([]byte(message), &sessionInfo)
	if err != nil {
		return handleError("Could not Join session")

	}
	u, err := api.User(sessionInfo.UserId)
	if err != nil {
		return handleError(err.Error())
	}
	err = api.JoinSession(sessionInfo.SessionId, u.Id)
	if err != nil {
		return handleError(err.Error())
	}
	api.Server.JoinRoom(c.Namespace(), sessionInfo.SessionId, c)
	api.Server.BroadcastToRoom(c.Namespace(), sessionInfo.SessionId, UserJoinedRoom, handleResponse(UserJoinedSessionResponse{UserName: u.Username}))
	return handleResponse(UserJoinSessionResponse{Id: sessionInfo.SessionId})
}

func (api *API) handleUserLeavesRoom(c socketio.Conn, message string) string {
	log.Println("leave session:", message)

	var userRoomInfo UserRoomInfo
	err := json.Unmarshal([]byte(message), &userRoomInfo)
	if err != nil {
		return handleError(err.Error())
	}
	u, err := api.User(userRoomInfo.UserId)
	if err != nil {
		return handleError(err.Error())
	}
	s, err := api.UserSession(u.Id)
	if err != nil {
		return handleError("Could not Leave Session")
	}
	err = api.LeaveSession(s.Id(), u.Id)
	if err != nil {
		return handleError(err.Error())
	}
	api.Server.BroadcastToRoom(c.Namespace(), s.Id(), UserLeftRoom, handleResponse(UserLeftSessionResponse{UserName: u.Username}))
	api.Server.LeaveRoom(c.Namespace(), s.Id(), c)
	return handleResponse(UserLeaveSessionResponse{Id: s.Id()})
}

func handleError(errMsg string) string {
	errMessage := ErrorMessage{Error: errMsg}
	parsedError, err := json.Marshal(errMessage)
	if err != nil {
		return err.Error()
	}
	return string(parsedError)
}

func handleResponse(resp any) string {
	parsedResp, err := json.Marshal(resp)
	if err != nil {
		return "could not parse"
	}
	return string(parsedResp)
}

func (api *API) handleChatMessage(c socketio.Conn, message string) string {
	log.Println("user message session:", message)

	var chatMessageInfo UserChatMessage
	err := json.Unmarshal([]byte(message), &chatMessageInfo)
	if err != nil {
		return handleError("Could not send message")
	}
	u, err := api.User(chatMessageInfo.UserId)
	if err != nil {
		return handleError(err.Error())
	}
	s, err := api.UserSession(u.Id)
	if err != nil {
		return handleError(err.Error())
	}
	resp := handleResponse(UserChatResponse{Message: chatMessageInfo.Message, UserName: u.Username})

	api.Server.BroadcastToRoom(c.Namespace(), s.Id(), ChatMessage, resp)
	return resp
}

func (api *API) handleUserUpdate(message string) string {
	var userInfo UserUpdateInfo
	err := json.Unmarshal([]byte(message), &userInfo)
	if err != nil {
		return handleError("Could not update User")
	}
	u, err := api.Application.UpdateUser(userInfo.UserId, userInfo.Username, userInfo.Password)
	if err != nil {
		return handleError(err.Error())
	}
	return handleResponse(u)
}

func (api *API) handleSessionsRequest() string {
	sessions := api.GetSessions()
	data := make([]SessionsInfo, len(sessions))
	i := 0
	for _, v := range sessions {
		data[i] = SessionsInfo{
			Id:         v.Id(),
			UsersCount: len(v.Users()),
			Capacity:   v.Capacity(),
			MinRating:  v.MinRating(),
			Constraint: v.Constraint(),
			Owner:      v.Owner(),
		}
		i++
	}
	return handleResponse(data)
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
		return api.handleUserAuthentication(msg)
	})

	api.Server.OnEvent("/", UpdateUser, func(c socketio.Conn, msg string) string {
		return api.handleUserUpdate(msg)
	})

	api.Server.OnEvent("/", CreateRoom, func(c socketio.Conn, msg string) string {
		return api.handleUserRoomCreation(c, msg)
	})

	api.Server.OnEvent("/", UserJoinRoom, func(c socketio.Conn, msg string) string {
		return api.handleUserJoinedRoom(c, msg)
	})

	api.Server.OnEvent("/", UserLeaveRoom, func(c socketio.Conn, msg string) string {
		return api.handleUserLeavesRoom(c, msg)
	})

	api.Server.OnEvent("/", UserMessage, func(c socketio.Conn, msg string) string {
		return api.handleChatMessage(c, msg)
	})

	api.Server.OnEvent("/", getRooms, func(c socketio.Conn) string {
		return api.handleSessionsRequest()
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
