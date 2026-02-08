# PO Backend

Go REST API backend for **PO**, a Tron-themed social network. Built with Echo v5, GORM, and PostgreSQL.

## Tech Stack

- **Go 1.25** with [Echo v5](https://echo.labstack.com/) HTTP framework
- **GORM** ORM with PostgreSQL driver
- **Redis** for caching
- **JWT** (HS256, 7-day expiry) for authentication
- **WebSocket** for real-time notifications
- **bcrypt** for password hashing

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL
- Redis

### Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your database credentials:

```
SERVER_PORT=8000
DB_HOST=localhost
DB_PORT=5432
DB_NAME=po
DB_USERNAME=postgres
DB_PASSWORD=postgres
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
JWT_SECRET=your-secret-key
```

3. Run the server (auto-migrates database tables on startup):

```bash
go run cmd/server.go
```

4. Seed sample data (optional):

```bash
go run seeder/seeder.go

# Reset and re-seed
go run seeder/seeder.go -reset
```

## Commands

```bash
go run cmd/server.go          # Start the server
go run seeder/seeder.go       # Seed database
go vet ./...                  # Vet
gofmt -l -w .                 # Format
golangci-lint run --fix       # Lint (used by pre-commit hook via lefthook)
```

## Architecture

Layered architecture with dependency injection wired in `routes/routes.go`:

```
Controllers -> Services -> Repositories -> GORM/PostgreSQL
```

```
po-backend/
├── cmd/            # Application entry point
├── configs/        # Database, Redis, and WebSocket configuration
├── controllers/    # HTTP request handlers
├── dtos/           # Request/response structs with validation tags
├── helper/         # JWT and password utilities
├── middlewares/     # Auth and ownership middleware
├── models/         # GORM entities
├── repositories/   # Data access layer
├── routes/         # Dependency injection and route registration
├── seeder/         # Database seeder
├── services/       # Business logic layer
├── utilities/      # Redis, WebSocket, and content parsing helpers
└── validation/     # Request validation
```

## API Endpoints

All routes are prefixed with `/api/v1`.

### Auth & Users

| Method | Endpoint                       | Auth | Description               |
|--------|--------------------------------|------|---------------------------|
| POST   | `/users`                       | No   | Register a new user       |
| POST   | `/login`                       | No   | Login and receive JWT     |
| GET    | `/verify`                      | Yes  | Verify token / get user   |
| GET    | `/users`                       | No   | List all users            |
| GET    | `/users/:id`                   | No   | Get user by ID            |
| GET    | `/users/username/:username`    | No   | Get user by username      |
| GET    | `/search?q=`                   | No   | Search users by name      |

### Follow

| Method | Endpoint           | Auth | Description                        |
|--------|--------------------|------|------------------------------------|
| GET    | `/following/users` | Yes  | Get users the current user follows |
| POST   | `/follow/:id`      | Yes  | Follow a user                      |
| DELETE | `/unfollow/:id`    | Yes  | Unfollow a user                    |

### Posts

| Method | Endpoint                    | Auth  | Description                 |
|--------|-----------------------------|-------|-----------------------------|
| GET    | `/content/posts`            | No    | Get all posts               |
| GET    | `/content/posts/:id`        | No    | Get post by ID              |
| POST   | `/content/posts`            | Yes   | Create post                 |
| PUT    | `/content/posts/:id`        | Owner | Update post                 |
| DELETE | `/content/posts/:id`        | Owner | Delete post                 |
| GET    | `/content/following/posts`  | Yes   | Get posts from followed users|

### Comments

| Method | Endpoint                  | Auth  | Description     |
|--------|---------------------------|-------|-----------------|
| POST   | `/content/comments`       | Yes   | Create comment  |
| PUT    | `/content/comments/:id`   | Owner | Update comment  |
| DELETE | `/content/comments/:id`   | Owner | Delete comment  |

### Likes

| Method | Endpoint                      | Auth | Description          |
|--------|-------------------------------|------|----------------------|
| POST   | `/content/like/posts/:id`     | Yes  | Like a post          |
| DELETE | `/content/unlike/posts/:id`   | Yes  | Unlike a post        |
| POST   | `/content/like/comments/:id`  | Yes  | Like a comment       |
| DELETE | `/content/unlike/comments/:id`| Yes  | Unlike a comment     |
| GET    | `/content/likes/posts/:id`    | No   | Get post likers      |
| GET    | `/content/likes/comments/:id` | No   | Get comment likers   |

### Reposts

| Method | Endpoint                        | Auth | Description      |
|--------|---------------------------------|------|------------------|
| POST   | `/content/repost/posts/:id`     | Yes  | Repost a post    |
| DELETE | `/content/unrepost/posts/:id`   | Yes  | Unrepost a post  |

### Hashtags

| Method | Endpoint                          | Auth | Description                  |
|--------|-----------------------------------|------|------------------------------|
| GET    | `/content/hashtags/:tag/posts`    | Yes  | Get posts by hashtag         |
| GET    | `/content/hashtags/trending`      | Yes  | Get trending hashtags        |

### Bookmarks

| Method | Endpoint                  | Auth | Description        |
|--------|---------------------------|------|--------------------|
| GET    | `/content/bookmarks`      | Yes  | Get user bookmarks |
| POST   | `/content/bookmarks/:id`  | Yes  | Bookmark a post    |
| DELETE | `/content/bookmarks/:id`  | Yes  | Remove bookmark    |

### Notifications

| Method | Endpoint                  | Auth | Description              |
|--------|---------------------------|------|--------------------------|
| GET    | `/content/notis`          | Yes  | Get notifications        |
| PUT    | `/content/notis/read`     | Yes  | Mark all as read         |
| PUT    | `/content/notis/read/:id` | Yes  | Mark one as read         |

### WebSocket

| Endpoint          | Description                                          |
|-------------------|------------------------------------------------------|
| `/ws/subscribe`   | Real-time notifications (send JWT as first message)  |

## Features

- **Posts** with CRUD, likes, comments, bookmarks, and reposts
- **Hashtags** — `#tag` in post content is auto-parsed and stored; supports hashtag search and trending
- **Mentions** — `@username` in posts and comments triggers notifications to mentioned users; the `/following/users` endpoint powers client-side mention autocomplete
- **Real-time notifications** via WebSocket for likes, comments, follows, reposts, and mentions
- **Follow system** with follower/following feeds
- **Ownership middleware** ensures only authors can edit/delete their content

## CORS

Configured for `localhost:5173` and `localhost:3000` by default.
