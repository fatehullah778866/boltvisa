'use client'

import { useState, useEffect, Suspense } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { apiRequest } from '@boltvisa/utils'

import { logger } from '@boltvisa/utils'
import { AppError } from '@boltvisa/types'
import { getUserFriendlyError } from '@/lib/errorHelpers'
import { Button, Input, Card } from '@boltvisa/ui'

function ResetPasswordForm() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const token = searchParams.get('token')

  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!token) {
      setError('Invalid reset token. Please request a new password reset.')
    }
  }, [token])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setSuccess(false)

    if (password !== confirmPassword) {
      setError('Passwords do not match')
      return
    }

    if (password.length < 8) {
      setError('Password must be at least 8 characters')
      return
    }

    if (!token) {
      setError('Invalid reset token')
      return
    }

    setLoading(true)

    try {
      await apiRequest('/api/v1/auth/reset-password', {
        method: 'POST',
        body: {
          token,
          password,
        },
      })

      setSuccess(true)
      setTimeout(() => {
        router.push('/login')
      }, 2000)
    } catch (err) {
      const e = err as AppError
      logger.error('Reset password error', e)
      setError(getUserFriendlyError(e))
    } finally {
      setLoading(false)
    }
  }

  if (!token) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
        <Card className="w-full max-w-md">
          <h1 className="text-2xl font-bold text-center mb-6">Invalid Reset Link</h1>
          <p className="text-gray-600 text-center mb-4">
            This password reset link is invalid or has expired.
          </p>
          <Button className="w-full" onClick={() => router.push('/forgot-password')}>
            Request New Reset Link
          </Button>
        </Card>
      </div>
    )
  }

  if (success) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
        <Card className="w-full max-w-md">
          <h1 className="text-2xl font-bold text-center mb-6 text-green-600">
            Password Reset Successful
          </h1>
          <p className="text-gray-600 text-center mb-4">
            Your password has been reset successfully. Redirecting to login...
          </p>
        </Card>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
      <Card className="w-full max-w-md">
        <h1 className="text-2xl font-bold text-center mb-6">Reset Password</h1>

        {error && (
          <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg">{error}</div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            type="password"
            label="New Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={8}
            disabled={loading}
          />

          <Input
            type="password"
            label="Confirm Password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            required
            minLength={8}
            disabled={loading}
          />

          <Button type="submit" className="w-full" disabled={loading}>
            {loading ? 'Resetting...' : 'Reset Password'}
          </Button>
        </form>

        <p className="mt-4 text-center text-sm text-gray-600">
          <a href="/login" className="text-primary-600 hover:underline">
            Back to Login
          </a>
        </p>
      </Card>
    </div>
  )
}

export default function ResetPasswordPage() {
  return (
    <Suspense fallback={
      <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
        <Card className="w-full max-w-md">
          <h1 className="text-2xl font-bold text-center mb-6">Loading...</h1>
        </Card>
      </div>
    }>
      <ResetPasswordForm />
    </Suspense>
  )
}

