'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiRequest } from '@boltvisa/utils';
import type { User, VisaApplication, VisaCategory } from '@boltvisa/types';
import { extractErrorDetails } from '@/lib/errorHelpers';

// Local DTO for /dashboard response shape
type DashboardStats = {
  total_applications: number;
  draft: number;
  submitted: number;
  in_review: number;
  approved: number;
  rejected: number;
  cancelled: number;
};
type DashboardResponse = {
  stats: DashboardStats;
  recent_applications: VisaApplication[];
};

export default function DashboardPage() {
  const router = useRouter();

  const [user, setUser] = useState<User | null>(null);
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [recent, setRecent] = useState<VisaApplication[]>([]);
  const [cats, setCats] = useState<VisaCategory[]>([]);

  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState<string | null>(null);

  async function loadAll() {
    setLoading(true);
    setErr(null);

    try {
      const token =
        typeof window !== 'undefined' ? localStorage.getItem('token') : null;
      if (!token) throw new Error('No auth token');

      // 1) Current user
      const me = await apiRequest<User>('/api/v1/users/me', {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!me || !me.email || !me.first_name) throw new Error('Invalid user data');
      setUser(me);

      // 2) Dashboard stats + recent
      const dash = await apiRequest<DashboardResponse>('/api/v1/dashboard', {
        headers: { Authorization: `Bearer ${token}` },
      });
      setStats(dash?.stats ?? null);
      setRecent(dash?.recent_applications ?? []);

      // 3) Visa categories strip
      const vc = await apiRequest<VisaCategory[]>('/api/v1/visa-categories', {
        headers: { Authorization: `Bearer ${token}` },
      });
      setCats(vc ?? []);
    } catch (e) {
      try { localStorage.removeItem('token'); } catch {}
      setErr(extractErrorDetails(e));
      router.replace('/login');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadAll();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (loading) {
    return (
      <main className="min-h-screen grid place-items-center bg-gradient-to-b from-gray-50 to-white">
        <div className="animate-pulse space-y-4 w-full max-w-6xl px-6">
          <div className="h-10 w-1/3 rounded-lg bg-gray-200" />
          <div className="h-24 w-full rounded-2xl bg-gray-200" />
          <div className="grid grid-cols-3 gap-4">
            <div className="h-28 rounded-2xl bg-gray-200" />
            <div className="h-28 rounded-2xl bg-gray-200" />
            <div className="h-28 rounded-2xl bg-gray-200" />
          </div>
          <div className="h-64 rounded-2xl bg-gray-200" />
        </div>
      </main>
    );
  }

  if (!user) {
    return (
      <main className="min-h-screen grid place-items-center p-6 bg-gradient-to-b from-gray-50 to-white">
        <p className="text-sm text-red-600">{err ?? 'Redirecting to login…'}</p>
      </main>
    );
  }

  const initials = `${user.first_name?.[0] ?? ''}${user.last_name?.[0] ?? ''}`.toUpperCase();

  return (
    <main className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
      {/* Hero */}
      <section className="relative">
        <div className="absolute inset-0 -z-10 bg-[radial-gradient(1200px_600px_at_50%_-200px,rgba(59,130,246,0.15),transparent)]" />
        <div className="mx-auto max-w-6xl px-6 pt-10 pb-6">
          <div className="flex items-center gap-4">
            <div className="size-12 rounded-full bg-white shadow ring-1 ring-gray-200 grid place-items-center text-gray-700 font-semibold">
              {initials || '🙂'}
            </div>
            <div>
              <h1 className="text-2xl md:text-3xl font-semibold tracking-tight">
                Welcome, {user.first_name}
              </h1>
              <p className="text-sm text-gray-600">{user.email}</p>
            </div>
          </div>

          {/* Quick actions (routes you already have) */}
          <div className="mt-6 flex flex-wrap gap-3">
            <button
              className="rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm shadow-sm hover:shadow transition"
              onClick={() => router.push('/applications/new')}
            >
              ➕ New Application
            </button>
            <button
              className="rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm shadow-sm hover:shadow transition"
              onClick={() => router.push('/visa-categories')}
            >
              🌍 Browse Visa Categories
            </button>
            <button
              className="rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm shadow-sm hover:shadow transition"
              onClick={() => router.push('/applications')}
            >
              📄 View Applications
            </button>
            <button
              className="rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm shadow-sm hover:shadow transition"
              onClick={() => router.push('/notifications')}
            >
              🔔 Notifications
            </button>
          </div>
        </div>
      </section>

      {/* Content */}
      <section className="mx-auto max-w-6xl px-6 pb-12">
        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <StatCard label="Applications" value={stats?.total_applications ?? 0} hint="Total submitted" />
          <StatCard label="In review" value={stats?.in_review ?? 0} hint="Currently being reviewed" />
          <StatCard label="Approved" value={stats?.approved ?? 0} hint="Approved so far" />
        </div>

        {/* Recent + Side */}
        <div className="mt-6 grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 rounded-2xl bg-white p-6 shadow-sm ring-1 ring-gray-100">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold">Recent Applications</h2>
              <button
                className="text-sm text-blue-600 hover:underline"
                onClick={() => router.push('/applications')}
              >
                View all
              </button>
            </div>

            {recent.length === 0 ? (
              <p className="mt-4 text-sm text-gray-500">No recent applications yet.</p>
            ) : (
              <ul className="mt-4 divide-y divide-gray-100">
                {recent.map((a) => (
                  <li key={a.id} className="py-3 flex items-center justify-between">
                    <div>
                      <p className="text-sm font-medium">
                        #{a.id} — {a.category?.name ?? 'Category'} ({a.category?.country ?? '—'})
                      </p>
                      <p className="text-xs text-gray-500">
                        Status: <span className="font-medium">{a.status}</span> • Updated{' '}
                        {new Date(a.updated_at).toLocaleString()}
                      </p>
                    </div>
                    <button
                      className="text-sm rounded-lg border border-gray-200 bg-white px-3 py-1.5 shadow-sm hover:shadow transition"
                      onClick={() => router.push('/applications')}
                    >
                      Open
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </div>

          <aside className="rounded-2xl bg-white p-6 shadow-sm ring-1 ring-gray-100">
            <h3 className="text-lg font-semibold">Quick glance</h3>
            <div className="mt-3 space-y-2 text-sm text-gray-700">
              <Progress label="Draft" value={stats?.draft ?? 0} total={stats?.total_applications ?? 0} />
              <Progress label="Submitted" value={stats?.submitted ?? 0} total={stats?.total_applications ?? 0} />
              <Progress label="In review" value={stats?.in_review ?? 0} total={stats?.total_applications ?? 0} />
              <Progress label="Approved" value={stats?.approved ?? 0} total={stats?.total_applications ?? 0} />
              <Progress label="Rejected" value={stats?.rejected ?? 0} total={stats?.total_applications ?? 0} />
            </div>
          </aside>
        </div>

        {/* Categories strip */}
        <div className="mt-8">
          <h3 className="text-lg font-semibold">Popular Categories</h3>
          {cats.length === 0 ? (
            <p className="mt-3 text-sm text-gray-500">No categories available.</p>
          ) : (
            <div className="mt-3 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {cats.slice(0, 6).map((c) => (
                <div
                  key={c.id}
                  className="rounded-2xl bg-white p-5 shadow-sm ring-1 ring-gray-100"
                >
                  <p className="text-xs uppercase tracking-wide text-gray-500">{c.country}</p>
                  <h4 className="mt-1 text-base font-semibold">{c.name}</h4>
                  <p className="mt-1 text-sm text-gray-600 line-clamp-2">{c.description}</p>
                  <p className="mt-3 text-sm text-gray-700">
                    Duration: <span className="font-medium">{c.duration}</span>
                  </p>
                  <p className="text-sm text-gray-700">
                    Fee: <span className="font-medium">${c.price}</span>
                  </p>
                  <div className="mt-4">
                    <button
                      className="rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm shadow-sm hover:shadow transition"
                      onClick={() => router.push('/applications/new')}
                    >
                      Apply
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </section>
    </main>
  );
}

/* ---------------- UI helpers (pure presentational) ---------------- */

function StatCard({ label, value, hint }: { label: string; value: number; hint: string }) {
  return (
    <div className="rounded-2xl bg-white p-5 shadow-sm ring-1 ring-gray-100">
      <p className="text-xs uppercase tracking-wide text-gray-500">{label}</p>
      <div className="mt-2 text-3xl font-semibold">{value}</div>
      <p className="mt-1 text-xs text-gray-500">{hint}</p>
    </div>
  );
}

function Progress({ label, value, total }: { label: string; value: number; total: number }) {
  const pct = total > 0 ? Math.min(100, Math.round((value / total) * 100)) : 0;
  return (
    <div>
      <div className="flex items-center justify-between text-xs text-gray-600 mb-1">
        <span>{label}</span>
        <span>{value} {total > 0 ? `(${pct}%)` : ''}</span>
      </div>
      <div className="h-2 w-full bg-gray-100 rounded-full overflow-hidden">
        <div className="h-full w-0 rounded-full bg-blue-500" style={{ width: `${pct}%` }} />
      </div>
    </div>
  );
}
