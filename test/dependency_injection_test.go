package solanum

import (
	"database/sql"
	"github.com/annuums/solanum"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sync"
	"testing"
)

type providerEntry struct {
	// factory constructs a new instance of the provider.
	factory func() interface{}

	// singleton indicates whether to reuse the same instance across resolves.
	singleton bool

	// instance holds the created singleton instance after first resolution.
	instance interface{}

	// hookCalled tracks if the initHook has already been executed.
	hookCalled bool

	// initHook is an optional callback that runs once after creating the instance.
	initHook func(interface{})

	// interfaceType, if non-nil, registers this provider under a Go interface type.
	interfaceType reflect.Type

	// providerType is the concrete type returned by the factory or provided directly.
	providerType reflect.Type
}

type container struct {
	mu           sync.RWMutex              // protects all maps below
	providers    map[string]*providerEntry // key -> providerEntry
	interfaceMap map[reflect.Type]string   // interface type -> key
	typeMap      map[reflect.Type]string   // concrete type -> key
}

// globalContainer is the shared, package-level DI container instance.
var globalContainer = &container{
	providers:    make(map[string]*providerEntry),
	interfaceMap: make(map[reflect.Type]string),
	typeMap:      make(map[reflect.Type]string),
}

// resetContainer clears the global DI container state for test isolation.
func resetContainer() {
	globalContainer = &container{
		providers:    make(map[string]*providerEntry),
		interfaceMap: make(map[reflect.Type]string),
		typeMap:      make(map[reflect.Type]string),
	}
}

// TestValidateDependencies_OK ensures that ValidateDependencies passes when all deps are registered.
func TestValidateDependencies_OK(t *testing.T) {
	resetContainer()

	// Register a dummy provider for key "foo"
	solanum.Register("foo", func() int { return 123 }, solanum.WithSingleton())

	// Create a module that depends on "foo"
	mod := solanum.NewModule("/test")
	mod.SetDependencies(solanum.Dep[int]("foo"))

	// Setup runner with the module
	r := solanum.NewSolanum(0)
	runnerIface := *r // Runner interface
	runnerConcrete := runnerIface.(solanum.Runner)
	runnerConcrete.SetModules(mod)

	// Validate should pass (no error)
	err := runnerConcrete.ValidateDependencies()
	assert.NoError(t, err)
}

// TestValidateDependencies_FailureMissing reports an error when a dependency is not registered.
func TestValidateDependencies_FailureMissing(t *testing.T) {
	resetContainer()

	// Do not register any provider for "missing"

	mod := solanum.NewModule("/test")
	mod.SetDependencies(solanum.Dep[string]("missing"))

	r := solanum.NewSolanum(0)
	runnerIface := *r
	runnerConcrete := runnerIface.(solanum.Runner)
	runnerConcrete.SetModules(mod)

	err := runnerConcrete.ValidateDependencies()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dependency validation failed for key=\"missing\"")
}

// DummyService depends on *sql.DB for demonstration
type DummyService struct {
	db *sql.DB
}

func NewDummyService(db *sql.DB) *DummyService {
	return &DummyService{db: db}
}

// TestWithDepInjection verifies that WithDep option causes correct injection by key.
func TestWithDepInjection(t *testing.T) {
	resetContainer()

	// Register *sql.DB as singleton under key "db"
	solanum.Register(
		"db",
		func() *sql.DB {
			// Return a dummy *sql.DB (nil is ok for test)
			return &sql.DB{}
		},
		solanum.WithSingleton(),
	)

	// Register DummyService with transient scope and WithDep for "db"
	solanum.Register(
		"svc",
		func(ds *sql.DB) *DummyService {
			return NewDummyService(ds)
		},
		solanum.WithTransient(),
		solanum.WithDep[*sql.DB]("db"),
	)

	// Resolve the service
	inst, err := solanum.Resolve("svc")
	assert.NoError(t, err)

	svc, ok := inst.(*DummyService)
	assert.True(t, ok, "expected *DummyService type")
	assert.NotNil(t, svc.db, "db should have been injected")
}

// TestAutomaticTypeInjection verifies automatic injection when WithDep is not used.
func TestAutomaticTypeInjection(t *testing.T) {
	resetContainer()

	// Register *sql.DB under key "db"
	solanum.Register(
		"db",
		func() *sql.DB { return &sql.DB{} },
		solanum.WithSingleton(),
	)

	// Register service without WithDep, relying on reflect-based auto deps
	solanum.Register(
		"svc2",
		func(ds *sql.DB) *DummyService {
			return NewDummyService(ds)
		},
		solanum.WithTransient(),
	)

	inst, err := solanum.Resolve("svc2")
	assert.NoError(t, err)

	svc, ok := inst.(*DummyService)
	assert.True(t, ok)
	assert.NotNil(t, svc.db)
}
