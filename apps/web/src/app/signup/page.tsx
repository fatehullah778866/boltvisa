'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiRequest } from '@boltvisa/utils'

import type { RegisterRequest, AuthResponse } from '@boltvisa/types';
import { extractErrorMessage } from '@/lib/errorHelpers';

export default function SignupPage() {
  const r = useRouter();
  const [form, setForm] = useState<RegisterRequest>({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
  });
  const [loading, setLoading] = useState(false);
  const [errMsg, setErrMsg] = useState<string | null>(null);

  const onChange =
    (k: keyof RegisterRequest) =>
    (e: React.ChangeEvent<HTMLInputElement>) =>
      setForm((s) => ({ ...s, [k]: e.target.value }));

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrMsg(null);
    setLoading(true);
    try {
      const data = await apiRequest<AuthResponse>('/api/v1/auth/register', {
        method: 'POST',
        body: form,
      });
      // store token for demo/dev; in real app prefer httpOnly cookies via API
      localStorage.setItem('token', data.token);
      r.push('/dashboard');
    } catch (err) {
      setErrMsg(extractErrorMessage(err));
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen flex items-center justify-center p-6">
      <form onSubmit={onSubmit} className="w-full max-w-md space-y-4 border p-6 rounded-xl">
        <h1 className="text-2xl font-semibold">Create account</h1>

        <div className="space-y-1">
          <label className="block text-sm">First name</label>
          <input
            className="w-full border rounded-md px-3 py-2"
            value={form.first_name}
            onChange={onChange('first_name')}
            required
          />
        </div>

        <div className="space-y-1">
          <label className="block text-sm">Last name</label>
          <input
            className="w-full border rounded-md px-3 py-2"
            value={form.last_name}
            onChange={onChange('last_name')}
            required
          />
        </div>

        <div className="space-y-1">
          <label className="block text-sm">Email</label>
          <input
            type="email"
            className="w-full border rounded-md px-3 py-2"
            value={form.email}
            onChange={onChange('email')}
            required
          />
        </div>

        <div className="space-y-1">
          <label className="block text-sm">Password</label>
          <input
            type="password"
            className="w-full border rounded-md px-3 py-2"
            value={form.password}
            onChange={onChange('password')}
            minLength={8}
            required
          />
        </div>

        {errMsg && (
          <p className="text-sm text-red-600 border border-red-200 bg-red-50 p-2 rounded">
            {errMsg}
          </p>
        )}

        <button
          type="submit"
          disabled={loading}
          className="w-full rounded-md px-3 py-2 border bg-black text-white disabled:opacity-60"
        >
          {loading ? 'Creating...' : 'Sign up'}
        </button>

        <p className="text-xs text-gray-500">
          API: {process.env.NEXT_PUBLIC_API_URL ?? '(not set)'}
        </p>
      </form>
    </main>
  );
}
