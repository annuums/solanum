package solanum

import (
	"log"

	"github.com/gin-gonic/gin"
)

var helathCheckModule *SolaModule

func NewHealthCheckModule(router *gin.RouterGroup, uri string) (*SolaModule, error) {
	if helathCheckModule == nil {
		helathCheckModule, _ = NewModule(router, uri)
		attachControllers()
	}

	return helathCheckModule, nil
}

func attachControllers() {
	//* Attatching Controller Directly
	var (
		ctr Controller
		err error
	)
	ctr, err = NewHealthCheckController()

	if err != nil {
		log.Fatal(err)
	}
	// ctr2, _ := NewAnotherController()
	//	...

	helathCheckModule.SetControllers(&ctr)
}
