package db

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/redis/go-redis/v9"
)

func NewPostgresPool(connStr string, maxConns int32) (*pgxpool.Pool, error) {
    cfg, err := pgxpool.ParseConfig(connStr)
    if err != nil {
        return nil, err
    }
    cfg.MaxConns = maxConns
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    pool, err := pgxpool.NewWithConfig(ctx, cfg)
    if err != nil {
        return nil, err
    }
    return pool, nil
}

func NewRedisClient(addr, password string) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       0,
    })
}
