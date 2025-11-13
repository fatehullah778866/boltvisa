'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiRequest } from '@boltvisa/utils';
import type { VisaApplication } from '@boltvisa/types';
import { extractErrorDetails } from '@/lib/errorHelpers';

export default function ApplicationsPage() {
  const router = useRouter();
  const [apps, setApps] = useState<VisaApplication[]>([]);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState<string | null>(null);

  async function load() {
    setLoading(true);
    setErr(null);
    try {
      const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
      if (!token) throw new Error('No auth token');
      const data = await apiRequest<VisaApplication[]>('/api/v1/applications', {
        headers: { Authorization: `Bearer ${token}` },
      });
      setApps(data ?? []);
    } catch (e) {
      try { localStorage.removeItem('token'); } catch {}
      setErr(extractErrorDetails(e));
      router.replace('/login');
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => { load(); /* eslint-disable-next-line */ }, []);

  if (loading) {
    return <main className="min-h-screen grid place-items-center">Loading applications…</main>;
  }

  return (
    <main className="min-h-screen bg-gray-50">
      <div className="mx-auto max-w-6xl px-6 py-8">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-semibold">Your Applications</h1>
          <button
            onClick={() => router.push('/applications/new')}
            className="rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm shadow-sm hover:shadow transition"
          >
            ➕ New Application
          </button>
        </div>

        {err && <p className="mt-4 text-sm text-red-600">{err}</p>}

        {apps.length === 0 ? (
          <div className="mt-6 rounded-2xl border border-dashed p-10 text-center text-gray-500 bg-white">
            No applications yet. Click “New Application” to get started.
          </div>
        ) : (
          <div className="mt-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {apps.map((a) => (
              <div key={a.id} className="rounded-2xl bg-white p-5 shadow-sm ring-1 ring-gray-100">
                <p className="text-xs uppercase tracking-wide text-gray-500 mb-1">
                  {a.category?.country || '—'} • {a.category?.name || 'Category'}
                </p>
                <div className="text-lg font-semibold">#{a.id}</div>
                <p className="text-sm text-gray-600 mt-1">
                  Status: <span className="font-medium">{a.status}</span>
                </p>
                <p className="text-xs text-gray-500 mt-2">
                  Updated: {new Date(a.updated_at).toLocaleString()}
                </p>
              </div>
            ))}
          </div>
        )}
      </div>
    </main>
  );
}
