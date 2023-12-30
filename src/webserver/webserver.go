package webserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	qrcode "github.com/skip2/go-qrcode"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebServer() {
	engine := gin.Default()
	if !static_ui {
		engine.NoRoute(uiReverseProxy())
	} else {
		engine.Use(static.Serve("/", static.LocalFile("./ui/dist/ui", true)))
	}

	router := engine.Group("api")

	router.POST("/login", login)
	router.GET("/logout", tokenMiddleware(), logout)
	router.GET("/userinfo", tokenMiddleware(), userInfo)
	router.POST("/create", tokenMiddleware(), create)
	router.DELETE("/:id", tokenMiddleware(), delete)
	router.GET("/code", tokenMiddleware(), code)
	router.GET("/users", tokenMiddleware(), allUsers)
	router.GET("/logs", tokenMiddleware(), allLogs)
	router.GET("/ws-logs", tokenMiddleware(), wsLogs)
	router.GET("/ws-users", tokenMiddleware(), wsUsers)
	router.GET("/validate", validate)

	if err := engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func validate(ctx *gin.Context) {
	q := ctx.Query("code")
	u, err := GetUserByPasscode([]byte(q))
	if err != nil {
		Logger("%s (%s)", err.Error(), q)
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.String(http.StatusOK, fmt.Sprintf("correct code for user %s", u.Username))
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
				"isAdmin":  u.IsAdmin(),
				"username": u.Username,
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

func logout(ctx *gin.Context) {

	ctx.SetCookie(_COOKIE_TOKEN, "", 0, "/", "", false, true)

	ctx.AbortWithStatus(http.StatusUnauthorized)
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

	if _, err := createUser(req.Username, req.Password); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.AbortWithStatus(http.StatusCreated)
}

func delete(ctx *gin.Context) {
	if !isAdmin(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	id_raw := ctx.Param("id")
	id, err := uuid.Parse(id_raw)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = deleteUser(id)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
}

func code(ctx *gin.Context) {
	userId, err := ctxUserID(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	for _, u := range users {
		if u.ID == userId {
			var png []byte
			content := fmt.Sprintf("otpauth://totp/ESP32-%s?secret=%s", u.Username, u.Secret)
			png, err := qrcode.Encode(content, qrcode.Medium, 512)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"secret": png})
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

func allLogs(ctx *gin.Context) {
	if !isAdmin(ctx) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"logs": logs})
}
