# Expense Tracker API (Go)

This is a runnable scaffold project for a scalable Go-based expense tracker API.
It includes:
- HTTP server (chi)
- Postgres connection (pgxpool)
- Redis client for rate-limiting/caching
- JWT auth manager
- Basic handlers for register/login and expenses (create/get)
- Dockerfile and docker-compose for local development
- Kubernetes deployment + HPA snippets

To run locally with docker-compose:
1. docker-compose up --build
2. Apply DB migration: connect to Postgres and run migrations/001_create_tables.sql
3. Use /api/v1/auth/register and /api/v1/auth/login to create a demo user and get a token.
4. Use the token in Authorization: Bearer <token> to call /api/v1/expenses endpoints.
