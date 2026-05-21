package solanum_test

import (
	"testing"

	solanum "github.com/annuums/solanum"
)

// Compile-time checks that core types satisfy their interfaces.
var (
	_ solanum.Module     = (*solanum.SolaModule)(nil)
	_ solanum.Controller = (*solanum.SolaController)(nil)
	_ solanum.Runner     = solanum.NewSolanum(
		solanum.WithPort(5050),
	)
)

// TestInterfaceCompliance passes if the code compiles successfully.
func TestInterfaceCompliance(t *testing.T) {}
