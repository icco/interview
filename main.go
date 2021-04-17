package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/icco/gutil/logging"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

var (
	service = "interview"
	project = "icco-cloud"
	log     = logging.Must(logging.NewLogger(service))
)

func main() {
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}
	log.Infow("Starting up", "host", fmt.Sprintf("http://localhost:%s", port))

	zgl := zapgorm2.New(log.Desugar())
	zgl.SetAsDefault()
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
		Logger: zgl,
	})
	if err != nil {
		log.Fatalw("cannot connect to database server", zap.Error(err))
	}

	if err := db.AutoMigrate(&ToDo{}); err != nil {
		log.Fatalw("cannot migrate todo", zap.Error(err))
	}

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(logging.Middleware(log.Desugar(), project))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok."))
	})

	tr := &todosResource{db: db}
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
	JSON(w, map[string]string{"error": fmt.Sprintf("%v", err)})
	w.WriteHeader(500)
}
