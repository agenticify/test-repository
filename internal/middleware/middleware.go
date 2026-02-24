package middleware

import (
	"log"
	"net/http"
	runtime/debug"
	"time"
)

// Logger logs the details of each request.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("%s %s %s %s", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
		}()
		next.ServeHTTP(w, r)
	})
}

// Recovery recovers from panics and logs the error.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Printf("PANIC: %s
%s", err, debug.Stack())
			}
		}()
		next.ServeHTTP(w, r)
	})
}
