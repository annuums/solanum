package solanum

var helloWorldController *controller

func NewHelloWorldController() (*controller, error) {
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
