package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"laboratorywork10/go-service/internal/auth"
	"laboratorywork10/go-service/internal/server"
)

func main() {
	if err := run(
		os.Getenv("JWT_SECRET"),
		os.Getenv("GO_SERVICE_PORT"),
		makeSignalChannel(),
	); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}

func run(secret, port string, stop <-chan os.Signal) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	if secret != "" {
		jwtSecret = secret
	}
	if jwtSecret == "" {
		jwtSecret = "super-secret-key"
	}

	serverPort := "8080"
	if port != "" {
		serverPort = port
	}

	authService := auth.NewService(jwtSecret, time.Hour)
	router := server.NewRouter(authService)

	httpServer := &http.Server{
		Addr:              ":" + serverPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("go-service is running on port %s", serverPort)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		return err
	case sig := <-stop:
		log.Printf("received signal %s, shutting down", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		return err
	}

	log.Print("go-service stopped gracefully")
	return nil
}

func makeSignalChannel() <-chan os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	return stop
}
