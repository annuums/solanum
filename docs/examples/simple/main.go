package main

import (
	"github.com/annuums/solanum"
	"github.com/gin-gonic/gin"
	"net/http"
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
			Handler: func(c *gin.Context) {
				c.String(http.StatusOK, "pong")
			},
		},
	)
	pingModule.SetControllers(ctrl)

	server := solanum.NewSolanum(
		solanum.WithPort(5050),
	)

	server.SetModules(pingModule)
	server.Run()
}
