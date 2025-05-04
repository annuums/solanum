package user

import (
	"github.com/annuums/solanum"
	"sync"
)

var (
	module *solanum.SolaModule
	once   sync.Once
)

func NewModule(uri string) *solanum.SolaModule {
	once.Do(func() {
		if module == nil {
			module = solanum.NewModule(uri)
			attachControllers()
		}
	})

	return module
}

func attachControllers() {
	ctr := NewController()

	module.SetControllers(ctr)
	module.SetDependencies(
		solanum.Dep[UserRepository]("userRepository"),
	)
}
