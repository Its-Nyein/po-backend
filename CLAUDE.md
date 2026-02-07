# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go backend for a social media platform (Yaycha v2) built with Echo v5, GORM, and PostgreSQL. Provides REST APIs for users, posts, comments, likes, follows, and real-time notifications via WebSocket.

## Common Commands

```bash
# Run the server (loads .env, auto-migrates DB, starts on SERVER_PORT)
go run cmd/server.go

# Seed database with sample data
go run seeder/seeder.go

# Reset and re-seed database
go run seeder/seeder.go -reset

# Lint (used by pre-commit hook via lefthook)
golangci-lint run --fix

# Vet
go vet ./...

# Format
gofmt -l -w .
```

## Architecture

Layered architecture with dependency injection wired in `routes/routes.go`:

```
Controllers → Services → Repositories → GORM/PostgreSQL
```

- **models/**: GORM entities with relationship tags (User, Post, Comment, PostLike, CommentLike, Follow, Notification)
- **dtos/**: Request/response structs with `validate` tags (go-playground/validator)
- **repositories/**: Database queries using GORM; each entity has its own repository
- **services/**: Business logic wrapping repositories; handles password hashing, JWT generation, notification creation
- **controllers/**: Echo HTTP handlers; bind request, call service, return JSON
- **middlewares/**: `IsAuthenticated` (JWT from Authorization header → sets userID in context), `IsPostOwner`, `IsCommentOwner`
- **routes/**: Single `InitRoutes()` function wires all dependencies and registers routes under `/api/v1`

## Key Patterns

**Authentication flow**: Register → bcrypt hash (cost 14) → Login → JWT (HS256, 7-day expiry) → Bearer token in Authorization header. Helper functions in `helper/jwt.go` and `helper/password.go`.

**Notification system**: Controllers create notifications via `NotificationService` and broadcast to connected clients via `utilities/websocket.go` → `configs/websocket.go` (WSManager). WebSocket clients authenticate by sending JWT as their first message on `/api/v1/ws/subscribe`.

**Ownership middleware**: `IsPostOwner` and `IsCommentOwner` are middleware factories that take a service and verify the authenticated user owns the resource before allowing DELETE operations.

**Eager loading**: Repositories use GORM `Preload` extensively to load relationships (e.g., posts come with User, Comments, PostLikes).

## Configuration

Environment variables loaded from `.env` (see `.env.example`): SERVER_PORT, DB_HOST/PORT/NAME/USERNAME/PASSWORD, REDIS_HOST/PORT/PASSWORD/DB, JWT_SECRET. Auto-migration runs on startup in `configs/config.go`.

CORS is configured for `localhost:5173` and `localhost:3000`.

## API Route Structure

All routes under `/api/v1`. Content routes nested under `/api/v1/content`. Auth-required routes use `IsAuthenticated` middleware. Delete routes additionally use ownership middlewares.
