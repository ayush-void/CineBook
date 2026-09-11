package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/yourusername/cinebook/backend/internal/models"
	"github.com/yourusername/cinebook/backend/internal/ws"
)

func main() {
	// --- Postgres ---
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=cinebook port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// Auto-migrate schema on boot
	if err := db.AutoMigrate(
		&models.User{},
		&models.Movie{},
		&models.Screen{},
		&models.Seat{},
		&models.Show{},
		&models.Booking{},
	); err != nil {
		log.Fatalf("failed to migrate schema: %v", err)
	}

	// --- Redis ---
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	// --- WebSocket hub for live seat updates ---
	hub := ws.NewHub()

	mux := http.NewServeMux()

	mux.HandleFunc("/ws/shows/", func(w http.ResponseWriter, r *http.Request) {
		showID := r.URL.Path[len("/ws/shows/"):]
		hub.HandleConnection(w, r, showID)
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("cinebook api listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
