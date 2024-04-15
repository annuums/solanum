package solanum

import (
	"sync"
)

var helathCheckModule *SolaModule
var hzOnce sync.Once

func NewHealthCheckModule(uri string) *SolaModule {
	hzOnce.Do(func() {
		if helathCheckModule == nil {
			helathCheckModule = NewModule(uri)
			attachControllers()
		}
	})

	return helathCheckModule
}

func attachControllers() {
	//* Attatching Controller Directly
	ctr := NewHealthCheckController()
	// ctr2, _ := NewAnotherController()
	//	...

	helathCheckModule.SetControllers(ctr)
}
