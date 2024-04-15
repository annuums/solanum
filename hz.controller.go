package solanum

var healthCheckController *SolaController

func NewHealthCheckController() (*SolaController, error) {
	if healthCheckController == nil {
		healthCheckController = NewController()
		addHandlers()
	}

	return healthCheckController, nil
}

func addHandlers() {
	healthCheckHandler := NewHealthCheckHandler()
	// anotherHandler := NewHelloWorldHandler()
	//* ...

	healthCheckController.SetHandlers(*healthCheckHandler)
}
