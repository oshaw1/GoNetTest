package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		log.Printf(
			"[%s] %s %s %s %d %v",
			start.Format("2006-01-02 15:04:05.000"),
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			rw.status,
			duration.Round(time.Microsecond),
		)
	}
}
