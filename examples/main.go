package main

import (
	"net/http/httptest"

	"github.com/andywang2011/go-app-kit/core"
	kitmw "github.com/andywang2011/go-app-kit/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(kitmw.Recovery)
	r.GET("/ping", core.Handle(func(c *core.Context) {
		c.Success(gin.H{"message": "pong"})
	}))

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)
}
