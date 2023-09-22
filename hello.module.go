package solanum

import (
	"log"

	"github.com/gin-gonic/gin"
)

var helloWorldModule *SolaModule

func NewHelloWorldModule(router *gin.RouterGroup, uri string) (*SolaModule, error) {
	if helloWorldModule == nil {
		helloWorldModule, _ = NewModule(router, uri)
		attachControllers()
	}

	return helloWorldModule, nil
}

func attachControllers() {
	//* Attatching Controller Directly
	var (
		ctr Controller
		err error
	)
	ctr, err = NewHelloWorldController()

	if err != nil {
		log.Fatal(err)
	}
	// ctr2, _ := NewAnotherController()
	//	...

	helloWorldModule.SetControllers(&ctr)
}
