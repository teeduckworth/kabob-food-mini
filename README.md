# KabobFood Backend

Golang backend for the Telegram mini app food delivery service. The service exposes REST APIs for the client mini app and the admin panel, applies Redis caching for menu/regions, enforces idempotent order creation, and integrates with Telegram Bot API for notifications.

## Stack

- Go 1.22 + Gin
- PostgreSQL + golang-migrate migrations
- Redis for caching/idempotency
- Docker & docker-compose for local setup
- Zap for JSON logging, JWT for auth scaffolding

## Project Layout

```
cmd/app            # service entrypoint
internal/app       # application wiring (router, server)
internal/config    # env-driven configuration
internal/http      # HTTP router + handlers
internal/server    # HTTP server lifecycle
internal/db        # PostgreSQL helpers
internal/cache     # Redis helpers
migrations         # SQL migrations (golang-migrate compatible)
```

## Quick Start (Docker)

```bash
# build and run everything
docker compose up --build

# apply migrations only
docker compose run --rm migrate
```

The API will be available on `http://localhost:8080`. Health endpoints: `/healthz`, `/version`. Prometheus metrics exposed on `/metrics`.

Initial REST surface (full spec in `docs/openapi.yaml` / Postman collection in `docs/postman_collection.json`):

- `POST /auth/telegram` — обмен initData на JWT + профиль
- `GET /menu`, `GET /regions` — публичное меню и зоны
- `GET /profile` — профиль + адреса (требует `Authorization: Bearer <token>`)
- `GET/POST/PUT/DELETE /addresses` — CRUD адресов пользователя
- `POST /orders` — создание заказа с `client_request_id` (идемпотентно)
- `GET /orders`, `GET /orders/:id` — список и детали заказов пользователя
- `POST /admin/login` — авторизация оператора (логин/пароль → JWT с ролью admin)
- `POST/PUT/DELETE /admin/categories/:id` — управление категориями меню (требует admin JWT)
- `POST/PUT/DELETE /admin/products/:id` — управление товарами
- `POST/PUT/DELETE /admin/regions/:id` — управление зонами доставки
- `GET /admin/orders`, `PUT /admin/orders/:id/status` — лента заказов (query: `status`, `from`, `to`, `limit`, `offset`) и смена статусов оператором

## Running locally (without Docker)

```bash
export APP_ENV=local
export HTTP_PORT=8080
export DB_URL=postgres://postgres:postgres@localhost:5432/kabobfood?sslmode=disable
export REDIS_URL=redis://localhost:6379/0
export JWT_SECRET=supersecret

# run migrations (requires migrate CLI installed)
migrate -path=./migrations -database=$DB_URL up

# start the app
go run ./cmd/app
```

## Configuration

| Variable | Description |
| --- | --- |
| `APP_ENV` | `development`, `docker`, `prod`, affects Gin mode/logging |
| `HTTP_HOST` / `HTTP_PORT` | HTTP server bind address |
| `DB_URL` | PostgreSQL DSN |
| `REDIS_URL` | Redis URI (e.g. `redis://redis:6379/0`) |
| `JWT_SECRET` | JWT signing secret |
| `TELEGRAM_BOT_TOKEN` / `TELEGRAM_ADMIN_CHAT_ID` | Telegram Bot API credentials для уведомлений пользователей/операторов |
| `AUTH_TELEGRAM_INIT_TTL` | TTL for Telegram initData validation (default 1h) |
| `CACHE_MENU_TTL` / `CACHE_REGIONS_TTL` | Redis TTLs for menu/regions caches |
| `ADMIN_DEFAULT_USERNAME` / `ADMIN_DEFAULT_PASSWORD` | bootstrap admin credentials created при старте |
| `ADMIN_JWT_EXPIRATION` | TTL админского JWT |
| `RATE_USER_LIMIT` / `RATE_ADMIN_LIMIT` / `RATE_WINDOW` | rate limit (requests per window) для user/admin API |
| `SENTRY_DSN` | Sentry DSN for error reporting |
| `SHUTDOWN_TIMEOUT` | Graceful shutdown timeout (default 10s) |

## API Documentation

- `docs/openapi.yaml` — OpenAPI 3.0 спецификация (можно импортировать в Swagger UI/Stoplight)
- `docs/postman_collection.json` — готовая коллекция Postman (подставьте `base_url`, `token`, `admin_token`)
- `mini-app` — Next.js клиент (README внутри каталога)

## Deployment & HTTPS

- Используйте reverse-proxy (Nginx/Traefik/Cloudflare) с TLS offload'ом перед сервисом (`HTTP_PORT=8080`).
- Выставляйте переменные `RATE_*` для ограничения RPS и подключайте firewall/Cloudflare для L7 защиты.
- Для VPS: описать systemd unit, использовать `docker compose` или `make deploy`. Для PaaS (Railway/Render/Fly.io) достаточно задать env и подключить volume для PG/Redis.
- Добавьте `SENTRY_DSN` и Prometheus scrape job (`/metrics`) в проде.

## Next Steps

- Добавить rate limiting/HTTPS конфигурацию и RBAC для admin/user
- Подготовить OpenAPI спецификацию и Postman коллекцию per ТЗ
- Написать unit-тесты бизнес-логики (заказы, кеш, статусы)
- Расширить админский API (фильтры заказов, поиск, отчёты)
- Документировать сценарии деплоя (VPS/Railway/Render) и настройку мониторинга
