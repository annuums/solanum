package main

import (
	"net/http"

	"github.com/annuums/solanum"
	"github.com/gin-gonic/gin"
)

func main() {

	pingModule := solanum.NewModule(
		solanum.WithUri("/ping"),
	)

	ctrl := solanum.NewController()
	ctrl.SetHandlers(
		&solanum.SolaService{
			Uri:    "",
			Method: http.MethodGet,
			Handler: func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "pong")
			},
		},
	)
	pingModule.AddControllers(ctrl)

	server := solanum.NewSolanum(
		solanum.WithPort(5050),
	)

	server.SetModules(pingModule)
	server.Run()
}
