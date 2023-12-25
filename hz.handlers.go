package solanum

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHealthCheckHandler() *SolaService {
	return &SolaService{
		Uri:        "/",
		Method:     http.MethodGet,
		Handler:    hzHandler,
		Middleware: hzMiddleware,
	}
}

func hzHandler(ctx *gin.Context) {
	ctx.String(200, "pong")
}

func hzMiddleware(ctx *gin.Context) {
	log.Println("Health Checking...")
	ctx.Next()
}
