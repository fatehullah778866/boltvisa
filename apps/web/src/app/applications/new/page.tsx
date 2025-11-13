'use client';

import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { apiRequest } from '@boltvisa/utils';
import type { VisaCategory, VisaApplication } from '@boltvisa/types';
import { extractErrorDetails } from '@/lib/errorHelpers';

type CreateApplicationRequest = {
  category_id: number;
  passport_number?: string;
  date_of_birth?: string;
  nationality?: string;
  travel_date?: string;
  notes?: string;
  submit?: boolean;
};

export default function NewApplicationPage() {
  const router = useRouter();
  const search = useSearchParams();

  const [cats, setCats] = useState<VisaCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  const [form, setForm] = useState<CreateApplicationRequest>({
    category_id: 0,
    passport_number: '',
    date_of_birth: '',
    nationality: '',
    travel_date: '',
    notes: '',
    submit: false,
  });

  // NEW: preselect from ?category_id=...
  useEffect(() => {
    const cid = Number(search.get('category_id') || 0);
    if (cid > 0) {
      setForm((s) => ({ ...s, category_id: cid }));
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const onChange =
    (k: keyof CreateApplicationRequest) =>
    (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
      const value =
        k === 'submit'
          ? (e as React.ChangeEvent<HTMLInputElement>).target.checked
          : e.target.value;
      setForm((s) => ({ ...s, [k]: k === 'category_id' ? Number(value) : (value as any) }));
    };

  async function loadCats() {
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

  useEffect(() => { loadCats(); /* eslint-disable-line react-hooks/exhaustive-deps */ }, []);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setErr(null);
    setSubmitting(true);
    try {
      const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
      if (!token) throw new Error('No auth token');

      await apiRequest<VisaApplication>('/api/v1/applications', {
        method: 'POST',
        headers: { Authorization: `Bearer ${token}` },
        body: form,
      });

      router.replace('/applications');
    } catch (e) {
      setErr(extractErrorDetails(e));
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) {
    return <main className="min-h-screen grid place-items-center">Loading…</main>;
  }

  return (
    <main className="min-h-screen bg-gray-50 grid place-items-center p-6">
      <form onSubmit={submit} className="w-full max-w-2xl rounded-2xl bg-white p-6 shadow-sm ring-1 ring-gray-100 space-y-5">
        <h1 className="text-2xl font-semibold">New Application</h1>

        {err && <p className="text-sm text-red-600 border border-red-200 bg-red-50 p-2 rounded">{err}</p>}

        <div>
          <label className="block text-sm mb-1">Visa category</label>
          <select
            className="w-full border rounded-md px-3 py-2"
            required
            value={form.category_id}
            onChange={onChange('category_id')}
          >
            <option value={0} disabled>Select a category…</option>
            {cats.map((c) => (
              <option key={c.id} value={c.id}>
                {c.country} — {c.name} ({c.duration})
              </option>
            ))}
          </select>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm mb-1">Passport number</label>
            <input className="w-full border rounded-md px-3 py-2" value={form.passport_number ?? ''} onChange={onChange('passport_number')} />
          </div>
          <div>
            <label className="block text-sm mb-1">Nationality</label>
            <input className="w-full border rounded-md px-3 py-2" value={form.nationality ?? ''} onChange={onChange('nationality')} />
          </div>
          <div>
            <label className="block text-sm mb-1">Date of birth</label>
            <input type="date" className="w-full border rounded-md px-3 py-2" value={form.date_of_birth ?? ''} onChange={onChange('date_of_birth')} />
          </div>
          <div>
            <label className="block text-sm mb-1">Intended travel date</label>
            <input type="date" className="w-full border rounded-md px-3 py-2" value={form.travel_date ?? ''} onChange={onChange('travel_date')} />
          </div>
        </div>

        <div>
          <label className="block text-sm mb-1">Notes</label>
          <textarea className="w-full border rounded-md px-3 py-2 min-h-[90px]" value={form.notes ?? ''} onChange={onChange('notes')} />
        </div>

        <div className="flex items-center justify-between">
          <label className="flex items-center gap-2 text-sm">
            <input type="checkbox" checked={!!form.submit} onChange={onChange('submit')} />
            Submit now (otherwise save as draft)
          </label>
          <div className="flex gap-3">
            <button
              type="button"
              className="rounded-xl border border-gray-200 bg-white px-4 py-2 text-sm shadow-sm hover:shadow transition"
              onClick={() => router.push('/applications')}
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={submitting || form.category_id === 0}
              className="rounded-xl border border-gray-200 bg-black text-white px-4 py-2 text-sm shadow-sm hover:shadow transition disabled:opacity-60"
            >
              {submitting ? 'Creating…' : 'Create application'}
            </button>
          </div>
        </div>
      </form>
    </main>
  );
}
