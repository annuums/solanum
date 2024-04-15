package solanum

import (
	"log"
)

var helathCheckModule *SolaModule

func NewHealthCheckModule(uri string) *SolaModule {
	if helathCheckModule == nil {
		helathCheckModule = NewModule(uri)
		attachControllers()
	}

	return helathCheckModule
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

	helathCheckModule.SetControllers(ctr)
}
