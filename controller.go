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

	// Handler is the Gin handler function to execute.
	Handler gin.HandlerFunc
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
