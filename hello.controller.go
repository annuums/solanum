package solanum

var helloWorldController Controller

func NewHelloWorldController() (Controller, error) {
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
