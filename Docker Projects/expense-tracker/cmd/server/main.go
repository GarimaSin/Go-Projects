package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "expense-tracker/internal/api"
    "expense-tracker/internal/auth"
    "expense-tracker/internal/db"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    cfg := loadConfigFromEnv()
    // Initialize DB pool
    pg, err := db.NewPostgresPool(cfg.PostgresURL, cfg.DBMaxConns)
    if err != nil {
        log.Fatalf("db connect: %v", err)
    }
    defer pg.Close()

    // Initialize Redis client (for rate-limiting, caching)
    redisClient := db.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword)
    defer redisClient.Close()

    authManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry)

    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    // Add Prometheus endpoint
    r.Handle("/metrics", promhttp.Handler())

    // Health
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    // API routes
    api.RegisterRoutes(r, pg, redisClient, authManager, cfg)

    srv := &http.Server{
        Addr:         cfg.BindAddr,
        Handler:      r,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    go func() {
        log.Printf("listening on %s", cfg.BindAddr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit
    log.Println("shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("server shutdown: %v", err)
    }
    log.Println("server stopped gracefully")
}

// Simple config loader (replace with env library or Viper)
type Config struct {
    BindAddr      string
    PostgresURL   string
    DBMaxConns    int32
    RedisAddr     string
    RedisPassword string
    JWTSecret     string
    JWTExpiry     time.Duration
}

func loadConfigFromEnv() *Config {
    return &Config{
        BindAddr:      getEnv("BIND_ADDR", ":8080"),
        PostgresURL:   getEnv("POSTGRES_URL", "postgres://postgres:password@localhost:5432/expenses?sslmode=disable"),
        DBMaxConns:    20,
        RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
        RedisPassword: getEnv("REDIS_PASSWORD", ""),
        JWTSecret:     getEnv("JWT_SECRET", "verysecret"),
        JWTExpiry:     24 * time.Hour,
    }
}

func getEnv(k, d string) string {
    if v := os.Getenv(k); v != "" {
        return v
    }
    return d
}
