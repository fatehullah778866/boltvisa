'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { apiRequest } from '@boltvisa/utils';
import type { AuthResponse, LoginRequest } from '@boltvisa/types';
import { extractErrorDetails } from '@/lib/errorHelpers';

export default function LoginPage() {
  const router = useRouter();
  const [form, setForm] = useState<LoginRequest>({ email: '', password: '' });
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  const onChange =
    (k: keyof LoginRequest) =>
    (e: React.ChangeEvent<HTMLInputElement>) =>
      setForm((s) => ({ ...s, [k]: e.target.value }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErr(null);
    setLoading(true);
    try {
      const res = await apiRequest<AuthResponse>('/api/v1/auth/login', {
        method: 'POST',
        body: form,
      });
      localStorage.setItem('token', res.token);
      router.replace('/dashboard');
    } catch (e) {
      setErr(extractErrorDetails(e));
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen flex items-center justify-center p-6">
      <form
        onSubmit={handleSubmit}
        className="w-full max-w-md space-y-4 border p-6 rounded-xl"
      >
        <h1 className="text-2xl font-semibold">Sign in</h1>

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

        {err && (
          <p className="text-sm text-red-600 border border-red-200 bg-red-50 p-2 rounded">
            {err}
          </p>
        )}

        <button
          type="submit"
          disabled={loading}
          className="w-full rounded-md px-3 py-2 border bg-black text-white disabled:opacity-60"
        >
          {loading ? 'Signing in…' : 'Sign in'}
        </button>

        {/* 🔗 New Sign Up Link */}
        <p className="text-sm text-center text-gray-600">
          Don’t have an account?{' '}
          <Link
            href="/signup"
            className="text-blue-600 hover:underline font-medium"
          >
            Sign up
          </Link>
        </p>

        <p className="text-xs text-gray-500">
          API: {process.env.NEXT_PUBLIC_API_URL ?? '(not set)'}
        </p>
      </form>
    </main>
  );
}
