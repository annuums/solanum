package solanum

import (
	"github.com/gin-gonic/gin"
	"log"
	"sync"
)

var healthCheckModule *SolaModule
var hzOnce sync.Once

func NewHealthCheckModule(uri string) *SolaModule {
	hzOnce.Do(func() {
		if healthCheckModule == nil {
			healthCheckModule = NewModule(uri)
			attachControllers()
			setPreMiddlewares()
			setPostMiddlewares()
		}
	})

	return healthCheckModule
}

func attachControllers() {
	//* Attatching Controller Directly
	ctr := NewHealthCheckController()
	// ctr2, _ := NewAnotherController()
	//	...

	healthCheckModule.SetControllers(ctr)
}

func setPreMiddlewares() {
	healthCheckModule.SetPreMiddlewares(
		func(ctx *gin.Context) {
			log.Println("Health Checking...")
		},
	)
}

func setPostMiddlewares() {
	healthCheckModule.SetPostMiddlewares(
		func(ctx *gin.Context) {
			log.Println("Health Check Done!")
		},
	)
}
