package middlerware

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[INFO] %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
