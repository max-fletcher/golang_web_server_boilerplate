package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
	"github.com/max-fletcher/golang_web_server_boilerplate/helpers/responses"
)

// All chi middleware always returns a "func(http.Handler) http.Handler" i.e it returns a function that accepts a "http.Handler" and returns a
// "http.Handler" in turn
func RateLimiter(requestsCount int, windowMins int) func(http.Handler) http.Handler {
	return httprate.LimitBy(
		requestsCount,                         // number of requests
		time.Duration(windowMins)*time.Minute, // time window in mins
		httprate.KeyByRealIP,
		// The third parameter and onwards accepts multiple option functions(httprate.Limit is variadic funtion), which are handler functions
		// that have many variations. e.g "httprate.WithLimitHandler" is a custom handler function that you can add more logic
		// to(if you need to). Here, it sends back a custom json response. But if you want to choose which key to use to rate-limit,
		// you can use "httprate.WithKeyFuncs" where you can define what keys should be used to rate-limit requests like this:
		// httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
		// 	return r.Header.Get("X-User-ID"), nil
		// }),
		// You can also pass multiple of these; each option function controlling an aspect of the rate-limiter e.g:
		// httprate.Limit(
		// 	100,
		// 	time.Minute,
		// 	httprate.WithKeyFuncs(httprate.KeyByIP),
		// 	httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
		// 		respondWithJSON(w, 429, map[string]string{
		// 			"error": "rate limit exceeded",
		// 		})
		// 	}),
		// ),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			responses.RespondWithJSON(w, http.StatusTooManyRequests, responses.Response{
				Code:    http.StatusTooManyRequests, // 429 status code
				Status:  "error",
				Message: "Too many requests. Rate limit exceeded.",
			})
		}),
	)
}
