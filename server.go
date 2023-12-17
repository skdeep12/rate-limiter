package main

import (
	"fmt"
	"io"
	"net/http"
	ratelimiter "rate-limiter/algorithms"
	"time"
)

func rateLimiter(next http.Handler) http.Handler {
	bucket := ratelimiter.NewLeakyBucket(10, 1)
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			isRequestAllowedByRateLimiter := bucket.ConsumeTokens(1, time.Now())
			if isRequestAllowedByRateLimiter {
				next.ServeHTTP(w, r)
			} else {
				fmt.Println("rate limit exceeded")
				w.WriteHeader(http.StatusTooManyRequests)
				io.WriteString(w, "rate limit exceeded.")
			}
		},
	)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("service hello world")
	io.WriteString(w, "Hello World")
}

func main() {
	http.Handle("/hello", rateLimiter(http.HandlerFunc(helloHandler)))
	fmt.Println("running server at :3333")
	_ = http.ListenAndServe(":3333", nil)

}
