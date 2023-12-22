package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ClientConnection[T any] struct {
	UserId uuid.UUID
	C      chan T
}

var logConnections = make([]*ClientConnection[string], 0, 10)

func wsLogs(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer conn.Close()

	userId, _ := ctxUserID(ctx)

	var c = make(chan string)
	logConnections = append(logConnections, &ClientConnection[string]{
		UserId: userId,
		C:      c,
	})

	for {
		l := <-c
		conn.WriteMessage(websocket.TextMessage, []byte(l))
	}
}

var userConnections = make([]*ClientConnection[User], 0, 10)

func wsUsers(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer conn.Close()

	userId, _ := ctxUserID(ctx)

	var c = make(chan User)
	userConnections = append(userConnections, &ClientConnection[User]{
		UserId: userId,
		C:      c,
	})

	for {
		l := <-c
		conn.WriteJSON(l)
	}
}
