package main

import "github.com/annuums/solanum"

func main() {
	server := *solanum.NewSolanum(5050)

	var healthCheckModule solanum.Module
	healthCheckUri := "/ping"
	healthCheckModule, _ = solanum.NewHealthCheckModule(
		server.GinEngine(),
		healthCheckUri,
	)

	server.SetModules(&healthCheckModule)

	server.Run()
}
