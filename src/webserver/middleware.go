package webserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

func isAdmin(ctx *gin.Context) bool {
	userId, _ := ctxUserID(ctx)
	for _, u := range users {
		if u.ID == userId {
			return u.Username == "admin"
		}
	}

	panic("user not found")
}
