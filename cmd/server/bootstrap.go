package server

import (
	"context"

	"github.com/karthiknatarajan/plinth/internal/rest/handler"
	"github.com/karthiknatarajan/plinth/internal/rest/router"
	"github.com/karthiknatarajan/plinth/internal/rest/server"
	"github.com/karthiknatarajan/plinth/internal/types"
	"github.com/rs/zerolog"
)

func provideSystem(ctx context.Context, cfg *types.Config) (*system, error) {
	apiHandler := handler.NewAPIHandler(ctx, cfg)
	router := router.NewRouter(apiHandler)
	server := server.NewServer(cfg, router)

	return &system{
		server: server,
		logger: zerolog.Logger{},
	}, nil
}
