package solanum

var healthCheckController *SolaController

func NewHealthCheckController() (*SolaController, error) {
	if healthCheckController == nil {
		healthCheckController, _ = NewController()
		addHandlers()
	}

	return healthCheckController, nil
}

func addHandlers() {
	healthCheckHandler := NewHealthCheckHandler()
	// anotherHandler := NewHelloWorldHandler()
	//* ...

	healthCheckController.AddHandler(healthCheckHandler)
}
