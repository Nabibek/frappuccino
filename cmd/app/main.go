package main

import (
	"database/sql"
	"fmt"
	"frappuccino/internal/api"
	"frappuccino/internal/api/handlers"
	"frappuccino/internal/repo"
	"frappuccino/internal/service"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// logger — простая middleware для логирования запросов
func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}

	repo := repo.New(db)
	svc := service.New(repo)
	handler := handlers.New(svc)

	mux := api.Router(handler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", logger(mux)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
