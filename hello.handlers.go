package solanum

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHelloWorldHandler() *SolaService {
	return &SolaService{
		Uri:        "/",
		Method:     http.MethodGet,
		Handler:    indexHandler,
		Middleware: indexMiddleware,
	}
}

func indexHandler(ctx *gin.Context) {
	ctx.JSON(200, "Hello, World! From HelloWorld Index Handler.")
}

func indexMiddleware(ctx *gin.Context) {
	log.Println("Hello Index Middleware")
	ctx.Next()
}
