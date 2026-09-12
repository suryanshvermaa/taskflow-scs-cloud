package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/scscloud/taskflow/internal/config"
	"github.com/scscloud/taskflow/internal/database"
	"github.com/scscloud/taskflow/internal/handlers"
)

func main() {
	cfg := config.Load()

	// ── Database ──────────────────────────────────────────────────────────
	db, err := database.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()
	log.Printf("database: opened %s", cfg.DBPath)

	// ── Handlers ──────────────────────────────────────────────────────────
	h := handlers.New(db)

	// ── Router ────────────────────────────────────────────────────────────
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/api/health", allowMethods(h.HealthCheck, http.MethodGet))

	// Tasks collection
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetTasks(w, r)
		case http.MethodPost:
			h.CreateTask(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Single task
	mux.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			h.UpdateTask(w, r)
		case http.MethodDelete:
			h.DeleteTask(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// ── Middleware stack ───────────────────────────────────────────────────
	handler := corsMiddleware(cfg.CORSOrigin, mux)
	handler = loggingMiddleware(handler)

	// ── Server ────────────────────────────────────────────────────────────
	addr := "0.0.0.0:" + cfg.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server: listening on http://%s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-quit
	log.Println("server: shutting down…")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server: forced shutdown: %v", err)
	}
	log.Println("server: stopped")
}

// corsMiddleware sets CORS headers for the configured origin.
func corsMiddleware(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs every request with its method, path and duration.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// allowMethods is a simple helper for single-method endpoints.
func allowMethods(fn http.HandlerFunc, methods ...string) http.HandlerFunc {
	allowed := make(map[string]bool, len(methods))
	for _, m := range methods {
		allowed[m] = true
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if !allowed[r.Method] {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		fn(w, r)
	}
}

