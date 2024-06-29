package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"supertal-tha-parking-app/logger"
	"time"
)

// Server is the basic structure of the server
type Server struct {
	name    string
	port    int
	timeout time.Duration
	cleanup func()
	handler http.Handler
}

// NewServer creates a new server
func NewServer(name string, port int, timeout time.Duration, h http.Handler) *Server {
	return &Server{
		name:    name,
		port:    port,
		timeout: timeout,
		handler: h,
	}
}

// Run runs the server
func (svr *Server) Run() {
	logger.GetLogger().Infof("starting %s server...", svr.name)

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", svr.port),
		Handler:           svr.handler,
		ReadTimeout:       svr.timeout,
		ReadHeaderTimeout: svr.timeout,
		WriteTimeout:      svr.timeout,
		IdleTimeout:       svr.timeout,
	}

	go func() {
		logger.GetLogger().Infof("%s server listening on port %d", svr.name, svr.port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger().Errorf("%s server stopped listening: %v", svr.name, server.ListenAndServe())
		}
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop

	logger.GetLogger().Infof("%s server shutdown initiated...", svr.name)
	ctx, cancel := context.WithTimeout(context.Background(), svr.timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.GetLogger().Errorf("%s server shutdown error: %v", svr.name, err)
	}
	logger.GetLogger().Infof("%s server shutdown complete", svr.name)
}
