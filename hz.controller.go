package solanum

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func NewHealthCheckController() *SolaController {
	healthCheckController := NewController()

	healthCheckController.SetHandlers(SolaService{
		Uri:        "",
		Method:     http.MethodGet,
		Handler:    hzHandler,
		Middleware: hzMiddleware,
	})

	return healthCheckController
}

func hzMiddleware(ctx *gin.Context) {
	log.Println("Health Checking...")
	ctx.Next()
}
