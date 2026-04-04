package middleware

import "net/http"

// defaultMaxBodySize is the default maximum request body size (10 MB).
const defaultMaxBodySize int64 = 10 << 20

// MaxBodySize returns middleware that limits the size of incoming request bodies.
// It wraps r.Body with http.MaxBytesReader, which returns an error if the body
// exceeds maxBytes. If maxBytes is 0, the default limit of 10 MB is used.
func MaxBodySize(maxBytes int64) func(http.Handler) http.Handler {
	if maxBytes <= 0 {
		maxBytes = defaultMaxBodySize
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}
