package solanum

type CorsOption struct {
	Urls             []string
	Headers          []string
	Methods          []string
	AllowCredentials bool
	OriginFunc       func(origin string) bool
	MaxAge           int
}

func WithUrl(urls []string) func(*CorsOption) {
	return func(s *CorsOption) {
		s.Urls = urls
	}
}

func WithHeaders(headers []string) func(*CorsOption) {
	return func(s *CorsOption) {
		s.Headers = headers
	}
}

func WithMethods(methods []string) func(*CorsOption) {
	return func(s *CorsOption) {
		s.Methods = methods
	}
}

func WithAllowCredentials(allowCredentials bool) func(*CorsOption) {
	return func(s *CorsOption) {
		s.AllowCredentials = allowCredentials
	}
}

func WithOriginFunc(originFunc func(origin string) bool) func(*CorsOption) {
	return func(s *CorsOption) {
		s.OriginFunc = originFunc
	}
}

func WithMaxAge(maxAge int) func(*CorsOption) {
	return func(s *CorsOption) {
		s.MaxAge = maxAge
	}
}

func CorsOptions(opts ...func(*CorsOption)) *CorsOption {
	var options CorsOption
	for _, opt := range opts {
		opt(&options)
	}

	if options.Urls == nil || len(options.Urls) == 0 {
		options.Urls = []string{"*"}
	}

	if options.Headers == nil || len(options.Headers) == 0 {
		options.Headers = CorsDefaultHeaders
	}

	if options.Methods == nil || len(options.Methods) == 0 {
		options.Methods = CorsDefaultMethods
	}

	if options.OriginFunc == nil {
		options.OriginFunc = func(origin string) bool {
			if len(options.Urls) == 1 && options.Urls[0] == "*" {
				return true
			}

			for _, allowedOrigin := range options.Urls {
				if origin == allowedOrigin {
					return true
				}
			}
			return false
		}
	}

	return &options
}
