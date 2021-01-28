package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}
	log.Printf("Starting up on http://localhost:%s", port)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok."))
	})

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL is empty!")
	}
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("db open: %+v", err)
	}
	tr := &todosResource{
		db: db,
	}
	r.Mount("/todos", tr.Routes())

	log.Fatal(http.ListenAndServe(":"+port, r))
}

// JSON encodes the object and writes it to the response.
func JSON(w http.ResponseWriter, obj interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}

// JSONError sends a 500 with the error.
func JSONError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	JSON(w, map[string]string{"error": fmt.Sprintf("%v", err)})
}
