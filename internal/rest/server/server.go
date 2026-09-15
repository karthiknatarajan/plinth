package server

import (
	"github.com/karthiknatarajan/plinth/internal/rest/router"
	"github.com/karthiknatarajan/plinth/internal/types"
	"github.com/karthiknatarajan/plinth/internal/utils/http"
)

type Server struct {
	*http.Server
}

func NewServer(cfg *types.Config, router *router.Router) *Server {
	return &Server{
		http.NewServer(
			http.Config{
				Port:              cfg.Server.HTTP.Port,
				ReadHeaderTimeout: cfg.Server.HTTP.ReadHeaderTimeout,
				ReadTimeout:       cfg.Server.HTTP.ReadTimeout,
				WriteTimeout:      cfg.Server.HTTP.WriteTimeout,
				IdleTimeout:       cfg.Server.HTTP.IdleTimeout,
			},
			router,
		),
	}
}
