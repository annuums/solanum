package cors_test

import (
	"testing"

	"github.com/annuums/solanum/middleware/cors"
	"github.com/stretchr/testify/assert"
)

// TestCorsOptionsDefaultValues validates defaults when no options are provided.
func TestCorsOptionsDefaultValues(t *testing.T) {
	op := cors.Options()

	assert.Equal(t, cors.DefaultHeaders, op.Headers)
	assert.Equal(t, cors.DefaultMethods, op.Methods)

	// All origins are allowed with the wildcard (not by reflecting the origin) when no URLs set
	assert.Nil(t, op.OriginFunc)
	assert.Equal(t, []string{cors.AllowAllOrigins}, op.Urls)
	assert.True(t, op.AllowsAllOrigins())
	assert.NoError(t, op.Validate())
}

// TestCorsOptionsOriginFuncOnly keeps urls empty when only originFunc is configured.
func TestCorsOptionsOriginFuncOnly(t *testing.T) {

	op := cors.Options(
		cors.WithOriginFunc(func(origin string) bool { return origin == "https://ok.com" }),
	)

	assert.Empty(t, op.Urls)
	assert.False(t, op.AllowsAllOrigins())
}

// TestCorsOptionsValidateCredentials rejects credentials combined with all origins.
func TestCorsOptionsValidateCredentials(t *testing.T) {

	assert.Error(t, cors.Options(cors.WithAllowCredentials(true)).Validate())
	assert.Error(t, cors.Options(
		cors.WithUrls([]string{"https://a.com", "*"}),
		cors.WithAllowCredentials(true),
	).Validate())
	assert.NoError(t, cors.Options(
		cors.WithUrls([]string{"https://a.com"}),
		cors.WithAllowCredentials(true),
	).Validate())
}

// TestCorsOptionsCustomSettings verifies each functional option overrides defaults.
func TestCorsOptionsCustomSettings(t *testing.T) {
	urls := []string{"https://a.com", "https://b.com"}
	headers := []string{"X-Test"}
	methods := []string{"PATCH"}

	op := cors.Options(
		cors.WithUrls(urls),
		cors.WithHeaders(headers),
		cors.WithMethods(methods),
		cors.WithAllowCredentials(true),
		cors.WithOriginFunc(func(origin string) bool { return origin == "ok" }),
		cors.WithMaxAge(5),
	)

	assert.Equal(t, urls, op.Urls)
	assert.Equal(t, headers, op.Headers)
	assert.Equal(t, methods, op.Methods)
	assert.True(t, op.AllowCredentials)
	assert.False(t, op.OriginFunc("nope"))
	assert.True(t, op.OriginFunc("ok"))
	assert.Equal(t, 5, op.MaxAge)
}
