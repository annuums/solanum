package solanum

import (
	"github.com/gin-gonic/gin"
	"log"
)

func hzHandler(ctx *gin.Context) {
	ctx.String(200, "pong")
}

func hzMiddleware(ctx *gin.Context) {
	log.Println("Health Checking...")
	ctx.Next()
}
