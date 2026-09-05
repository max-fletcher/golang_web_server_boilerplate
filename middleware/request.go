package middleware

import "net/http"

func MaxBodySizeMiddleware(maxMegaBytes int) func(http.Handler) http.Handler { // Middleware for setting
	maxBytes := int64(maxMegaBytes) << 20 // 10 << 20 means 10 multiplied by 2 to the power of 20 so 10MB in bytes
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}
