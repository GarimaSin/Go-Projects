package api

import (
    "context"
    "net/http"
    "strconv"
    "time"

    "expense-tracker/internal/auth"

    "github.com/redis/go-redis/v9"
    "github.com/go-chi/chi/v5/middleware"
)

func AuthMiddleware(j *auth.JWTManager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := r.Header.Get("Authorization")
            if token == "" {
                http.Error(w, "missing auth", http.StatusUnauthorized)
                return
            }
            if len(token) > 7 && token[:7] == "Bearer " {
                token = token[7:]
            }
            claims, err := j.Verify(token)
            if err != nil {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }
            ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func RateLimitMiddleware(redisClient *redis.Client) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := r.Context().Value("user_id")
            key := "rl:"
            if user != nil {
                key += "u:" + strconv.FormatInt(user.(int64), 10)
            } else {
                key += "anon"
            }
            limit := 100
            ttl := time.Minute
            cnt, err := redisClient.Incr(r.Context(), key).Result()
            if err == nil {
                if cnt == 1 {
                    redisClient.Expire(r.Context(), key, ttl)
                }
                if cnt > int64(limit) {
                    http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
                    return
                }
            }
            next.ServeHTTP(w, r)
        })
    }
}
