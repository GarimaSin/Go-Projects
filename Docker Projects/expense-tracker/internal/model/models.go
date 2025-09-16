package model

import "time"

type User struct {
    ID        int64     `json:"id"`
    Email     string    `json:"email"`
    Password  string    `json:"-"` // hashed
    CreatedAt time.Time `json:"created_at"`
}

type Expense struct {
    ID         int64     `json:"id"`
    UserID     int64     `json:"user_id"`
    Title      string    `json:"title"`
    Amount     float64   `json:"amount"`
    Category   string    `json:"category"`
    OccurredAt time.Time `json:"occurred_at"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
