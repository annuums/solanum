package main

import (
	"github.com/annuums/solanum"
	"github.com/annuums/solanum/docs/examples/dependency_injection/web/mvc/user"
)

func main() {

	RegisterDependencies()

	app := *solanum.NewSolanum(5050)
	app.SetModules(
		user.NewModule(""),
	)

	app.Run()
}
