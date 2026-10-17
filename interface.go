package solanum

import (
	"context"
	"net/http"

	"github.com/annuums/solanum/middleware/cors"
	"github.com/gin-gonic/gin"
)

type (

	// Controller declares a set of service handlers for a logical grouping of routes.
	// SolaController is the default implementation.
	Controller interface {
		// Handlers returns the slice of service definitions this controller manages.
		Handlers() []*SolaService
	}

	// Module represents a self-contained HTTP module mounted under its own URI prefix.
	// SolaModule is the default implementation with controllers and pre/post middlewares.
	Module interface {
		// Uri returns the base URI path for this module (e.g., "/users").
		Uri() string

		// SetRoutes binds the module's controllers, middleware, and routes onto a RouterGroup.
		SetRoutes(router *gin.RouterGroup)
	}

	// Runner is the application entrypoint interface for Solanum.
	// It manages module initialization, CORS, and server start.
	// Global middlewares are configured with NewSolanum options (WithMiddlewares, WithRecovery).
	Runner interface {
		// InitModules sets up all registered modules and their routes (idempotent).
		InitModules()

		// Handler initializes module routes (once) and returns an http.Handler for tests or custom servers.
		Handler() http.Handler

		// Modules returns a slice of pointers to registered Modules.
		Modules() []Module

		// SetModules replaces all registered Modules with the provided list.
		SetModules(m ...Module)

		// AddModules appends one or more Modules to the Runner.
		AddModules(m ...Module)

		// GinEngine exposes the underlying *gin.Engine for custom setup.
		GinEngine() *gin.Engine

		// Port exposes the configured port for the HTTP server.
		Port() int

		// Cors applies CORS configuration to the Gin engine using functional options.
		// It panics when the options are insecure (e.g. credentials with all origins)
		// or when called after routes are initialized.
		Cors(opts ...func(*cors.Option))

		// Run boots the HTTP server, initializing modules and listening on the configured port.
		// It shuts down gracefully on SIGINT/SIGTERM and returns an error when serving fails.
		Run() error

		// RunContext is like Run but shuts down gracefully when ctx is done.
		RunContext(ctx context.Context) error
	}
)
