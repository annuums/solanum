package solanum

var helloWorldController *SolaController

func NewHelloWorldController() (*SolaController, error) {
	if helloWorldController == nil {
		helloWorldController, _ = NewController()
		addHandlers()
	}

	return helloWorldController, nil
}

func addHandlers() {
	helloHandler := NewHelloWorldHandler()
	// anotherHandler := NewHelloWorldHandler()
	//* ...

	helloWorldController.AddHandler(helloHandler)
}
