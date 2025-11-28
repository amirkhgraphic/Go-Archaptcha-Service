# go-arcaptcha-service

Small Gin + Gorm service that demonstrates arcaptcha-like validation and basic CRUD/list/group operations.

## Quick start
1) Copy `.env.example` to `.env` and adjust (defaults: `DB_PATH=data/data.db`, `PORT=8080`).
2) Run Database migration: `go run migrations/migration_001.go`.
3) Seed test data (Optional): `go run seeds/seed_users.go` will insert a handful of demo users (idempotent).
4) Start the API: `go run main.go` (listens on `:8080`).

### Docker Compose
```bash
docker compose up --build
```
Environment defaults come from `.env` (copy from `.env.example`).


## Endpoints
- `GET /ping` - health check.
- `GET /__fake/arcaptcha/challenge` - mint a one-time `challenge_id`.
- `POST /__fake/arcaptcha/verify` - check a token without consuming it.
- `POST /api/users` - create user (requires `challenge_id`).
- `GET /api/users` - list users with `page`, `page_size`, `sort`, `search`, `username`, `email`.
- `GET /api/users/:id` - fetch a user.
- `PATCH /api/users/:id` - update user (requires `challenge_id`).
- `GET /api/users/group` - aggregate users by gender/nationality (e.g., `?group_by=gender,nationality`).

## Captcha simulation rules
- Omit or empty `challenge_id` -> 400.
- Unknown/expired `challenge_id` -> 400.
- `challenge_id` ending with `-neterr` -> 503 to mimic network failure.
- One-time use: `ValidateChallenge` consumes the token. The verify endpoint only peeks and keeps it usable.

Typical flow:
```bash
curl -s http://localhost:8080/__fake/arcaptcha/challenge
# -> copy challenge_id

curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"alice\",\"email\":\"alice@example.com\",\"bio\":\"demo\",\"challenge_id\":\"<challenge_id>\"}"
```

## Grouping endpoint
`GET /api/users/group` accepts `group_by` combinations of `gender` and `nationality`, and returns counts per group.

## Postman collection
Import `docs/postman_collection.json` and set `base_url` (default `http://localhost:8080`) and `challenge_id` variables. Use the fake arcaptcha endpoints to refresh tokens for protected requests.
