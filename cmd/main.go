package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/wys1203/go-gorilla-example/users/delivery"
	"github.com/wys1203/go-gorilla-example/users/entity"
	"github.com/wys1203/go-gorilla-example/users/repository"
	"github.com/wys1203/go-gorilla-example/users/usecase"
)

func main() {
	dsn := "host=localhost user=ui_test password=mysecretpassword dbname=ui_test port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&entity.User{})

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := delivery.NewUserHandler(userUsecase)

	router := mux.NewRouter()
	delivery.RegisterUserRoutes(router, userHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Channel to listen for OS signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	// Run the server in a goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Block until an OS signal is received
	<-stopChan
	log.Println("Shutting down the server...")

	// Create a context with a timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
