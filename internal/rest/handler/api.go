package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/karthiknatarajan/plinth/internal/rest/middleware/audit"
	"github.com/karthiknatarajan/plinth/internal/rest/middleware/logging"
	"github.com/karthiknatarajan/plinth/internal/rest/middleware/nocache"
	"github.com/karthiknatarajan/plinth/internal/types"
	"github.com/rs/zerolog/hlog"
)

// APIHandler is an abstraction of an http handler that handles api calls.
type APIHandler interface {
	http.Handler
}

// NewAPIHandler returns a new APIHandler.
func NewAPIHandler(
	appCtx context.Context,
	config *types.Config,
) APIHandler {

	// Use go-chi router for inner routing.
	r := chi.NewRouter()

	// Apply common middleware.
	r.Use(nocache.NoCache)
	r.Use(middleware.Recoverer)

	// configure logging middleware.
	r.Use(hlog.URLHandler("http.url"))
	r.Use(hlog.MethodHandler("http.method"))
	r.Use(logging.HLogRequestIDHandler())
	r.Use(logging.HLogAccessLogHandler())

	// configure cors middleware
	r.Use(corsHandler(config))

	r.Use(audit.Middleware())

	r.Route("/v1", func(r chi.Router) {
	})

	return r
}

func corsHandler(config *types.Config) func(http.Handler) http.Handler {
	return cors.New(
		cors.Options{
			AllowedOrigins:   config.Cors.AllowedOrigins,
			AllowedMethods:   config.Cors.AllowedMethods,
			AllowedHeaders:   config.Cors.AllowedHeaders,
			ExposedHeaders:   config.Cors.ExposedHeaders,
			AllowCredentials: config.Cors.AllowCredentials,
			MaxAge:           config.Cors.MaxAge,
		},
	).Handler
}
