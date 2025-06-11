package util

import (
	"log"
	"strings"
)

// Default CORS settings
var (
	CorsDefaultMethods = []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"}
	CorsDefaultHeaders = []string{
		"Access-Control-Allow-Headers",
		"Origin",
		"Accept",
		"X-Requested-With",
		"Content-Type",
		"Access-Control-Request-Method",
		"Access-Control-Request-Headers",
	}
	CorsDefaultCredentials  = false
	CorsDefaultOriginalFunc = func(origin string) bool {
		// Default origin function allows any localhost origin
		return strings.Contains(origin, "://localhost")
	}
)

// CorsOption defines configuration settings for Cross-Origin Resource Sharing (CORS).
type CorsOption struct {
	// Urls list of allowed origin URLs
	Urls []string

	// Headers list of allowed HTTP headers
	Headers []string

	// Methods list of allowed HTTP methods
	Methods []string

	// AllowCredentials whether cookies and credentials are allowed
	AllowCredentials bool

	// OriginFunc custom function to validate origin
	OriginFunc func(origin string) bool

	// MaxAge preflight cache duration in hours
	MaxAge int
}

// WithUrls sets the allowed origin URLs for CORS.
func WithUrls(urls []string) func(*CorsOption) {

	return func(s *CorsOption) {
		s.Urls = urls
	}
}

// WithHeaders sets the allowed HTTP headers for CORS.
func WithHeaders(headers []string) func(*CorsOption) {

	return func(s *CorsOption) {
		s.Headers = headers
	}
}

// WithMethods sets the allowed HTTP methods for CORS.
func WithMethods(methods []string) func(*CorsOption) {

	return func(s *CorsOption) {
		s.Methods = methods
	}
}

// WithAllowCredentials enables or disables transmission of credentials (cookies) in CORS requests.
func WithAllowCredentials(allowCredentials bool) func(*CorsOption) {

	return func(s *CorsOption) {
		s.AllowCredentials = allowCredentials
	}
}

// WithOriginFunc sets a custom origin validation function for CORS.
func WithOriginFunc(originFunc func(origin string) bool) func(*CorsOption) {

	return func(s *CorsOption) {
		s.OriginFunc = originFunc
	}
}

// WithMaxAge sets the maximum age (in hours) for CORS preflight requests to be cached.
func WithMaxAge(maxAge int) func(*CorsOption) {

	return func(s *CorsOption) {
		s.MaxAge = maxAge
	}
}

// CorsOptions applies a list of option functions to a CorsOption instance and
// fills in defaults for any missing settings.
func CorsOptions(opts ...func(*CorsOption)) *CorsOption {

	var options CorsOption
	for _, opt := range opts {

		opt(&options)
	}

	// Use default headers if none specified
	if options.Headers == nil || len(options.Headers) == 0 {

		options.Headers = CorsDefaultHeaders
	}

	// Use default methods if none specified
	if options.Methods == nil || len(options.Methods) == 0 {

		options.Methods = CorsDefaultMethods
	}

	// Use default origin validation function
	if options.OriginFunc == nil {

		// If no URLs provided, allow all origins
		if options.Urls == nil || len(options.Urls) == 0 {

			log.Println("Both urls and originfunc for cors are not defined. allowing all origins...")
			options.OriginFunc = func(origin string) bool {

				return true
			}
		} else {

			// Restrict to configured URLs or wildcard
			options.OriginFunc = func(origin string) bool {

				// '*' wildcard allows all origins
				if len(options.Urls) == 1 && options.Urls[0] == "*" {

					return true
				}

				// Permit origin if it matches one in the list
				for _, allowed := range options.Urls {

					if origin == allowed {
						return true
					}
				}

				return false
			}
		}
	}

	return &options
}
