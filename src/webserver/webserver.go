package webserver

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
}

var clientChannels = make([]Client, 0, 10)

func tokenMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie(_COOKIE_TOKEN)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("login"))
			return
		}

		userId, err := uuid.Parse(token)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		found := false
		for _, u := range users {
			if u.ID == userId {
				found = true
				break
			}
		}

		if !found {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set(_COOKIE_TOKEN, token)
	}
}

func ctxUserID(ctx *gin.Context) (uuid.UUID, error) {
	s := ctx.GetString(_COOKIE_TOKEN)
	return uuid.Parse(s)
}

func WebServer() {
	engine := gin.Default()
	if !static_ui {
		engine.NoRoute(uiReverseProxy())
	} else {
		engine.Use(static.Serve("/", static.LocalFile("./ui/dist", true)))
	}

	router := engine.Group("api")

	router.POST("/login", login)
	router.GET("/userinfo", tokenMiddleware(), userInfo)
	router.POST("/create", tokenMiddleware(), create)
	router.GET("/code", tokenMiddleware(), code)
	router.GET("/users", tokenMiddleware(), allUsers)
	router.GET("/logs", tokenMiddleware(), allLogs)
	router.GET("/ws", tokenMiddleware(), ws)

	if err := engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func ws(ctx *gin.Context) {
	userId, err := ctxUserID(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	clientChannels = append(clientChannels, Client{
		UserID: userId,
		Conn:   conn,
	})
}

func userInfo(ctx *gin.Context) {
	userId, err := ctxUserID(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	for _, u := range users {
		if u.ID == userId {
			ctx.JSON(http.StatusCreated, gin.H{
				"isAdmin": u.Username == "admin",
			})
			return
		}
	}

	ctx.Status(http.StatusUnauthorized)
}

func login(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	for _, u := range users {
		if u.Username == req.Username && u.Password == req.Password {
			ctx.SetCookie(_COOKIE_TOKEN, fmt.Sprint(u.ID), 1000, "/", "", false, true)

			ctx.JSON(http.StatusCreated, gin.H{
				"isAdmin": u.Username == "admin",
			})

			return
		}
	}

	ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("username/password is wrong"))
}

func create(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	for _, u := range users {
		if u.Username == req.Username {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("username already exists"))
			return
		}
	}

	secret := uuid.NewString()

	new_user := User{
		ID:       uuid.New(),
		Username: req.Username,
		Password: req.Password,
		Code:     generatePassCode(secret),
		Secret:   secret,
	}

	users = append(users, &new_user)

	ctx.AbortWithStatus(http.StatusCreated)
}

func code(ctx *gin.Context) {
	userId, err := ctxUserID(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	for _, u := range users {
		if u.ID == userId {
			ctx.JSON(http.StatusOK, gin.H{"code": u.Code})
			return
		}
	}

	ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("user not found"))
}

func allUsers(ctx *gin.Context) {
	if !isAdmin(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func isAdmin(ctx *gin.Context) bool {
	userId, _ := ctxUserID(ctx)
	for _, u := range users {
		if u.ID == userId {
			return u.Username == "admin"
		}
	}

	panic("user not found")
}

func allLogs(ctx *gin.Context) {
	if !isAdmin(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	raw := ctx.Query("lastId")
	lastLogId, err := strconv.ParseInt(raw, 10, 64)
	if raw == "" || err != nil {
		ctx.JSON(http.StatusOK, gin.H{"logs": logs})
		return
	}

	_logs := make([]*Log, 0, len(logs))

	for _, l := range logs {
		if l.ID > lastLogId {
			_logs = append(_logs, l)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"logs": _logs})
}

func GetUserByPasscode(passcode string) (*User, error) {
	for _, u := range users {
		if u.Code == passcode {
			return u, nil
		}
	}
	return nil, fmt.Errorf("wrong code")
}
