package server

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/karthiknatarajan/plinth/internal/rest/server"
	"github.com/karthiknatarajan/plinth/internal/utils/config"
	"github.com/karthiknatarajan/plinth/internal/utils/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type system struct {
	logger zerolog.Logger
	server *server.Server
}

func Run() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger.SetupLogger(cfg)

	system, err := provideSystem(ctx, cfg)
	if err != nil {
		panic(err)
	}

	// gCtx is canceled if any of flowwing occurs: any go routine launched with g encrounters error or ctx is canceled
	g, gCtx := errgroup.WithContext(ctx)

	// e.g. for starting a background go routine
	// g.Go(func() error {
	// 	return background.Start(ctx)
	// })

	// start server
	gHTTP, shutdownHTTP := system.server.ListenAndServe()
	g.Go(gHTTP.Wait)

	log.Info().
		Int("port", cfg.Server.HTTP.Port).
		Msg("server started")

	// wait until the error group context is done
	<-gCtx.Done()

	// restore default behavior on interrept signal
	stop()
	log.Info().Msg("shutting down gracefully (press Ctrl+C again to force)")

	// shutdown servers gracefully
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTime)
	defer cancel()

	if sErr := shutdownHTTP(shutdownCtx); sErr != nil {
		log.Error().Err(sErr).Msg("failed to shutdown http server")
	}

	log.Info().Msg("Waiting for subroutines to complete")
	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
	return
}
