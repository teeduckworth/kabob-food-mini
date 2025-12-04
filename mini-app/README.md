# KabobFood Mini App UI

Next.js + Tailwind клиент для Telegram Mini App. Покрывает сценарии меню, корзины, оформления заказа и профиля.

## Запуск

```bash
npm install
npm run dev
```

Необходимые переменные окружения (создайте `.env.local`):

```
NEXT_PUBLIC_API_URL=http://localhost:8080
# опционально для локальной отладки без Telegram
NEXT_PUBLIC_DEV_TOKEN=your-local-jwt
```

## Структура

- `src/app/page.tsx` — меню, категории, карточки товаров
- `src/app/checkout/page.tsx` — оформление заказа
- `src/app/profile/page.tsx` — профиль/адреса (демо-токен)
- `src/components/*` — UI-компоненты (карусель категорий, карточки, bottom-nav, корзина)
- `src/store/cart.ts` — Zustand стора корзины
- `src/lib/api.ts` — минимальный API-клиент

## TODO

- Настроить деплой (Vercel/Pages) и Mini App манифест
- Добавить полноценное редактирование профиля (имя/телефон) и состояние заказа
- Подготовить UI-тесты или стори для основных экранов

## Telegram WebApp интеграция

- `@twa-dev/sdk` сообщает, что мини-апп готово: вызываем `WebApp.ready()`, `expand()` и слушаем `themeChanged`, чтобы подхватывать цвета Telegram.
- Кнопка Back внутри Telegram синхронизирована с роутингом Next.js: на внутренних страницах возвращает назад, на главной — закрывает мини-апп.
- На интерфейсе есть отдельная кнопка «Закрыть», CTA-кнопки используют цветовую схему Telegram через CSS-переменные.
