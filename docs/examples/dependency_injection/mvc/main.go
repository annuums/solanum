package main

import (
	"github.com/annuums/solanum"
	"github.com/annuums/solanum/docs/examples/dependency_injection/mvc/user"
)

func main() {

	RegisterDependencies()

	app := *solanum.NewSolanum(5050)
	app.SetModules(
		user.NewModule(""),
	)

	app.Run()
}
