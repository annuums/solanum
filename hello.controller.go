package solanum

type helloWorldController struct {
	controller
}

func NewHelloWorldController() *helloWorldController {
	return &helloWorldController{}
}
