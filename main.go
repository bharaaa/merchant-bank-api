package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"merchant-bank-api/handlers"
	"merchant-bank-api/middleware"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/payment", middleware.Authenticate(handlers.Payment)).Methods("POST")
	router.HandleFunc("/logout", middleware.Authenticate(handlers.Logout)).Methods("POST")

	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on port 8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	<-quit // Block until we receive a signal
	log.Println("Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Println("Server exited gracefully")
}
