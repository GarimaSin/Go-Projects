package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/prometheus/client_golang/prometheus/promhttp"

    "github.com/example/urlshort/internal/api"
    "github.com/example/urlshort/internal/db"
    "github.com/example/urlshort/internal/events"
    "github.com/example/urlshort/internal/short"
)

type Config struct {
    BindAddr    string
    PostgresURL string
    RedisAddr   string
    RedisPassword string
    KafkaBrokers []string
    DBMaxConns  int32
}

func loadConfigFromEnv() Config {
    return Config{
        BindAddr:    ":8080",
        PostgresURL: "postgres://postgres:password@localhost:5432/urlshort",
        RedisAddr:   "localhost:6379",
        KafkaBrokers: []string{"localhost:9092"},
        DBMaxConns:  10,
    }
}

func main() {
    cfg := loadConfigFromEnv()

    pg, err := db.NewPostgresPool(cfg.PostgresURL, cfg.DBMaxConns)
    if err != nil {
        log.Fatalf("pg connect: %v", err)
    }
    defer pg.Close()

    redisClient := db.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword)
    defer redisClient.Close()

    producer, err := events.NewKafkaProducer(cfg.KafkaBrokers)
    if err != nil {
        log.Fatalf("kafka producer: %v", err)
    }
    defer producer.Close()

    generator := short.NewGenerator(pg)
    svc := short.NewService(pg, redisClient, producer, generator)

    r := chi.NewRouter()
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Handle("/metrics", promhttp.Handler())

    api.RegisterRoutes(r, svc)

    srv := &http.Server{
        Addr:         cfg.BindAddr,
        Handler:      r,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 5 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        log.Printf("listening on %s", cfg.BindAddr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit
    log.Println("shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("server shutdown: %v", err)
    }
}
