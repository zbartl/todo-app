package middleware

import (
	"github.com/gobridge-kr/todo-app/server/utils"
	"net/http"
)

type RequireAuthMiddleware struct {
	next http.Handler
	jwt *jwtea.Provider
}

func (mw *RequireAuthMiddleware) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	CorsMiddleware{Next: NewJwtMiddleware(mw.next, mw.jwt)}.ServeHTTP(w, r)
}

func RequiredAuth(next http.Handler, jwt *jwtea.Provider) *JwtMiddleware {
	return &JwtMiddleware{
		next: next,
		jwt: jwt,
	}
}
