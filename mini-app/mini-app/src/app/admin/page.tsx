'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { adminApi } from '@/lib/admin-api';
import { useAdminSession } from '@/hooks/useAdminSession';

export default function AdminLoginPage() {
  const router = useRouter();
  const { token, setToken, ready } = useAdminSession();
  const [username, setUsername] = useState('admin');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (ready && token) {
      router.replace('/admin/dashboard');
    }
  }, [ready, token, router]);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (loading) return;
    setError(null);
    setLoading(true);
    try {
      const res = await adminApi.login(username, password);
      setToken(res.token);
      router.push('/admin/dashboard');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Не удалось войти');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="mx-auto max-w-2xl">
      <div className="rounded-3xl border border-white/10 bg-white/5 p-8 shadow-2xl shadow-amber-500/10 backdrop-blur">
        <p className="text-sm text-white/60">Эксклюзивный доступ</p>
        <h2 className="mt-2 text-2xl font-semibold">Войдите в админ-панель</h2>
        <p className="mt-1 text-sm text-white/60">
          Используйте корпоративный аккаунт, чтобы управлять меню и заказами KabobFood.
        </p>

        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-2">
            <label className="text-sm text-white/70" htmlFor="username">
              Логин
            </label>
            <input
              id="username"
              type="text"
              autoComplete="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full rounded-2xl border border-white/20 bg-white/5 px-4 py-3 text-white outline-none transition focus:border-amber-200"
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm text-white/70" htmlFor="password">
              Пароль
            </label>
            <input
              id="password"
              type="password"
              autoComplete="current-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full rounded-2xl border border-white/20 bg-white/5 px-4 py-3 text-white outline-none transition focus:border-amber-200"
            />
          </div>

          {error && <p className="text-sm text-rose-300">{error}</p>}

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-2xl bg-gradient-to-r from-amber-400 via-orange-500 to-amber-500 px-4 py-3 text-center font-semibold text-slate-900 shadow-lg shadow-amber-500/30 transition hover:opacity-90 disabled:opacity-60"
          >
            {loading ? 'Проверяем доступ...' : 'Войти'}
          </button>
        </form>

        <div className="mt-8 grid gap-4 text-sm text-white/60 md:grid-cols-2">
          <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
            <p className="text-xs uppercase tracking-widest text-white/40">Статус</p>
            <p className="mt-1 text-lg font-semibold text-white">99.9% uptime</p>
            <p className="text-xs text-white/50">Инфраструктура на Kubernetes + CDN</p>
          </div>
          <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
            <p className="text-xs uppercase tracking-widest text-white/40">Мониторинг</p>
            <p className="mt-1 text-lg font-semibold text-white">Live metrics</p>
            <p className="text-xs text-white/50">Grafana + Telegram alerts 24/7</p>
          </div>
        </div>
      </div>
    </div>
  );
}
