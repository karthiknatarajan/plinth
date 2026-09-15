package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	DefaultReadHeaderTimeout = 2 * time.Second
	DefaultReadTimeout       = 5 * time.Second
	DefaultWriteTimeout      = 15 * time.Second
	DefaultIdleTimeout       = 0
)

// Config defines the config of an http server.
type Config struct {
	Port              int
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

// Server is a wrapper around http.Server that exposes different async ListenAndServe methods
// that return corresponding ShutdownFunctions.
type Server struct {
	config  Config
	handler http.Handler
}

// ShutdownFunction defines a function that is called to shutdown the server.
type ShutdownFunction func(context.Context) error

func NewServer(config Config, handler http.Handler) *Server {
	if config.ReadHeaderTimeout == 0 {
		config.ReadHeaderTimeout = DefaultReadHeaderTimeout
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = DefaultReadTimeout
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = DefaultWriteTimeout
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = DefaultIdleTimeout
	}

	return &Server{
		config:  config,
		handler: handler,
	}
}

// ListenAndServe initializes a server to respond to HTTP network requests.
func (s *Server) ListenAndServe() (*errgroup.Group, ShutdownFunction) {
	return s.listenAndServe()
}

func (s *Server) listenAndServe() (*errgroup.Group, ShutdownFunction) {
	var g errgroup.Group
	s1 := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.config.Port),
		ReadHeaderTimeout: s.config.ReadHeaderTimeout,
		ReadTimeout:       s.config.ReadTimeout,
		WriteTimeout:      s.config.WriteTimeout,
		IdleTimeout:       s.config.IdleTimeout,
		Handler:           s.handler,
	}
	g.Go(func() error {
		return s1.ListenAndServe()
	})

	return &g, s1.Shutdown
}
