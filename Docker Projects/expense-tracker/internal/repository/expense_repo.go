package repository

import (
    "context"
    "time"

    "expense-tracker/internal/model"

    "github.com/jackc/pgx/v5/pgxpool"
)

type ExpenseRepo struct {
    db *pgxpool.Pool
}

func NewExpenseRepo(db *pgxpool.Pool) *ExpenseRepo {
    return &ExpenseRepo{db: db}
}

func (r *ExpenseRepo) Create(ctx context.Context, e *model.Expense) error {
    q := `INSERT INTO expenses (user_id, title, amount, category, occurred_at, created_at, updated_at)
          VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`
    now := time.Now().UTC()
    e.CreatedAt = now
    e.UpdatedAt = now
    row := r.db.QueryRow(ctx, q, e.UserID, e.Title, e.Amount, e.Category, e.OccurredAt, e.CreatedAt, e.UpdatedAt)
    return row.Scan(&e.ID)
}

func (r *ExpenseRepo) GetByID(ctx context.Context, id int64) (*model.Expense, error) {
    q := `SELECT id, user_id, title, amount, category, occurred_at, created_at, updated_at FROM expenses WHERE id=$1`
    var e model.Expense
    if err := r.db.QueryRow(ctx, q, id).Scan(&e.ID, &e.UserID, &e.Title, &e.Amount, &e.Category, &e.OccurredAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
        return nil, err
    }
    return &e, nil
}
