package main

import (
	"context"
	"errors"
	"flag"
	"github.com/tinkermode/tsserv/pkg/logger"
	"github.com/tinkermode/tsserv/pkg/tsserv"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	port := flag.Int("p", 8080, "port number of server")
	flag.Parse()

	logger.InfoLogger.Printf("Starting server on port %d\n", *port)

	s := tsserv.New(*port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.InfoLogger.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			logger.ErrorLogger.Fatalf("Server forced to shutdown: %v", err)
		}
		logger.InfoLogger.Println("Server gracefully stopped")
	}()

	if err := s.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.ErrorLogger.Fatalf("Server exited with error: %v", err)
	}
}
