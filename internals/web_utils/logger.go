package webutils

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf(
			"[%s] %s %s from %s - User-Agent: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			r.Header.Get("User-Agent"),
		)
		next.ServeHTTP(w, r)
		log.Printf(
			"Request for %s %s completed in %s",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}
