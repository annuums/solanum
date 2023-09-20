package solanum

type helloWorldModule struct {
	module
}

func NewHelloWorld() *helloWorldModule {
	controller := NewHelloWorldController()

	var helloModule Module = &helloWorldModule{}

	helloModule.SetControllers(controller)

	return helloModule.(*helloWorldModule)
}
