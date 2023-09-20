package modules

type helloWorldController struct {
	controller
}

func NewHelloWorldController() *helloWorldController {
	return &helloWorldController{}
}
