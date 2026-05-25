package solanum

import "github.com/gin-gonic/gin"

// SolaController groups one or more SolaService handlers under a logical controller.
// It implements the Controller interface.
type SolaController struct {
	handlers []*SolaService
}

// SolaService represents a single HTTP route handler configuration.
type SolaService struct {
	// Uri is the route path relative to the module prefix (e.g., "/:id").
	Uri string

	// Method is the HTTP method to bind (GET, POST, etc.).
	Method string

	// PreHandlers are middleware functions to run before Handler for this route only.
	PreHandlers []gin.HandlerFunc

	// Handler is the Gin handler function to execute.
	Handler gin.HandlerFunc

	// PostHandlers are middleware functions to run after Handler for this route only.
	// Gin executes all handlers in the chain sequentially; ctx.Next() is not required.
	// Use ctx.Abort() in any handler to stop further execution.
	PostHandlers []gin.HandlerFunc
}

// NewController constructs an empty SolaController ready to receive handlers.
func NewController() *SolaController {

	return &SolaController{
		handlers: make([]*SolaService, 0),
	}
}

// SetHandlers replaces the controller's handler list with the provided entries.
func (ctr *SolaController) SetHandlers(handlers ...*SolaService) {

	ctr.handlers = handlers
}

// Handlers returns the slice of SolaService entries managed by this controller.
func (ctr *SolaController) Handlers() []*SolaService {

	return ctr.handlers
}
