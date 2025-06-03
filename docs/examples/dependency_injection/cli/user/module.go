package user

import (
	"github.com/annuums/solanum"
	"github.com/annuums/solanum/container"
	"sync"
)

var (
	module *solanum.SolaModule
	once   sync.Once
)

func NewModule() *solanum.SolaModule {
	once.Do(func() {
		if module == nil {
			module = solanum.NewModule(
				solanum.WithDependency(
					container.DepConfig[UserRepository]("userRepository"),
				),
			)
		}
	})

	return module
}
