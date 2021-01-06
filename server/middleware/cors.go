package middleware

import "net/http"

type CorsMiddleware struct {
	Next http.Handler
}

// Cors middleware enables cross-origin requests
func (mw CorsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("access-control-allow-origin", "*")
	w.Header().Set("access-control-allow-methods", "GET, POST, PATCH, DELETE")
	w.Header().Set("access-control-allow-headers", "accept, content-type")
	w.Header().Set("content-type", "application/json; charset=UTF-8")
	mw.Next.ServeHTTP(w, r)
}
