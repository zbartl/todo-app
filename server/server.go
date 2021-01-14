package server

import (
	"github.com/gobridge-kr/todo-app/server/config"
	"github.com/gobridge-kr/todo-app/server/controller"
	"github.com/gobridge-kr/todo-app/server/database"
	"github.com/gobridge-kr/todo-app/server/handler"
	"github.com/gobridge-kr/todo-app/server/middleware"
	"github.com/zbartl/jwtea"
	"net/http"
)

// Server represents current server status
type Server struct {
	baseURL     string
	middlewares []func(w http.ResponseWriter, r *http.Request)
}

// New creates a new Server with given URL
func New(baseURL string) *Server {
	return &Server{
		baseURL: baseURL,
	}
}

// Middleware configures middleware to process requests
func (s *Server) Middleware(middleware func(w http.ResponseWriter, r *http.Request)) {
	s.middlewares = append(s.middlewares, middleware)
}

func (s *Server) ConfigureRoutes(mux *http.ServeMux,
	database *database.Database,
	jwt *jwtea.Provider,
	thirdPartyAuthConfig *config.ThirdPartyAuthConfiguration,
) {
	mux.Handle(
		"/todo",
		jwtea.RequiredJwtAuthentication(controller.Todo(database), jwt),
	)
	
	mux.Handle(
		"/user/auth",
		http.StripPrefix(
			"/user",
			&middleware.AllowAnonymous{Next: handler.Auth(jwt, thirdPartyAuthConfig)},
		),
	)
}

// Serve starts the actual serving job
func (s *Server) Serve(mux *http.ServeMux, port string) {
	http.ListenAndServe(":"+port, mux)
}
