package middleware

import "net/http"

type AllowAnonymous struct {
	Next http.Handler
}

func (mw *AllowAnonymous) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	CorsMiddleware{Next: mw.Next}.ServeHTTP(w, r)
}
