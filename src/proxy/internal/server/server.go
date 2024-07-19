package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port     int
	basePath string
	router   *chi.Mux
}

func NewServer() *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Invalid PORT: %v", err)
	}

	basePath := os.Getenv("BUCKET_BASE_PATH")
	if basePath == "" {
		log.Fatal("BUCKET_BASE_PATH environment variable is not set")
	}

	s := &Server{
		port:     port,
		basePath: basePath,
		router:   chi.NewRouter(),
	}

	s.setupRoutes()

	log.Println("Reverse proxy running on", port)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}
