package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	readR   = 5 * time.Second   // max time to read request from the client
	writeR  = 12 * time.Second  // max time to write response to the client
	keepA   = 120 * time.Second // max time for connections using TCP Keep-Alive
	timeout = 10 * time.Second  // max time to complete tasks before shutdown
)

//SERVE starts the application
func (c *chiRouter) SERVE() {
	port := getPort()

	s := http.Server{
		Addr:         port,
		Handler:      c.mux,
		ReadTimeout:  readR,
		WriteTimeout: writeR,
		IdleTimeout:  keepA,
	}

	log.Printf("Starting http server on port %s\n", port)
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error starting server %v", err)
		}
	}()
	log.Print("Server Started")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Print("Signal closing server received")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Println("Server shutdown failed", "error", err)
	}
	log.Println("Server shutdown gracefully")
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("environment variable PORT is required")
	}
	return ":" + port
}
