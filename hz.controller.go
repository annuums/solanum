package solanum

import (
	"net/http"
)

func NewHealthCheckController() *SolaController {
	healthCheckController := NewController()

	healthCheckController.SetHandlers(SolaService{
		Uri:     "",
		Method:  http.MethodGet,
		Handler: hzHandler,
	})

	return healthCheckController
}
