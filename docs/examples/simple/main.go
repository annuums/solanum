package main

import (
	"github.com/annuums/solanum"
	"github.com/annuums/solanum/http/health"
)

func main() {
	server := *solanum.NewSolanum(5050)

	healthCheckUri := "/ping"
	healthCheckModule := health.NewHealthCheckModule(
		healthCheckUri,
	)

	server.SetModules(healthCheckModule)

	server.Run()
}
