package solanum

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHelloWorldHandler() *service {
	return &service{
		uri:        "/",
		method:     http.MethodGet,
		handler:    indexHandler,
		middleware: indexMiddleware,
	}
}

func indexHandler(ctx *gin.Context) {
	ctx.JSON(200, "Hello, World! From HelloWorld Index Handler.")
}

func indexMiddleware(ctx *gin.Context) {
	log.Println("Hello Index Middleware")
	ctx.Next()
}
