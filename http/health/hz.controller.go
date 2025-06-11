package health

import (
	"github.com/annuums/solanum"
	"net/http"
)

func NewHealthCheckController() *solanum.SolaController {

	healthCheckController := solanum.NewController()

	healthCheckController.SetHandlers(
		&solanum.SolaService{
			Uri:     "",
			Method:  http.MethodGet,
			Handler: hzHandler,
		},
	)

	return healthCheckController
}
