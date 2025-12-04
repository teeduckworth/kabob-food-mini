'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import { useRouter } from 'next/navigation';
import { adminApi, type AdminCategoryInput, type AdminProductInput } from '@/lib/admin-api';
import { useAdminSession } from '@/hooks/useAdminSession';
import type { MenuCategory, Product } from '@/types/api';

interface EnrichedProduct extends Product {
  categoryName: string;
  categoryEmoji?: string;
}

const createCategoryForm = (): Omit<AdminCategoryInput, 'sort_order'> & { sort_order: string } => ({
  name: '',
  emoji: '',
  is_active: true,
  sort_order: '10',
});

const createProductForm = (): Omit<AdminProductInput, 'price' | 'old_price' | 'sort_order' | 'category_id'> & {
  category_id: string;
  price: string;
  old_price: string;
  sort_order: string;
} => ({
  category_id: '',
  name: '',
  description: '',
  price: '',
  old_price: '',
  image_url: '',
  is_active: true,
  sort_order: '100',
});

export default function AdminDashboardPage() {
  const router = useRouter();
  const { token, ready, logout } = useAdminSession();
  const [menu, setMenu] = useState<MenuCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [banner, setBanner] = useState<string | null>(null);
  const [categoryForm, setCategoryForm] = useState(createCategoryForm);
  const [productForm, setProductForm] = useState(createProductForm);
  const [editingCategoryId, setEditingCategoryId] = useState<number | null>(null);
  const [editingProductId, setEditingProductId] = useState<number | null>(null);
  const [submittingCategory, setSubmittingCategory] = useState(false);
  const [submittingProduct, setSubmittingProduct] = useState(false);
  const [productActionId, setProductActionId] = useState<number | null>(null);
  const [categoryActionId, setCategoryActionId] = useState<number | null>(null);

  const loadMenu = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await adminApi.getMenu();
      setMenu(data.categories);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Не удалось загрузить меню');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (!ready) return;
    if (!token) {
      router.replace('/admin');
      return;
    }
    loadMenu();
  }, [ready, token, router, loadMenu]);

  const flatProducts = useMemo<EnrichedProduct[]>(() => {
    return menu.flatMap((category) =>
      category.products.map((product) => ({
        ...product,
        categoryName: category.name,
        categoryEmoji: category.emoji,
      }))
    );
  }, [menu]);

  const stats = useMemo(() => {
    const activeProducts = flatProducts.filter((product) => product.is_active).length;
    const avgPrice = flatProducts.length
      ? Math.round(
          (flatProducts.reduce((sum, product) => sum + product.price, 0) / flatProducts.length) * 100
        ) / 100
      : 0;
    return [
      { label: 'Категорий', value: menu.length },
      { label: 'Позиции в меню', value: flatProducts.length },
      { label: 'Активных блюд', value: activeProducts },
      { label: 'Средний чек, ₽', value: avgPrice },
    ];
  }, [menu.length, flatProducts]);

  function showBanner(message: string) {
    setBanner(message);
    setTimeout(() => setBanner(null), 4000);
  }

  function resetCategoryForm() {
    setCategoryForm(createCategoryForm());
    setEditingCategoryId(null);
  }

  function resetProductForm() {
    setProductForm(createProductForm());
    setEditingProductId(null);
  }

  async function handleCategorySubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!token || submittingCategory) return;
    if (!categoryForm.name.trim()) {
      setError('Введите название категории');
      return;
    }
    setSubmittingCategory(true);
    setError(null);
    try {
      const payload: AdminCategoryInput = {
        name: categoryForm.name.trim(),
        emoji: categoryForm.emoji?.trim() ?? '',
        sort_order: Number(categoryForm.sort_order) || 0,
        is_active: categoryForm.is_active,
      };
      if (editingCategoryId) {
        await adminApi.updateCategory(token, editingCategoryId, payload);
        showBanner('Категория обновлена');
      } else {
        await adminApi.createCategory(token, payload);
        showBanner('Категория добавлена');
      }
      resetCategoryForm();
      await loadMenu();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка создания категории');
    } finally {
      setSubmittingCategory(false);
    }
  }

  async function handleProductSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!token || submittingProduct) return;
    if (!productForm.name.trim()) {
      setError('Введите название блюда');
      return;
    }
    if (!productForm.category_id) {
      setError('Выберите категорию');
      return;
    }
    if (!productForm.price) {
      setError('Укажите цену');
      return;
    }
    setSubmittingProduct(true);
    setError(null);
    try {
      const payload: AdminProductInput = {
        category_id: Number(productForm.category_id),
        name: productForm.name.trim(),
        description: productForm.description?.trim() ?? '',
        price: Number(productForm.price),
        old_price: productForm.old_price ? Number(productForm.old_price) : 0,
        image_url: productForm.image_url?.trim() ?? '',
        is_active: productForm.is_active,
        sort_order: Number(productForm.sort_order) || 0,
      };
      if (editingProductId) {
        await adminApi.updateProduct(token, editingProductId, payload);
        showBanner('Блюдо обновлено');
      } else {
        await adminApi.createProduct(token, payload);
        showBanner('Блюдо добавлено');
      }
      resetProductForm();
      await loadMenu();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка создания блюда');
    } finally {
      setSubmittingProduct(false);
    }
  }

  async function handleToggleProduct(product: EnrichedProduct) {
    if (!token) return;
    setProductActionId(product.id);
    setError(null);
    try {
      const payload: AdminProductInput = {
        category_id: product.category_id,
        name: product.name,
        description: product.description ?? '',
        price: product.price,
        old_price: product.old_price ?? 0,
        image_url: product.image_url ?? '',
        is_active: !product.is_active,
        sort_order: product.sort_order,
      };
      await adminApi.updateProduct(token, product.id, payload);
      showBanner(product.is_active ? 'Блюдо скрыто в меню' : 'Блюдо снова активно');
      await loadMenu();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка обновления блюда');
    } finally {
      setProductActionId(null);
    }
  }

  function handleEditCategory(category: MenuCategory) {
    setCategoryForm({
      name: category.name,
      emoji: category.emoji ?? '',
      sort_order: String(category.sort_order ?? 0),
      is_active: category.is_active ?? true,
    });
    setEditingCategoryId(category.id);
  }

  async function handleDeleteCategory(id: number) {
    if (!token) return;
    if (!window.confirm('Удалить категорию? Блюда внутри тоже исчезнут.')) return;
    setCategoryActionId(id);
    setError(null);
    try {
      await adminApi.deleteCategory(token, id);
      showBanner('Категория удалена');
      if (editingCategoryId === id) {
        resetCategoryForm();
      }
      await loadMenu();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка удаления категории');
    } finally {
      setCategoryActionId(null);
    }
  }

  function handleEditProduct(product: EnrichedProduct) {
    setProductForm({
      category_id: String(product.category_id),
      name: product.name,
      description: product.description ?? '',
      price: String(product.price),
      old_price: product.old_price ? String(product.old_price) : '',
      image_url: product.image_url ?? '',
      is_active: product.is_active,
      sort_order: String(product.sort_order),
    });
    setEditingProductId(product.id);
  }

  async function handleDeleteProduct(id: number) {
    if (!token) return;
    if (!window.confirm('Удалить блюдо?')) return;
    setProductActionId(id);
    setError(null);
    try {
      await adminApi.deleteProduct(token, id);
      showBanner('Блюдо удалено');
      if (editingProductId === id) {
        resetProductForm();
      }
      await loadMenu();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка удаления блюда');
    } finally {
      setProductActionId(null);
    }
  }

  if (!ready || !token) {
    return null;
  }

  return (
    <div className="space-y-10">
      <div className="flex flex-col gap-4 rounded-3xl border border-white/10 bg-white/5 p-6 backdrop-blur md:flex-row md:items-center md:justify-between">
        <div>
          <p className="text-sm text-white/60">Операционный обзор</p>
          <h2 className="mt-1 text-2xl font-semibold">Актуальное меню KabobFood</h2>
        </div>
        <button
          onClick={() => {
            logout();
            router.push('/admin');
          }}
          className="rounded-full border border-white/20 px-5 py-2 text-sm font-medium text-white/80 transition hover:border-white hover:text-white"
        >
          Выйти из панели
        </button>
      </div>

      {banner && (
        <div className="rounded-2xl border border-emerald-300/40 bg-emerald-400/10 px-4 py-3 text-sm text-emerald-100">
          {banner}
        </div>
      )}

      {error && (
        <div className="rounded-2xl border border-rose-300/40 bg-rose-500/10 px-4 py-3 text-sm text-rose-100">
          {error}
        </div>
      )}

      <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {stats.map((stat) => (
          <div key={stat.label} className="rounded-3xl border border-white/10 bg-white/5 p-5">
            <p className="text-xs uppercase tracking-[0.3em] text-white/40">{stat.label}</p>
            <p className="mt-3 text-3xl font-semibold text-white">{stat.value || 0}</p>
          </div>
        ))}
      </section>

      <section className="grid gap-6 lg:grid-cols-2">
        <div className="rounded-3xl border border-white/10 bg-white/5 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-white/60">{editingCategoryId ? 'Редактировать категорию' : 'Добавить категорию'}</p>
              <h3 className="text-xl font-semibold">
                {editingCategoryId ? 'Обновление витрины' : 'Новая витрина'}
              </h3>
            </div>
            {editingCategoryId ? (
              <button
                type="button"
                onClick={resetCategoryForm}
                className="text-xs uppercase tracking-[0.3em] text-white/60 hover:text-white"
              >
                Отменить
              </button>
            ) : (
              <span className="text-xs uppercase tracking-[0.3em] text-white/40">Luxury tier</span>
            )}
          </div>
          <form className="mt-6 space-y-4" onSubmit={handleCategorySubmit}>
            <div className="grid gap-3 md:grid-cols-2">
              <label className="space-y-2 text-sm text-white/70">
                Название
                <input
                  value={categoryForm.name}
                  onChange={(e) => setCategoryForm((prev) => ({ ...prev, name: e.target.value }))}
                  className="w-full rounded-2xl border border-white/15 bg-white/5 px-4 py-3 text-white focus:border-amber-200"
                />
              </label>
              <label className="space-y-2 text-sm text-white/70">
                Emoji
                <input
                  value={categoryForm.emoji}
                  onChange={(e) => setCategoryForm((prev) => ({ ...prev, emoji: e.target.value }))}
                  className="w-full rounded-2xl border border-white/15 bg-white/5 px-4 py-3 text-white focus:border-amber-200"
                />
              </label>
            </div>
            <div className="grid gap-3 md:grid-cols-2">
              <label className="space-y-2 text-sm text-white/70">
                Порядок сортировки
                <input
                  type="number"
                  value={categoryForm.sort_order}
                  onChange={(e) => setCategoryForm((prev) => ({ ...prev, sort_order: e.target.value }))}
                  className="w-full rounded-2xl border border-white/15 bg-white/5 px-4 py-3 text-white focus:border-amber-200"
                />
              </label>
              <label className="flex items-center gap-3 text-sm text-white/70">
                <input
                  type="checkbox"
                  checked={categoryForm.is_active}
                  onChange={(e) => setCategoryForm((prev) => ({ ...prev, is_active: e.target.checked }))}
                  className="size-5 rounded border-white/30 bg-white/10"
                />
                Показать в меню
              </label>
            </div>
            <button
              type="submit"
              disabled={submittingCategory}
              className="w-full rounded-2xl bg-gradient-to-r from-amber-400 via-orange-500 to-amber-500 px-4 py-3 text-center font-semibold text-slate-900 shadow-lg shadow-amber-500/30 hover:opacity-90 disabled:opacity-60"
            >
              {submittingCategory
                ? 'Сохраняем...'
                : editingCategoryId
                  ? 'Обновить категорию'
                  : 'Создать категорию'}
            </button>
          </form>
        </div>

        <div className="rounded-3xl border border-white/10 bg-white/5 p-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-white/60">{editingProductId ? 'Редактировать блюдо' : 'Добавить блюдо'}</p>
              <h3 className="text-xl font-semibold">
                {editingProductId ? 'Обновление позиции' : 'Новая позиция'}
              </h3>
            </div>
            {editingProductId ? (
              <button
                type="button"
                onClick={resetProductForm}
                className="text-xs uppercase tracking-[0.3em] text-white/60 hover:text-white"
              >
                Отменить
              </button>
            ) : (
              <span className="text-xs uppercase tracking-[0.3em] text-white/40">Chef signature</span>
            )}
          </div>
          <form className="mt-6 space-y-4" onSubmit={handleProductSubmit}>
            <label className="space-y-2 text-sm text-white/70">
              Категория
              <select
                value={productForm.category_id}
                onChange={(e) => setProductForm((prev) => ({ ...prev, category_id: e.target.value }))}
                className="w-full rounded-2xl border border-white/15 bg-slate-900/70 px-4 py-3 text-white focus:border-amber-200"
              >
                <option value="">Выберите категорию</option>
                {menu.map((category) => (
                  <option key={category.id} value={category.id}>
                    {category.emoji ? `${category.emoji} ` : ''}
                    {category.name}
                  </option>
                ))}
              </select>
            </label>

            <div className="grid gap-3 md:grid-cols-2">
              <label className="space-y-2 text-sm text-white/70">
                Название
                <input
                  value={productForm.name}
                  onChange={(e) => setProductForm((prev) => ({ ...prev, name: e.target.value }))}
                  className="w-full rounded-2xl border border-white/15 bg-white/5 px-4 py-3 text-white focus:border-amber-200"
                />
              </label>
              <label className="space-y-2 text-sm text-white/70">
                Цена, ₽
                <input
                  type="number"
                  step="0.01"
                  value={productForm.price}
                  onChange={(e) => setProductForm((prev) => ({ ...prev, price: e.target.value }))}
                  className="w-full rounded-2xl border border-white/15 bg-white/5 px-4 py-3 text-white focus:border-amber-200"
                />
              </label>
            </div>

            <label className="space-y-2 text-sm text-white/70">
              Описание
              <textarea
                value={productForm.description}
                onChange={(e) => setProductForm((prev) => ({ ...prev, description: e.target.value }))}
                className="w-full rounded-2xl border border-white/15 bg-white/5 px-4 py-3 text-white focus:border-amber-200"
                rows={3}
              />
            </label>

            <div className="grid gap-3 md:grid-cols-3">
              <label className="space-y-2 text-sm text-white/70">
                Стар. цена
                <input
                  type="number"
                  step="0.01"
                  value={productForm.old_price}
                  onChange={(e) => setProductForm((prev) => ({ ...prev, old_price: e.target.value }))}
                  className="w-full rounded-2xl border border-white/15 bg-white/5 px-4 py-3 text-white focus:border-amber-200"
                />
              </label>
              <label className="space-y-2 text-sm text-white/70">
                Сортировка
                <input
                  type="number"
                  value={productForm.sort_order}
                  onChange={(e) => setProductForm((prev) => ({ ...prev, sort_order: e.target.value }))}
                  className="w-full rounded-2xl border border-white/15 bg-white/5 px-4 py-3 text-white focus:border-amber-200"
                />
              </label>
              <label className="space-y-2 text-sm text-white/70">
                Фото (URL)
                <input
                  type="url"
                  value={productForm.image_url}
                  onChange={(e) => setProductForm((prev) => ({ ...prev, image_url: e.target.value }))}
                  className="w-full rounded-2xl border border-white/15 bg-white/5 px-4 py-3 text-white focus:border-amber-200"
                />
              </label>
            </div>

            <label className="flex items-center gap-3 text-sm text-white/70">
              <input
                type="checkbox"
                checked={productForm.is_active}
                onChange={(e) => setProductForm((prev) => ({ ...prev, is_active: e.target.checked }))}
                className="size-5 rounded border-white/30 bg-white/10"
              />
              Показать в меню
            </label>

            <button
              type="submit"
              disabled={submittingProduct}
              className="w-full rounded-2xl bg-gradient-to-r from-amber-400 via-orange-500 to-amber-500 px-4 py-3 text-center font-semibold text-slate-900 shadow-lg shadow-amber-500/30 hover:opacity-90 disabled:opacity-60"
            >
              {submittingProduct
                ? 'Сохраняем...'
                : editingProductId
                  ? 'Обновить блюдо'
                  : 'Создать блюдо'}
            </button>
          </form>
        </div>
      </section>

      <section className="rounded-3xl border border-white/10 bg-white/5 p-6">
        <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
          <div>
            <p className="text-sm text-white/60">Категории</p>
            <h3 className="text-xl font-semibold">Витрины и подборки</h3>
          </div>
          <p className="text-sm text-white/60">Клиентский интерфейс показывает только активные позиции.</p>
        </div>

        {loading ? (
          <p className="mt-6 text-sm text-white/60">Загружаем структуру меню...</p>
        ) : (
          <div className="mt-6 grid gap-4 md:grid-cols-2">
            {menu.map((category) => (
              <div key={category.id} className="rounded-2xl border border-white/10 bg-white/5 p-4">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-lg font-semibold text-white">
                      {category.emoji && <span className="mr-2 text-xl">{category.emoji}</span>}
                      {category.name}
                    </p>
                    <p className="text-xs uppercase tracking-[0.3em] text-white/40">{category.products.length} позиций</p>
                  </div>
                  <span className="text-sm text-white/60">#{category.sort_order}</span>
                </div>
                <div className="mt-4 flex flex-wrap gap-2 text-xs text-white/60">
                  {category.products.slice(0, 4).map((product) => (
                    <span key={product.id} className="rounded-full border border-white/10 px-3 py-1">
                      {product.name}
                    </span>
                  ))}
                  {category.products.length === 0 && (
                    <span className="rounded-full border border-dashed border-white/20 px-3 py-1 text-white/40">
                      Пока пусто
                    </span>
                  )}
                </div>
                <div className="mt-4 flex flex-wrap gap-2 text-xs">
                  <button
                    type="button"
                    onClick={() => handleEditCategory(category)}
                    className="rounded-full border border-white/20 px-4 py-2 text-white/80 transition hover:border-white"
                  >
                    Изменить
                  </button>
                  <button
                    type="button"
                    disabled={categoryActionId === category.id}
                    onClick={() => handleDeleteCategory(category.id)}
                    className="rounded-full border border-rose-300/40 px-4 py-2 text-rose-200 transition hover:border-rose-200 disabled:opacity-50"
                  >
                    {categoryActionId === category.id ? 'Удаляем...' : 'Удалить'}
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </section>

      <section className="rounded-3xl border border-white/10 bg-white/5 p-6">
        <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
          <div>
            <p className="text-sm text-white/60">Позиции</p>
            <h3 className="text-xl font-semibold">Контроль состава меню</h3>
          </div>
          <p className="text-sm text-white/60">Мгновенно скрывайте и возвращайте блюда из выдачи.</p>
        </div>
        {flatProducts.length === 0 ? (
          <p className="mt-6 text-sm text-white/60">Добавьте хотя бы одну позицию, чтобы увидеть список.</p>
        ) : (
          <div className="mt-6 space-y-3">
            {flatProducts.map((product) => (
              <div
                key={product.id}
                className="flex flex-col gap-3 rounded-2xl border border-white/10 bg-slate-900/50 px-4 py-3 md:flex-row md:items-center md:justify-between"
              >
                <div>
                  <p className="text-lg font-semibold text-white">{product.name}</p>
                  <p className="text-sm text-white/60">
                    {product.categoryEmoji && <span className="mr-2">{product.categoryEmoji}</span>}
                    {product.categoryName}
                  </p>
                </div>
                <div className="flex flex-col gap-2 text-sm text-white/70 md:flex-row md:items-center">
                  <span className="font-semibold text-amber-200">{product.price} ₽</span>
                  <span className={`text-xs ${product.is_active ? 'text-emerald-300' : 'text-white/40'}`}>
                    {product.is_active ? 'Показывается клиентам' : 'Скрыто'}
                  </span>
                  <div className="flex flex-wrap gap-2 text-xs">
                    <button
                      onClick={() => handleEditProduct(product)}
                      className="rounded-full border border-white/20 px-4 py-2 text-white/80 transition hover:border-white"
                    >
                      Изменить
                    </button>
                    <button
                      onClick={() => handleToggleProduct(product)}
                      disabled={productActionId === product.id}
                      className="rounded-full border border-white/20 px-4 py-2 text-white/80 transition hover:border-white disabled:opacity-50"
                    >
                      {productActionId === product.id
                        ? 'Сохраняем...'
                        : product.is_active
                          ? 'Скрыть'
                          : 'Вернуть в меню'}
                    </button>
                    <button
                      onClick={() => handleDeleteProduct(product.id)}
                      disabled={productActionId === product.id}
                      className="rounded-full border border-rose-300/40 px-4 py-2 text-rose-200 transition hover:border-rose-200 disabled:opacity-50"
                    >
                      {productActionId === product.id ? 'Удаляем...' : 'Удалить'}
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
