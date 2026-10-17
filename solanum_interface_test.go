package solanum_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	solanum "github.com/annuums/solanum"
)

// Compile-time checks that core types satisfy their interfaces.
var (
	_ solanum.Module     = (*solanum.SolaModule)(nil)
	_ solanum.Controller = (*solanum.SolaController)(nil)
	_ solanum.Runner     = (*solanum.Server)(nil)
)

// TestInterfaceCompliance passes if the code compiles successfully.
func TestInterfaceCompliance(t *testing.T) {}

// healthModule is a minimal custom Module implementing only Uri and SetRoutes.
type healthModule struct{}

func (healthModule) Uri() string { return "/health" }

func (healthModule) SetRoutes(router *gin.RouterGroup) {

	router.GET("", func(c *gin.Context) { c.Status(http.StatusNoContent) })
}

// TestCustomModule verifies a hand-written Module plugs into the server.
func TestCustomModule(t *testing.T) {

	server := solanum.NewSolanum()
	server.AddModules(healthModule{})

	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	assert.Equal(t, http.StatusNoContent, rec.Code)
}
