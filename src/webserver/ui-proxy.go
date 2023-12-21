package webserver

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func uiReverseProxy() gin.HandlerFunc {
	target, err := url.Parse("http://localhost:4200/")
	if err != nil {
		return nil
	}

	return func(c *gin.Context) {
		rewrite := func(r *httputil.ProxyRequest) {
			r.SetURL(target)
		}
		proxy := &httputil.ReverseProxy{Rewrite: rewrite}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
