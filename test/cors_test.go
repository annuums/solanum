package solanum_test

import (
	"github.com/annuums/solanum/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCorsOptionsDefaultValues validates defaults when no options are provided.
func TestCorsOptionsDefaultValues(t *testing.T) {
	op := util.CorsOptions()

	// Headers and Methods should pick up global defaults
	assert.Equal(t, util.CorsDefaultHeaders, op.Headers)
	assert.Equal(t, util.CorsDefaultMethods, op.Methods)

	// OriginFunc should allow all when no URLs set
	assert.True(t, op.OriginFunc("https://random.com"))
}

// TestCorsOptionsCustomSettings verifies each functional option overrides defaults.
func TestCorsOptionsCustomSettings(t *testing.T) {
	urls := []string{"https://a.com", "https://b.com"}
	headers := []string{"X-Test"}
	methods := []string{"PATCH"}

	op := util.CorsOptions(
		util.WithUrls(urls),
		util.WithHeaders(headers),
		util.WithMethods(methods),
		util.WithAllowCredentials(true),
		util.WithOriginFunc(func(origin string) bool { return origin == "ok" }),
		util.WithMaxAge(5),
	)

	assert.Equal(t, urls, op.Urls)
	assert.Equal(t, headers, op.Headers)
	assert.Equal(t, methods, op.Methods)
	assert.True(t, op.AllowCredentials)
	assert.False(t, op.OriginFunc("nope"))
	assert.True(t, op.OriginFunc("ok"))
	assert.Equal(t, 5, op.MaxAge)
}
