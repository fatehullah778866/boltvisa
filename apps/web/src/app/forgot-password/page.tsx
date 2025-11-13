'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { apiRequest } from '@boltvisa/utils'

import { logger } from '@boltvisa/utils'
import { AppError } from '@boltvisa/types'
import { getUserFriendlyError } from '@/lib/errorHelpers'
import { Button, Input, Card } from '@boltvisa/ui'

export default function ForgotPasswordPage() {
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setSuccess(false)
    setLoading(true)

    try {
      await apiRequest('/api/v1/auth/forgot-password', {
        method: 'POST',
        body: { email },
      })

      setSuccess(true)
    } catch (err) {
      const e = err as AppError
      logger.error('Forgot password error', e)
      setError(getUserFriendlyError(e))
    } finally {
      setLoading(false)
    }
  }

  if (success) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
        <Card className="w-full max-w-md">
          <h1 className="text-2xl font-bold text-center mb-6 text-green-600">
            Check Your Email
          </h1>
          <p className="text-gray-600 text-center mb-4">
            If an account with that email exists, we&apos;ve sent a password reset link.
            Please check your email and follow the instructions.
          </p>
          <Button className="w-full" onClick={() => router.push('/login')}>
            Back to Login
          </Button>
        </Card>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
      <Card className="w-full max-w-md">
        <h1 className="text-2xl font-bold text-center mb-6">Forgot Password</h1>
        <p className="text-center text-gray-600 mb-6">
          Enter your email address and we&apos;ll send you a link to reset your password.
        </p>

        {error && (
          <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg">{error}</div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            type="email"
            label="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            disabled={loading}
          />

          <Button type="submit" className="w-full" disabled={loading}>
            {loading ? 'Sending...' : 'Send Reset Link'}
          </Button>
        </form>

        <p className="mt-4 text-center text-sm text-gray-600">
          Remember your password?{' '}
          <a href="/login" className="text-primary-600 hover:underline">
            Login
          </a>
        </p>
      </Card>
    </div>
  )
}

