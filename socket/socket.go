package socket

import (
	"encoding/json"
	socketio "github.com/googollee/go-socket.io"
	"github.com/scarlettmiss/engine-w/application"
	"github.com/scarlettmiss/engine-w/application/domain/achievement"
	"github.com/scarlettmiss/engine-w/application/domain/session"
	"github.com/scarlettmiss/engine-w/application/domain/user"
	"log"
)

type API struct {
	*application.Application
	*socketio.Server
}

const (
	CreateUser             string = "createUser"
	UserAuthentication     string = "userAuthentication"
	UserLogOut             string = "userLogOut"
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
	RequestUsersScore      string = "requestUsersScore"
	AddUserAchivement      string = "addUserAchivement"
)

func New(application *application.Application) (*API, error) {
	server := socketio.NewServer(nil)
	api := &API{
		Application: application,
		Server:      server,
	}
	return api, nil
}

type UserContext struct {
	UserId string `json:"userId"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

type UserInfoResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Online   bool   `json:"online"`
	Skill    int    `json:"skill"`
}

type GenericResponse struct {
	Message string `json:"message"`
}

type BroadcastActionUsername struct {
	UserName string `json:"username"`
}

type BroadcastUserMessage struct {
	Message  string `json:"message"`
	UserName string `json:"username"`
}

type UserAuthenticationInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserUpdateInfo struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type SessionCreationInfo struct {
	Capacity   int    `json:"capacity"`
	MinRating  int    `json:"minRating"`
	MaxRating  int    `json:"maxRating"`
	Constraint string `json:"constraint"`
}

type UserSessionInfo struct {
	SessionId string `json:"session_id"`
}

type AchievementInfo struct {
	Id string `json:"id"`
}

type UserChatMessage struct {
	Message string `json:"message"`
}

type SessionsInfo struct {
	Id            string             `json:"id"`
	UsersCount    int                `json:"usersCount"`
	Capacity      int                `json:"capacity"`
	MinRating     int                `json:"minRating"`
	Constraint    string             `json:"constraint"`
	OwnerUsername string             `json:"owner"`
	Users         []UserInfoResponse `json:"users"`
}

func (api *API) handleDefaultMessage() string {
	return handleError("Error")
}

func (api *API) handleUserCreation(c socketio.Conn, message string) string {
	log.Println("create user:", message)

	var userInfo UserAuthenticationInfo
	err := json.Unmarshal([]byte(message), &userInfo)
	if err != nil {
		return handleError("Could not create User")
	}
	u, err := api.Application.CreateUser(userInfo.Username, userInfo.Password)
	if err != nil {
		return handleError(err.Error())
	}
	u.SetOnline(true)

	err = api.UserSetOnline(u.Id(), true)
	if err != nil {
		log.Println("Err: ", err.Error())
	}

	c.SetContext(UserContext{UserId: u.Id()})
	return handleResponse(UserInfoResponse{u.Id(), u.Username(), u.Online(), u.SkillPoints()})
}

func (api *API) handleUserAuthentication(c socketio.Conn, message string) string {
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

	err = api.UserSetOnline(u.Id(), true)
	if err != nil {
		log.Println("Err: ", err.Error())
	}

	c.SetContext(UserContext{UserId: u.Id()})
	return handleResponse(UserInfoResponse{Id: u.Id(), Username: u.Username(), Online: u.Online(), Skill: u.SkillPoints()})
}

func (api *API) handleUserDisconnect(c socketio.Conn) string {
	log.Println("log out user:")
	ctx := c.Context().(UserContext)
	u, err := api.User(ctx.UserId)
	if err != nil {
		return handleError(err.Error())
	}
	s, err := api.UserSession(u.Id())
	if err != nil {
		return handleError(err.Error())
	}
	err = api.userLeavesRoomImpl(c, u, s)
	if err != nil {
		return handleError(err.Error())
	}
	err = api.UserSetOnline(u.Id(), false)
	if err != nil {
		log.Println("Err: ", err.Error())
	}

	c.SetContext(UserContext{UserId: ""})

	return handleResponse(GenericResponse{Message: "Disconnected"})
}

func (api *API) handleUserRoomCreation(c socketio.Conn, message string) string {
	log.Println("create session:", message)

	var sessionInfo SessionCreationInfo
	ctx := c.Context().(UserContext)
	err := json.Unmarshal([]byte(message), &sessionInfo)
	if err != nil {
		return handleError(err.Error())
	}
	u, err := api.User(ctx.UserId)
	if err != nil {
		return handleError(err.Error())
	}
	sess, err := api.CreateSession(u, sessionInfo.Capacity, sessionInfo.MinRating, sessionInfo.MaxRating, sessionInfo.Constraint)
	if err != nil {
		return handleError(err.Error())
	}
	api.Server.JoinRoom(c.Namespace(), sess.Id(), c)

	users := api.getUsersResponse(sess.Users())
	return handleResponse(
		SessionsInfo{
			Id:            sess.Id(),
			Capacity:      sess.Capacity(),
			Constraint:    sess.Constraint(),
			MinRating:     sess.MinRating(),
			OwnerUsername: sess.Owner().Username(),
			Users:         users,
		})
}
func (api *API) getUsersResponse(usersMap map[string]*user.User) []UserInfoResponse {
	var users = make([]UserInfoResponse, len(usersMap))
	i := 0
	for _, v := range usersMap {
		users[i] = UserInfoResponse{
			Id:       v.Id(),
			Username: v.Username(),
			Online:   v.Online(),
			Skill:    v.SkillPoints(),
		}
		i++
	}
	return users
}

func (api *API) broadcastSessionMessages(sess *session.Session, c socketio.Conn) {
	for _, v := range sess.Messages() {
		u, err := api.User(v.UserId())
		if err != nil {
			log.Println("Err: " + err.Error())
			continue
		}
		resp := handleResponse(BroadcastUserMessage{Message: v.Message(), UserName: u.Username()})
		log.Println("Message: " + v.Message())

		api.Server.BroadcastToRoom(c.Namespace(), sess.Id(), ChatMessage, resp)
	}
}

func (api *API) joinRoomImpl(sess *session.Session, u *user.User, c socketio.Conn) string {
	sessId := sess.Id()
	api.Server.JoinRoom(c.Namespace(), sessId, c)
	api.broadcastSessionMessages(sess, c)
	api.Server.BroadcastToRoom(c.Namespace(), sessId, UserJoinedRoom, handleResponse(BroadcastActionUsername{UserName: u.Username()}))
	users := api.getUsersResponse(sess.Users())
	return handleResponse(
		SessionsInfo{
			Id:            sess.Id(),
			Capacity:      sess.Capacity(),
			Constraint:    sess.Constraint(),
			MinRating:     sess.MinRating(),
			OwnerUsername: sess.Owner().Username(),
			Users:         users,
		})
}

func (api *API) handleUserJoinedRoom(c socketio.Conn, message string) string {
	log.Println("join session:", message)

	var sessionInfo UserSessionInfo
	err := json.Unmarshal([]byte(message), &sessionInfo)
	if err != nil {
		return handleError("Could not Join session")

	}
	ctx := c.Context().(UserContext)
	u, err := api.User(ctx.UserId)
	if err != nil {
		return handleError(err.Error())
	}

	sess, err := api.JoinSession(sessionInfo.SessionId, u.Id())
	if err != nil {
		return handleError(err.Error())
	}
	return api.joinRoomImpl(sess, u, c)

}

func (api *API) handleMatchMaking(c socketio.Conn) string {
	log.Println("Match-making:")
	ctx := c.Context().(UserContext)
	u, err := api.User(ctx.UserId)
	if err != nil {
		return handleError(err.Error())
	}
	sess, err := api.RequestJoinSession(u.Id())
	if err != nil {
		return handleError(err.Error())
	}
	return api.joinRoomImpl(sess, u, c)
}

func (api *API) handleUserLeavesRoom(c socketio.Conn) string {
	log.Println("leave session:")
	ctx := c.Context().(UserContext)
	u, err := api.User(ctx.UserId)
	if err != nil {
		return handleError(err.Error())
	}
	s, err := api.UserSession(u.Id())
	if err != nil {
		return handleError("Could not Leave Session")
	}
	err = api.userLeavesRoomImpl(c, u, s)
	if err != nil {
		return handleError(err.Error())
	}
	return handleResponse(GenericResponse{Message: "Success"})
}

func (api *API) userLeavesRoomImpl(c socketio.Conn, u *user.User, s *session.Session) error {
	err := api.LeaveSession(s.Id(), u.Id())
	if err != nil {
		return err
	}
	api.Server.BroadcastToRoom(c.Namespace(), s.Id(), UserLeftRoom, handleResponse(BroadcastActionUsername{UserName: u.Username()}))
	api.Server.LeaveRoom(c.Namespace(), s.Id(), c)
	return nil
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
	ctx := c.Context().(UserContext)

	u, err := api.User(ctx.UserId)
	if err != nil {
		return handleError(err.Error())
	}
	s, err := api.UserSession(u.Id())
	if err != nil {
		return handleError(err.Error())
	}

	m, err := api.SessionAddMessage(s.Id(), u.Id(), chatMessageInfo.Message)
	if err != nil {
		return handleError(err.Error())
	}

	resp := handleResponse(BroadcastUserMessage{Message: m.Message(), UserName: u.Username()})

	api.Server.BroadcastToRoom(c.Namespace(), s.Id(), ChatMessage, resp)
	return resp
}

func (api *API) handleUserUpdate(c socketio.Conn, message string) string {
	var userInfo UserUpdateInfo
	err := json.Unmarshal([]byte(message), &userInfo)
	if err != nil {
		return handleError("Could not update User")
	}
	ctx := c.Context().(UserContext)
	u, err := api.User(ctx.UserId)
	if err != nil {
		return handleError(err.Error())
	}

	err = api.UserSetAccountInfo(u.Id(), userInfo.Username, userInfo.Password)
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
		users := api.getUsersResponse(v.Users())
		data[i] = SessionsInfo{
			Id:            v.Id(),
			UsersCount:    len(v.Users()),
			Capacity:      v.Capacity(),
			MinRating:     v.MinRating(),
			Constraint:    v.Constraint(),
			OwnerUsername: v.Owner().Username(),
			Users:         users,
		}
		i++
	}
	return handleResponse(data)
}

func (api *API) handleUsersScoresRequest() string {
	users := api.Users()
	data := make([]UserInfoResponse, len(users))
	i := 0
	for _, v := range users {
		data[i] = UserInfoResponse{
			Id:       v.Id(),
			Username: v.Username(),
			Online:   v.Online(),
			Skill:    v.SkillPoints(),
		}
		i++
	}
	return handleResponse(data)
}

func (api *API) handleAchievementAssignment(c socketio.Conn, message string) string {
	var achievementInfo AchievementInfo
	err := json.Unmarshal([]byte(message), &achievementInfo)
	if err != nil {
		return handleError(achievement.ErrNotFound.Error())

	}
	ctx := c.Context().(UserContext)
	u, err := api.User(ctx.UserId)
	if err != nil {
		handleError(err.Error())
	}

	err = api.UserAddAchievement(u.Id(), achievementInfo.Id)
	if err != nil {
		handleError(err.Error())
	}

	return handleResponse(GenericResponse{Message: "Updated"})
}

func (api *API) CreateHandlers() {
	api.Server.OnConnect("/", func(c socketio.Conn) error {
		c.SetContext(UserContext{})
		log.Println("connected:", c.ID())
		return nil
	})

	api.Server.OnDisconnect("/", func(c socketio.Conn, msg string) {
		log.Println("disconnected:", c.ID(), msg)
		api.handleUserDisconnect(c)
		err := c.Close()
		if err != nil {
			log.Println("Err: " + err.Error())
		}
	})

	api.Server.OnEvent("/", CreateUser, func(c socketio.Conn, msg string) string {
		result := api.handleUserCreation(c, msg)
		return result
	})

	api.Server.OnEvent("/", UserAuthentication, func(c socketio.Conn, msg string) string {
		return api.handleUserAuthentication(c, msg)
	})

	api.Server.OnEvent("/", UpdateUser, func(c socketio.Conn, msg string) string {
		return api.handleUserUpdate(c, msg)
	})

	api.Server.OnEvent("/", CreateRoom, func(c socketio.Conn, msg string) string {
		return api.handleUserRoomCreation(c, msg)
	})

	api.Server.OnEvent("/", UserJoinRoom, func(c socketio.Conn, msg string) string {
		return api.handleUserJoinedRoom(c, msg)
	})

	api.Server.OnEvent("/", UserLeaveRoom, func(c socketio.Conn, msg string) string {
		return api.handleUserLeavesRoom(c)
	})

	api.Server.OnEvent("/", UserMessage, func(c socketio.Conn, msg string) string {
		return api.handleChatMessage(c, msg)
	})

	api.Server.OnEvent("/", getRooms, func(c socketio.Conn) string {
		return api.handleSessionsRequest()
	})

	api.Server.OnEvent("/", UserRequestMatchMaking, func(c socketio.Conn, msg string) string {
		return api.handleMatchMaking(c)
	})

	api.Server.OnEvent("/", RequestUsersScore, func(c socketio.Conn, msg string) string {
		return api.handleUsersScoresRequest()
	})

	api.Server.OnEvent("/", AddUserAchivement, func(c socketio.Conn, msg string) string {
		return api.handleAchievementAssignment(c, msg)
	})

	api.Server.OnEvent("/", UserLogOut, func(c socketio.Conn) {
		log.Println("logged out:", c.ID())
		api.handleUserDisconnect(c)
		err := c.Close()
		if err != nil {
			log.Println("Err: " + err.Error())
		}
	})

	go func() {
		if err := api.Server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()

}
