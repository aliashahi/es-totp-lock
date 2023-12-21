package webserver

import "github.com/google/uuid"

const static_ui = false
const _COOKIE_TOKEN = "TOKEN"

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:",omitempty"`
	Code     string    `json:",omitempty"`
	Secret   string    `json:",omitempty"`
}

type Log struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

var users = make([]*User, 0, 10)
var logs = make([]*Log, 0, 10)

func init() {
	secret := uuid.NewString()
	users = append(users, &User{
		ID:       uuid.New(),
		Username: "admin",
		Password: "admin",
		Code:     generatePassCode(secret),
		Secret:   secret,
	})
}
