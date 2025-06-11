package health

import (
	"github.com/gin-gonic/gin"
)

func hzHandler(ctx *gin.Context) {

	ctx.String(200, "pong")
}
