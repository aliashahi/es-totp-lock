package webserver

import (
	"time"

	"github.com/google/uuid"
)

const static_ui = true
const _COOKIE_TOKEN = "TOKEN"

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Avatar    string    `json:"avatar"`
	Secret    string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Log struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}
