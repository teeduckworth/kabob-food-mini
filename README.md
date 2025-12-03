# KabobFood Monorepo

Этот репозиторий содержит бэкенд (Go) и клиентское Mini App (Next.js) для Telegram.

## Backend (Go + Gin)

Golang backend обслуживает Telegram mini app и админ-панель: REST API, кеширование в Redis, идемпотентные заказы, интеграция с Telegram Bot API.

### Стек
- Go 1.22 + Gin
- PostgreSQL + golang-migrate migrations
- Redis для кеша и идемпотентности
- Docker / docker-compose для локальной разработки
- Zap для логов, JWT для авторизации

### Структура
```
cmd/app            # entrypoint сервиса
internal/app       # связывание компонентов
internal/config    # загрузка env-конфига
internal/http      # маршруты и handlers
internal/server    # lifecycle HTTP-сервера
internal/db        # PostgreSQL helper'ы
internal/cache     # Redis helper'ы
migrations         # SQL миграции (golang-migrate)
```

### Быстрый старт (Docker)
```bash
docker compose up --build          # сервис + БД + Redis
# только миграции
docker compose run --rm migrate
```
API по умолчанию доступен на `http://localhost:8080` (`/healthz`, `/version`, `/metrics`).

### Основные эндпоинты
- `POST /auth/telegram` — обмен initData на JWT + профиль
- `GET /menu`, `GET /regions` — публичные справочники
- `GET /profile` — профиль и адреса (нужен Bearer JWT)
- CRUD ` /addresses`
- `POST /orders` — создание заказа (идемпотентность по `client_request_id`)
- Админские `POST /admin/login`, `POST/PUT/DELETE /admin/categories|products|regions`, `GET/PUT /admin/orders`

### Запуск без Docker
```bash
export APP_ENV=local
export HTTP_PORT=8080
export DB_URL=postgres://postgres:postgres@localhost:5432/kabobfood?sslmode=disable
export REDIS_URL=redis://localhost:6379/0
export JWT_SECRET=supersecret

migrate -path=./migrations -database=$DB_URL up
GO111MODULE=on go run ./cmd/app
```

### Переменные окружения
| Var | Описание |
| --- | --- |
| `APP_ENV` | режим (`development`, `docker`, `prod`) |
| `HTTP_HOST` / `HTTP_PORT` | bind адрес HTTP сервера |
| `DB_URL` | DSN PostgreSQL |
| `REDIS_URL` | URI Redis |
| `JWT_SECRET` | ключ подписи JWT |
| `TELEGRAM_BOT_TOKEN` / `TELEGRAM_ADMIN_CHAT_ID` | интеграция с Telegram Bot API |
| `AUTH_TELEGRAM_INIT_TTL` | TTL для initData |
| `CACHE_MENU_TTL` / `CACHE_REGIONS_TTL` | TTL кешей |
| `ADMIN_DEFAULT_USERNAME` / `ADMIN_DEFAULT_PASSWORD` | bootstrap админ |
| `ADMIN_JWT_EXPIRATION` | TTL админского JWT |
| `RATE_USER_LIMIT` / `RATE_ADMIN_LIMIT` / `RATE_WINDOW` | лимиты RPS |
| `SENTRY_DSN` | DSN для Sentry |
| `SHUTDOWN_TIMEOUT` | graceful shutdown |

### Документация
- `docs/openapi.yaml` — OpenAPI 3.0
- `docs/postman_collection.json` — Postman коллекция

### Продакшн советы
- ставьте reverse‑proxy (Nginx/Traefik) с TLS
- конфигурируйте rate-лимиты и firewall/Cloudflare
- следите за `/metrics`, подключайте Sentry
- задокументируйте деплой на VPS/Railway/Render, мониторинг, бэкапы

## Mini App (Next.js + Tailwind)

Клиент для Telegram Mini App: меню, корзина, оформление заказа, профиль/адреса.

### Запуск
```bash
cd mini-app
npm install
npm run dev
```

### Env (`mini-app/.env.local`)
```
NEXT_PUBLIC_API_URL=http://localhost:8080
# для локальной отладки без Telegram
NEXT_PUBLIC_DEV_TOKEN=your-local-jwt
NEXT_PUBLIC_DEV_INIT_DATA=query_string_from_telegram
```

### Структура
- `src/app/page.tsx` — меню и карточки
- `src/app/checkout/page.tsx` — оформление заказа
- `src/app/profile/page.tsx` — профиль и адреса
- `src/components/*` — UI (табы, карточки, bottom-nav, корзина)
- `src/store/cart.ts` — Zustand хранилище
- `src/lib/api.ts` — клиент к backend REST API

### TODO
- Деплой (Vercel/Pages) и Telegram Mini App manifest
- Полноценное редактирование профиля, отображение статуса заказа
- UI tests/storybook

### Telegram WebApp интеграция
- Используем `@twa-dev/sdk`: `WebApp.ready()`, `expand()`, реагируем на `themeChanged`
- BackButton внутри Telegram синхронизирован с роутингом (главная → `close`, внутренние → `router.back()`)
- Доп. кнопка «Закрыть» и стили CTA тянут цвета из Telegram через CSS переменные
