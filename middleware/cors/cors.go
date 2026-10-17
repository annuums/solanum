package cors

import "errors"

// AllowAllOrigins is the wildcard origin that permits every origin.
const AllowAllOrigins = "*"

// Default CORS settings
var (
	DefaultMethods = []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"}
	DefaultHeaders = []string{
		"Access-Control-Allow-Headers",
		"Origin",
		"Accept",
		"X-Requested-With",
		"Content-Type",
		"Access-Control-Request-Method",
		"Access-Control-Request-Headers",
	}
	DefaultCredentials = false
)

// Option defines configuration settings for Cross-Origin Resource Sharing (CORS).
type Option struct {
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
func WithUrls(urls []string) func(*Option) {

	return func(s *Option) {
		s.Urls = urls
	}
}

// WithHeaders sets the allowed HTTP headers for CORS.
func WithHeaders(headers []string) func(*Option) {

	return func(s *Option) {
		s.Headers = headers
	}
}

// WithMethods sets the allowed HTTP methods for CORS.
func WithMethods(methods []string) func(*Option) {

	return func(s *Option) {
		s.Methods = methods
	}
}

// WithAllowCredentials enables or disables transmission of credentials (cookies) in CORS requests.
func WithAllowCredentials(allowCredentials bool) func(*Option) {

	return func(s *Option) {
		s.AllowCredentials = allowCredentials
	}
}

// WithOriginFunc sets a custom origin validation function for CORS.
func WithOriginFunc(originFunc func(origin string) bool) func(*Option) {

	return func(s *Option) {
		s.OriginFunc = originFunc
	}
}

// WithMaxAge sets the maximum age (in hours) for CORS preflight requests to be cached.
func WithMaxAge(maxAge int) func(*Option) {

	return func(s *Option) {
		s.MaxAge = maxAge
	}
}

// Options applies a list of option functions to a CorsOption instance and
// fills in defaults for any missing settings.
// When neither urls nor originFunc is configured, all origins are allowed with the '*' wildcard
// (the request origin is never reflected).
func Options(opts ...func(*Option)) *Option {

	var options Option
	for _, opt := range opts {

		opt(&options)
	}

	// Use default headers if none specified
	if len(options.Headers) == 0 {

		options.Headers = DefaultHeaders
	}

	// Use default methods if none specified
	if len(options.Methods) == 0 {

		options.Methods = DefaultMethods
	}

	// Allow all origins only when neither urls nor originFunc is configured
	if options.OriginFunc == nil && len(options.Urls) == 0 {

		options.Urls = []string{AllowAllOrigins}
	}

	return &options
}

// AllowsAllOrigins reports whether the option permits every origin via the '*' wildcard.
func (o *Option) AllowsAllOrigins() bool {
	for _, url := range o.Urls {
		if url == AllowAllOrigins {

			return true
		}
	}

	return false
}

// Validate checks for insecure or conflicting CORS settings.
// Credentials cannot be combined with the '*' wildcard: browsers reject such responses,
// and reflecting arbitrary origins with credentials would expose authenticated endpoints.
func (o *Option) Validate() error {
	if o.AllowCredentials && o.AllowsAllOrigins() {

		return errors.New("cors: AllowCredentials cannot be used with all origins ('*'); specify urls or originFunc")
	}

	return nil
}
