// cmd/main.go
package main

import (
	"flag"
	"fmt"
	"github.com/annuums/solanum"
	"github.com/annuums/solanum/docs/examples/dependency_injection/cli/user"
	"os"
)

func main() {

	RegisterDependencies()

	runner := solanum.NewSolanum()
	runner.SetModules(
		user.NewModule(),
	)

	//ctx := context.Background()
	//if repoInst, err := container.Resolve("userRepository"); err == nil {
	//
	//	ctx = context.WithValue(ctx, container.NewContextKey("userRepository"), repoInst)
	//}

	cmd := flag.String("cmd", "", "Execute: list, add")
	name := flag.String("name", "", "User Name (required in add command)")
	email := flag.String("email", "", "User Email (Required in add command)")

	flag.Parse()

	if *cmd == "" {

		printUsageAndExit()
	}

	switch *cmd {
	case "list":

		if err := user.ListUsersCLI(); err != nil {

			fmt.Fprintf(os.Stderr, "Error listing users: %v\n", err)
			os.Exit(1)
		}

	case "add":

		if *name == "" || *email == "" {

			fmt.Fprintln(os.Stderr, "Error: -name and -email are required for add command")
			printUsageAndExit()
		}

		if err := user.AddUserCLI(*name, *email); err != nil {

			fmt.Fprintf(os.Stderr, "Error adding user: %v\n", err)
			os.Exit(1)
		}

	default:

		fmt.Fprintf(os.Stderr, "Unknown cmd: %s\n", *cmd)
		printUsageAndExit()
	}
}

func printUsageAndExit() {

	fmt.Println(`Usage:
  -cmd=list
     → List All Users

  -cmd=add -name=<NAME> -email=<EMAIL>
     → Add New User

examples:
  go run main.go -cmd=list
  go run main.go -cmd=add -name="Bob" -email="bob@example.com"`)

	os.Exit(1)
}
