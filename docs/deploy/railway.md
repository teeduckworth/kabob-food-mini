# Railway Deployment Guide

Railway больше не имеет полностью бесплатного тарифа: новый аккаунт получает ~$5 кредита, которого хватает примерно на 500–700 часов работы маленького контейнера (1 vCPU / 512 MB RAM). Как только кредит закончится, придётся привязать карту и оплачивать потребление. Ниже план, как задеплоить бэкенд на Railway с учётом текущей инфраструктуры (Go + PostgreSQL + Redis).

## 1. Подготовка репозитория
- В `Dockerfile` уже собрана Go-служба (`./cmd/app`), Railway умеет автоматически строить контейнер из него.
- Убедитесь, что ветка `main` содержит актуальный код: деплой будет триггериться на каждый push в репозиторий.

## 2. Создание проекта Railway
1. Зарегистрируйтесь/войдите на [railway.app](https://railway.app).
2. Нажмите **New Project → Deploy from GitHub repo**, подключите GitHub-аккаунт и выберите `kabobfood` (нужны права на чтение репо).
3. В разделе **Settings → Deployments** выберите `Dockerfile` как build source (по умолчанию Railway попытается применить Nixpacks, но Docker гарантированно собирается корректно).

## 3. Базы данных и кеш
1. Внутри проекта добавьте **PostgreSQL** сервис (`+ New → Database → PostgreSQL`). Railway создаст отдельный контейнер и выдаст переменные (`DATABASE_URL`, `PGHOST`, `PGUSER`, ...).
2. Аналогично добавьте **Redis** (`+ New → Database → Redis`). Получите переменные `REDIS_URL`, `REDIS_HOST`, `REDIS_PASSWORD` и т.д.
3. Позднее эти переменные пробросим в сервис бэкенда через интерполяции `{{ postgres.DATABASE_URL }}` и `{{ redis.REDIS_URL }}` (Railway UI позволяет ссылаться на другие сервисы).

## 4. Переменные окружения сервиса
Зайдите в созданный сервис `kabobfood` → **Variables** и задайте значения:

| Var | Значение |
| --- | --- |
| `APP_ENV` | `prod` |
| `HTTP_HOST` | `0.0.0.0` (дефолт, можно не задавать) |
| `HTTP_PORT` | `8080` (дефолт; Railway сам пробросит внешний порт) |
| `DB_URL` | `{{ postgres.DATABASE_URL }}` + при необходимости `?sslmode=disable` (или `sslmode=require`, если нужен TLS) |
| `REDIS_URL` | `{{ redis.REDIS_URL }}` |
| `JWT_SECRET` | случайная строка (64+ бит) |
| `TELEGRAM_BOT_TOKEN` | токен вашего бота |
| `TELEGRAM_ADMIN_CHAT_ID` | chat_id администратора для алертов |
| `ADMIN_DEFAULT_USERNAME` / `ADMIN_DEFAULT_PASSWORD` | сгенерируйте безопасные значения, чтобы не использовать дефолт `admin/admin123` |
| `SENTRY_DSN` | опционально, DSN проекта Sentry |
| `RATE_USER_LIMIT`, `RATE_ADMIN_LIMIT`, `RATE_WINDOW` | опционально, лимиты запросов |

> Railway автоматически прокидывает переменную `PORT`, но сервис слушает `HTTP_PORT` (по умолчанию 8080), поэтому ничего дополнительного делать не нужно.

## 5. Запуск миграций
Перед первым запуском выполните SQL-миграции на только что созданной БД. Самый простой путь — локально запустить контейнер `migrate/migrate`, подставив удалённый DSN от Railway:

```bash
# Подставьте полный postgres DSN со страницы Variables
export DB_URL="postgres://user:pass@host:port/db?sslmode=require"

docker run --rm \
  -v "$(pwd)/migrations:/migrations" \
  --network host \
  migrate/migrate:v4.17.1 \
  -path=/migrations -database "$DB_URL" up
```

Альтернатива: установить локально CLI `migrate` и запустить `migrate -path=./migrations -database "$DB_URL" up`. Миграции повторно запускать не нужно — достаточно выполнять их каждый раз, когда в репозитории появляются новые файлы в `migrations/`.

## 6. Деплой
1. Сделайте `git push origin main` — Railway автоматически соберёт Docker-образ (`golang:1.25` builder → `distroless`).
2. Следите за сборкой в **Deployments**. После удачного билда сервис стартует и станет доступен по домену вида `https://kabobfood.up.railway.app`.
3. Проверить здоровье можно по `GET https://kabobfood.up.railway.app/healthz`.
4. Логи доступны через вкладку **Logs** или `railway logs` в CLI.

## 7. CLI (опционально)
- Установите [Railway CLI](https://docs.railway.app/reference/cli), выполните `railway login`, затем `railway link` в корне репозитория.
- Команда `railway variables` покажет текущие ENV, `railway run "cmd"` выполнит локальную команду с этими переменными (удобно для миграций/сидов).

## 8. Что дальше
- Когда понадобится переехать на другой PaaS (Render/Fly/etc.), можно переиспользовать тот же `Dockerfile`; понадобится лишь новый блок env переменных и подключение к их Postgres/Redis.
- Следите за использованием кредитов Railway: как только бесплатные минуты закончатся, сервис будет остановлен.
- Добавьте оповещения (например, в Telegram через вашего бота) на сбои/ошибки после подключения Sentry/логов.

С этими шагами бэкенд можно раскатить за ~10–15 минут, после чего mini-app сможет стучаться в Railway-домен.
