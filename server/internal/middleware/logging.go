package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// statusRecorder wraps http.ResponseWriter so we can capture the status code
// the handler wrote (the standard ResponseWriter doesn't expose it).
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Logger logs one line per request: method, path, resulting status, how long it
// took, and the browser Origin header (handy when debugging CORS/preflights).
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Default to 200: if the handler never calls WriteHeader explicitly,
		// net/http sends 200, so that's the value we should report.
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"origin", r.Header.Get("Origin"),
			"duration", time.Since(start).String(),
		)
	})
}
