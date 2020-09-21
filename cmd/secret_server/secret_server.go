package main

import (
	"context"
	"github.com/pawmart/wp-atrd-task/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pawmart/wp-atrd-task/internal/http/app"
	"github.com/pawmart/wp-atrd-task/internal/http/router"
)

func main() {
	a := app.NewApp(storage.NewInMemoryStorage())
	r := router.NewRoutes(a)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":3001",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen  %v\n", err)
	}
}
