// Package tsserv is a simple HTTP server with endpoints that return some time series data.
package tsserv

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	errLogger = log.New(os.Stderr, "[TSServ] ", log.LstdFlags)
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
		errLogger.Printf("Failed to write body: %v", err)
	}
}

func New(port int) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", sayHello)
	mux.HandleFunc("/data", getRawDataPoints)

	httpServ := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        mux,
		ErrorLog:       errLogger,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &Server{mux: mux, httpServ: httpServ}
}

func (s *Server) Run() error {
	return s.httpServ.ListenAndServe()
}
