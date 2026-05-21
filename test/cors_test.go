package solanum_test

import (
	"github.com/annuums/solanum/middleware/cors"

	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCorsOptionsDefaultValues validates defaults when no options are provided.
func TestCorsOptionsDefaultValues(t *testing.T) {
	op := cors.Options()

	// Headers and Methods should pick up global defaults
	assert.Equal(t, cors.DefaultHeaders, op.Headers)
	assert.Equal(t, cors.DefaultMethods, op.Methods)

	// OriginFunc should allow all when no URLs set
	assert.True(t, op.OriginFunc("https://random.com"))
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
