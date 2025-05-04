package solanum_test

import (
	"testing"

	solanum "github.com/annuums/solanum"
)

// Compile-time checks that core types satisfy their interfaces.
var (
	_ solanum.Module     = (*solanum.SolaModule)(nil)
	_ solanum.Controller = (*solanum.SolaController)(nil)
	_ solanum.Runner     = *solanum.NewSolanum(5050)
)

// Dummy test to ensure this file is included in test suite.
func TestInterfaceCompliance(t *testing.T) {
	// This test passes if the code compiles successfully.
}
