package appkit_test

import (
	"fmt"
	"net/http/httptest"

	"github.com/andywang2011/go-app-kit/core"
	kitmw "github.com/andywang2011/go-app-kit/middleware"
	"github.com/gin-gonic/gin"
)

func Example_middleware() {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(kitmw.Recovery)
	r.GET("/ping", core.Handle(func(c *core.Context) {
		c.Success(gin.H{"ok": true})
	}))

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)
	fmt.Println(w.Code)
	// Output: 200
}
