// Package tsserv is a simple HTTP server with endpoints that return some time series data.
package tsserv

import (
	"context"
	"fmt"
	"github.com/tinkermode/tsserv/pkg/logger"
	"net/http"
	"time"
)

type (
	Server struct {
		mux      *http.ServeMux
		httpServ *http.Server
	}
)

func sendErrorResponse(response http.ResponseWriter, status int, errMsg string) {
	response.WriteHeader(status)
	if _, err := response.Write([]byte(errMsg + "\n")); err != nil {
		logger.ErrorLogger.Printf("Failed to write body: %v", err)
	}
}

func New(port int) *Server {
	// may be more new variables here like read time out, write time out, max header bytes
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", sayHello)
	mux.HandleFunc("/data", getRawDataPoints)

	httpServ := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        mux,
		ErrorLog:       logger.ErrorLogger,
		ReadTimeout:    10 * time.Second, // todo @tangxin may be need to change to a config or env variable
		WriteTimeout:   10 * time.Second, // todo @tangxin may be need to change to a config or env variable
		MaxHeaderBytes: 1 << 20,
	}

	return &Server{mux: mux, httpServ: httpServ}
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.InfoLogger.Println("Shutting down server...")
	return s.httpServ.Shutdown(ctx)
}

func (s *Server) Run() error {
	return s.httpServ.ListenAndServe()
}
