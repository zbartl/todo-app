package middleware

import (
	"github.com/gobridge-kr/todo-app/server/utils"
	"net/http"
)

type JwtMiddleware struct {
	next http.Handler
	jwt *jwtea.Provider
}

func (mw JwtMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	
	err := mw.jwt.Validate(token)
	if err != nil {
		mw.next.ServeHTTP(w, r)
		return
	}
	
	http.Error(w, err.Error(), http.StatusUnauthorized)
	return
}

func NewJwtMiddleware(next http.Handler, jwt *jwtea.Provider) *JwtMiddleware {
	return &JwtMiddleware{
		next: next,
		jwt: jwt,
	}
}
