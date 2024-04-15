package solanum

import "net/http"

var healthCheckController *SolaController

func NewHealthCheckController() (*SolaController, error) {
	if healthCheckController == nil {
		healthCheckController = NewController()
		addHandlers()
	}

	return healthCheckController, nil
}

func addHandlers() {
	healthCheckController.SetHandlers(SolaService{
		Uri:        "",
		Method:     http.MethodGet,
		Handler:    hzHandler,
		Middleware: hzMiddleware,
	})
}
