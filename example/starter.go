package main

import "github.com/annuums/solanum"

func main() {
	server := *solanum.NewSolanum(5050)

	var healthCheckModule solanum.Module
	healthCheckUri := "/ping"
	healthCheckModule = solanum.NewHealthCheckModule(
		healthCheckUri,
	)

	server.SetModules(&healthCheckModule)

	server.Run()
}
