'use client';

import { useEffect, useMemo, useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiRequest } from '@boltvisa/utils';
import type { VisaCategory } from '@boltvisa/types';
import { extractErrorDetails } from '@/lib/errorHelpers';

type SortKey = 'relevance' | 'priceAsc' | 'priceDesc' | 'nameAsc' | 'nameDesc';

export default function VisaCategoriesPage() {
  const router = useRouter();
  const [cats, setCats] = useState<VisaCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState<string | null>(null);

  // UI controls
  const [query, setQuery] = useState('');
  const [country, setCountry] = useState<'all' | string>('all');
  const [sort, setSort] = useState<SortKey>('relevance');

  async function load() {
    setLoading(true);
    setErr(null);
    try {
      const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
      if (!token) throw new Error('No auth token');

      const data = await apiRequest<VisaCategory[]>('/api/v1/visa-categories', {
        headers: { Authorization: `Bearer ${token}` },
      });
      setCats(data ?? []);
    } catch (e) {
      try { localStorage.removeItem('token'); } catch {}
      setErr(extractErrorDetails(e));
      router.replace('/login');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => { load(); /* eslint-disable-next-line */ }, []);

  const countries = useMemo(() => {
    const set = new Set<string>();
    cats.forEach(c => set.add(c.country));
    return Array.from(set).sort((a,b) => a.localeCompare(b));
  }, [cats]);

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    let arr = cats.filter(c => {
      const matchesText = !q ||
        c.name.toLowerCase().includes(q) ||
        c.description.toLowerCase().includes(q) ||
        c.country.toLowerCase().includes(q);
      const matchesCountry = country === 'all' || c.country === country;
      return matchesText && matchesCountry;
    });

    switch (sort) {
      case 'priceAsc':
        arr = [...arr].sort((a,b) => a.price - b.price);
        break;
      case 'priceDesc':
        arr = [...arr].sort((a,b) => b.price - a.price);
        break;
      case 'nameAsc':
        arr = [...arr].sort((a,b) => a.name.localeCompare(b.name));
        break;
      case 'nameDesc':
        arr = [...arr].sort((a,b) => b.name.localeCompare(a.name));
        break;
      case 'relevance':
      default:
        // naïve relevance: prioritize name/country matches
        arr = [...arr].sort((a,b) => {
          const aw = weight(a, q);
          const bw = weight(b, q);
          return bw - aw;
        });
    }
    return arr;
  }, [cats, query, country, sort]);

  function weight(c: VisaCategory, q: string) {
    if (!q) return 0;
    const name = c.name.toLowerCase();
    const country = c.country.toLowerCase();
    const desc = c.description.toLowerCase();
    let w = 0;
    if (name.includes(q)) w += 3;
    if (country.includes(q)) w += 2;
    if (desc.includes(q)) w += 1;
    return w;
  }

  if (loading) {
    return (
      <main className="min-h-screen grid place-items-center bg-gray-50">
        <div className="animate-pulse space-y-4 w-full max-w-6xl px-6">
          <div className="h-10 w-1/3 rounded-lg bg-gray-200" />
          <div className="grid grid-cols-3 gap-4">
            <div className="h-24 rounded-2xl bg-gray-200" />
            <div className="h-24 rounded-2xl bg-gray-200" />
            <div className="h-24 rounded-2xl bg-gray-200" />
          </div>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-gray-50">
      <div className="mx-auto max-w-6xl px-6 py-8">
        <div className="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
          <div>
            <h1 className="text-2xl font-semibold">Visa Categories</h1>
            <p className="text-sm text-gray-600">Browse, search, and apply in a click.</p>
          </div>

          {/* Controls */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <input
              value={query}
              onChange={e => setQuery(e.target.value)}
              placeholder="Search by name, country, description…"
              className="w-full border rounded-md px-3 py-2"
            />
            <select
              value={country}
              onChange={e => setCountry(e.target.value)}
              className="w-full border rounded-md px-3 py-2"
            >
              <option value="all">All countries</option>
              {countries.map(c => (
                <option key={c} value={c}>{c}</option>
              ))}
            </select>
            <select
              value={sort}
              onChange={e => setSort(e.target.value as SortKey)}
              className="w-full border rounded-md px-3 py-2"
            >
              <option value="relevance">Sort: Relevance</option>
              <option value="priceAsc">Sort: Price (low → high)</option>
              <option value="priceDesc">Sort: Price (high → low)</option>
              <option value="nameAsc">Sort: Name (A → Z)</option>
              <option value="nameDesc">Sort: Name (Z → A)</option>
            </select>
          </div>
        </div>

        {err && <p className="mt-4 text-sm text-red-600">{err}</p>}

        {filtered.length === 0 ? (
          <div className="mt-6 rounded-2xl border border-dashed p-10 text-center text-gray-500 bg-white">
            No categories found. Try changing your filters.
          </div>
        ) : (
          <div className="mt-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filtered.map((c) => (
              <div key={c.id} className="rounded-2xl bg-white p-5 shadow-sm ring-1 ring-gray-100">
                <p className="text-xs uppercase tracking-wide text-gray-500">{c.country}</p>
                <h3 className="mt-1 text-lg font-semibold">{c.name}</h3>
                <p className="mt-1 text-sm text-gray-600 line-clamp-3">{c.description}</p>
                <div className="mt-3 grid grid-cols-2 gap-2 text-sm">
                  <div className="rounded-lg bg-gray-50 p-2">
                    <p className="text-xs text-gray-500">Duration</p>
                    <p className="font-medium">{c.duration}</p>
                  </div>
                  <div className="rounded-lg bg-gray-50 p-2">
                    <p className="text-xs text-gray-500">Fee</p>
                    <p className="font-medium">${c.price}</p>
                  </div>
                </div>
                <div className="mt-4 flex gap-3">
                  <button
                    className="rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm shadow-sm hover:shadow transition"
                    onClick={() => router.push(`/applications/new?category_id=${c.id}`)}
                  >
                    Apply
                  </button>
                  <button
                    className="rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm shadow-sm hover:shadow transition"
                    onClick={() => router.push('/applications')}
                  >
                    View applications
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </main>
  );
}
