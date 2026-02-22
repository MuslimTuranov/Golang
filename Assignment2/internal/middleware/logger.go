package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Printf("ts=%s method=%s path=%s dur=%s",
				start.Format(time.RFC3339), r.Method, r.URL.Path, time.Since(start))
		})
	}
}
