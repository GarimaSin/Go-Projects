package api

import (
    "encoding/json"
    "net/http"
    "strconv"
    "context"
    "time"

    "expense-tracker/internal/repository"
    "expense-tracker/internal/auth"
    "expense-tracker/internal/model"

    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/redis/go-redis/v9"
)

type AppConfig struct {
    ExpenseRepo *repository.ExpenseRepo
    AuthManager *auth.JWTManager
    Redis       *redis.Client
}

func RegisterRoutes(r chi.Router, db *pgxpool.Pool, redisClient *redis.Client, authManager *auth.JWTManager, cfg interface{}) {
    repo := repository.NewExpenseRepo(db)
    app := &AppConfig{ExpenseRepo: repo, AuthManager: authManager, Redis: redisClient}

    r.Route("/api/v1", func(r chi.Router) {
        r.Post("/auth/register", app.Register)
        r.Post("/auth/login", app.Login)

        r.Group(func(r chi.Router) {
            r.Use(AuthMiddleware(authManager))
            r.Use(RateLimitMiddleware(redisClient))
            r.Post("/expenses", app.CreateExpense)
            r.Get("/expenses/{id}", app.GetExpense)
            r.Get("/expenses", app.ListExpenses)
            r.Put("/expenses/{id}", app.UpdateExpense)
            r.Delete("/expenses/{id}", app.DeleteExpense)
        })
    })
}

// Simple register/login implementations (in-memory user store for demo)
// In production, use a users table and repository.

type registerReq struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type loginReq struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// naive in-memory user store for example
var demoUsers = map[string]int64{} // email -> id
var demoUserPasswords = map[string]string{} // email -> hashed (bcrypt) - for demo we'll store plaintext to keep simple
var nextUserID int64 = 1

func (a *AppConfig) Register(w http.ResponseWriter, r *http.Request) {
    var req registerReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    if req.Email == "" || req.Password == "" {
        http.Error(w, "email/password required", http.StatusBadRequest)
        return
    }
    if _, exists := demoUsers[req.Email]; exists {
        http.Error(w, "user exists", http.StatusConflict)
        return
    }
    // NOTE: replace with bcrypt hashing in production
    demoUsers[req.Email] = nextUserID
    demoUserPasswords[req.Email] = req.Password
    nextUserID++
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(`{"status":"ok"}`))
}

func (a *AppConfig) Login(w http.ResponseWriter, r *http.Request) {
    var req loginReq
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    id, exists := demoUsers[req.Email]
    if !exists || demoUserPasswords[req.Email] != req.Password {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }
    // generate token
    token, err := a.AuthManager.Generate(id)
    if err != nil {
        http.Error(w, "failed to generate token", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Example handler: CreateExpense
func (a *AppConfig) CreateExpense(w http.ResponseWriter, r *http.Request) {
    var req model.Expense
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    // Extract user ID from context (set by Auth middleware)
    uid, ok := r.Context().Value("user_id").(int64)
    if !ok {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }
    req.UserID = uid
    if req.OccurredAt.IsZero() {
        req.OccurredAt = time.Now().UTC()
    }

    ctx := r.Context()
    if err := a.ExpenseRepo.Create(ctx, &req); err != nil {
        http.Error(w, "failed to create", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(req)
}

func (a *AppConfig) GetExpense(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }
    ctx := r.Context()
    e, err := a.ExpenseRepo.GetByID(ctx, id)
    if err != nil {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(e)
}

func (a *AppConfig) ListExpenses(w http.ResponseWriter, r *http.Request) {
    // For brevity: not implemented; return empty list
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode([]model.Expense{})
}

func (a *AppConfig) UpdateExpense(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (a *AppConfig) DeleteExpense(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "not implemented", http.StatusNotImplemented)
}
