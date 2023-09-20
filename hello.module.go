package solanum

import "github.com/gin-gonic/gin"

var helloWorldModule Module

func NewHelloWorldModule(router *gin.RouterGroup, uri string) (Module, error) {
	if helloWorldModule == nil {
		helloWorldModule, _ = NewModule(router, uri)
		attachControllers()
	}

	return helloWorldModule, nil
}

func attachControllers() {
	//* Attatching Controller Directly
	ctr, _ := NewHelloWorldController()
	// ctr2, _ := NewAnotherController()
	//	...

	helloWorldModule.SetControllers(ctr)
}
