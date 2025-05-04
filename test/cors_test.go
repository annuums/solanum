package solanum_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	solanum "github.com/annuums/solanum"
)

// TestCorsOptionsDefaultValues validates defaults when no options are provided.
func TestCorsOptionsDefaultValues(t *testing.T) {
	op := solanum.CorsOptions()

	// Headers and Methods should pick up global defaults
	assert.Equal(t, solanum.CorsDefaultHeaders, op.Headers)
	assert.Equal(t, solanum.CorsDefaultMethods, op.Methods)

	// OriginFunc should allow all when no URLs set
	assert.True(t, op.OriginFunc("https://random.com"))
}

// TestCorsOptionsCustomSettings verifies each functional option overrides defaults.
func TestCorsOptionsCustomSettings(t *testing.T) {
	urls := []string{"https://a.com", "https://b.com"}
	headers := []string{"X-Test"}
	methods := []string{"PATCH"}

	op := solanum.CorsOptions(
		solanum.WithUrls(urls),
		solanum.WithHeaders(headers),
		solanum.WithMethods(methods),
		solanum.WithAllowCredentials(true),
		solanum.WithOriginFunc(func(origin string) bool { return origin == "ok" }),
		solanum.WithMaxAge(5),
	)

	assert.Equal(t, urls, op.Urls)
	assert.Equal(t, headers, op.Headers)
	assert.Equal(t, methods, op.Methods)
	assert.True(t, op.AllowCredentials)
	assert.False(t, op.OriginFunc("nope"))
	assert.True(t, op.OriginFunc("ok"))
	assert.Equal(t, 5, op.MaxAge)
}
