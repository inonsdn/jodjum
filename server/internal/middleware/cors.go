package middleware

import "net/http"

// CORS returns a middleware that tells browsers which origins may call this API.
//
// It takes the allowed origin as an argument (instead of hardcoding it) so the
// same code works for http://localhost:5173 in dev and your real frontend URL
// in production — the value comes from config/env.
func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// These headers are what the browser reads to decide whether to
			// allow the frontend's JavaScript to see the response.
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			// Authorization MUST be here because you send "Bearer <token>";
			// Content-Type is needed for JSON request bodies.
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			// Let the browser cache the preflight result for 1 hour so it
			// doesn't send an OPTIONS request before every single call.
			w.Header().Set("Access-Control-Max-Age", "3600")

			// The preflight: the browser sends OPTIONS (with no auth) to ask
			// "am I allowed?". We've set the headers above, so just answer 204
			// and stop — do NOT pass it down to auth, which would 401 it.
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// A real request: hand it to the router (and then your auth).
			next.ServeHTTP(w, r)
		})
	}
}
