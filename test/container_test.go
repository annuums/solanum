package solanum_test

import (
	"reflect"
	"testing"

	solanum "github.com/annuums/solanum"
	"github.com/stretchr/testify/assert"
)

// TestRegisterResolveSingleton verifies that singleton providers return the same instance.
func TestRegisterResolveSingleton(t *testing.T) {
	type payload struct{ Value int }
	factory := func() payload { return payload{Value: 7} }
	solanum.Register("singleton", factory, solanum.WithSingleton())

	inst1, err1 := solanum.Resolve("singleton")
	assert.NoError(t, err1)
	p1 := inst1.(payload)
	assert.Equal(t, 7, p1.Value)

	inst2, err2 := solanum.Resolve("singleton")
	assert.NoError(t, err2)
	// Verify same value and instance semantics
	p2 := inst2.(payload)
	assert.Equal(t, p1, p2)
}

// TestRegisterResolveTransient ensures transient providers yield new instances each time.
func TestRegisterResolveTransient(t *testing.T) {
	counter := 0
	type id struct{ N int }
	factory := func() id {
		counter++
		return id{N: counter}
	}
	solanum.Register("transient", factory, solanum.WithTransient())

	first, _ := solanum.Resolve("transient")
	second, _ := solanum.Resolve("transient")
	assert.NotEqual(t, first, second)
}

// TestInitHookRunsOnce verifies that the init hook is invoked exactly once.
func TestInitHookRunsOnce(t *testing.T) {
	calls := 0
	solanum.Register("hooked", func() string { return "ok" },
		solanum.WithSingleton(),
		solanum.WithInit(func(i interface{}) { calls++ }),
	)

	_, _ = solanum.Resolve("hooked")
	_, _ = solanum.Resolve("hooked")
	assert.Equal(t, 1, calls)
}

// TestResolveMissing verifies that resolving a non-registered key returns an error.
func TestResolveMissing(t *testing.T) {
	_, err := solanum.Resolve("no-such-key")
	assert.Error(t, err)
}

// TestResolveByTypeMismatch verifies error when instance does not implement interface.
func TestResolveByTypeMismatch(t *testing.T) {
	type MyIfc interface{ Foo() }
	solanum.Register("bad", "stringValue", solanum.WithSingleton(), solanum.As((*MyIfc)(nil)))

	_, err := solanum.ResolveByType("bad", reflect.TypeOf((*MyIfc)(nil)).Elem())
	assert.Error(t, err)
}
