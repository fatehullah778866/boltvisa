'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { logger } from '@boltvisa/utils'
import { apiRequest } from '@boltvisa/utils'

import type { AppError } from '@boltvisa/types'
import { VisaApplication } from '@boltvisa/types'
import { Button, Card, Input } from '@boltvisa/ui'
import { getUserFriendlyError } from '@/lib/errorHelpers'

interface Payment {
  id: number
  application_id: number
  user_id: number
  amount: number
  currency: string
  status: string
  method: string
  transaction_id?: string
  created_at: string
  application?: VisaApplication
}

export default function PaymentsPage() {
  const router = useRouter()
  const [payments, setPayments] = useState<Payment[]>([])
  const [applications, setApplications] = useState<VisaApplication[]>([])
  const [loading, setLoading] = useState(true)
  const [showPaymentForm, setShowPaymentForm] = useState(false)
  const [selectedApplication, setSelectedApplication] = useState<number | null>(null)
  const [amount, setAmount] = useState('')
  const [paymentMethod, setPaymentMethod] = useState('stripe')

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }

    loadData()
  }, [router])

  const loadData = async () => {
    try {
      const [paymentsData, appsData] = await Promise.all([
        apiRequest<Payment[]>('/api/v1/payments'),
        apiRequest<VisaApplication[]>('/api/v1/applications'),
      ])
      setPayments(paymentsData)
      setApplications(appsData)
    } catch (err) {
      const e = err as AppError
      logger.error('Failed to load data:', e)
      alert(getUserFriendlyError(e))
    } finally {
      setLoading(false)
    }
  }

  const handleCreatePayment = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedApplication || !amount) return

    try {
      const result = await apiRequest<{
        payment_id: number
        payment_intent_id: string
        client_secret: string
        amount: number
        currency: string
      }>('/api/v1/payments', {
        method: 'POST',
        body: {
          application_id: selectedApplication,
          amount: parseFloat(amount),
          currency: 'USD',
          method: paymentMethod,
        },
      })

      // In a real implementation, you would integrate Stripe Elements here
      alert(`Payment created! Payment ID: ${result.payment_id}\nClient Secret: ${result.client_secret}`)
      
      setShowPaymentForm(false)
      setSelectedApplication(null)
      setAmount('')
      loadData()
    } catch (err) {
      const e = err as AppError
      logger.error('Failed to create payment:', e)
      alert(getUserFriendlyError(e))
    }
  }

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: 'bg-yellow-100 text-yellow-800',
      processing: 'bg-blue-100 text-blue-800',
      completed: 'bg-green-100 text-green-800',
      failed: 'bg-red-100 text-red-800',
      refunded: 'bg-gray-100 text-gray-800',
    }
    return colors[status] || 'bg-gray-100 text-gray-800'
  }

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-lg">Loading...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-6 flex justify-between items-center">
          <h1 className="text-3xl font-bold">Payments</h1>
          <Button onClick={() => setShowPaymentForm(!showPaymentForm)}>
            {showPaymentForm ? 'Cancel' : 'New Payment'}
          </Button>
        </div>

        {showPaymentForm && (
          <Card className="mb-6">
            <h2 className="text-xl font-semibold mb-4">Create Payment</h2>
            <form onSubmit={handleCreatePayment} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Application
                </label>
                <select
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg"
                  value={selectedApplication || ''}
                  onChange={(e) => setSelectedApplication(parseInt(e.target.value))}
                  required
                >
                  <option value="">Select an application</option>
                  {applications.map((app) => (
                    <option key={app.id} value={app.id}>
                      {app.category?.name || 'Unknown'} - {app.status}
                    </option>
                  ))}
                </select>
              </div>

              <Input
                type="number"
                label="Amount (USD)"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                step="0.01"
                min="0"
                required
              />

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Payment Method
                </label>
                <select
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg"
                  value={paymentMethod}
                  onChange={(e) => setPaymentMethod(e.target.value)}
                  required
                >
                  <option value="stripe">Stripe</option>
                  <option value="razorpay">Razorpay</option>
                </select>
              </div>

              <div className="flex gap-4">
                <Button type="submit" className="flex-1">
                  Create Payment
                </Button>
                <Button
                  type="button"
                  variant="secondary"
                  onClick={() => {
                    setShowPaymentForm(false)
                    setSelectedApplication(null)
                    setAmount('')
                  }}
                >
                  Cancel
                </Button>
              </div>
            </form>
          </Card>
        )}

        <div>
          <h2 className="text-xl font-semibold mb-4">Payment History</h2>
          {payments.length === 0 ? (
            <Card>
              <div className="text-center py-8">
                <p className="text-gray-600">No payments yet.</p>
              </div>
            </Card>
          ) : (
            <div className="space-y-4">
              {payments.map((payment) => (
                <Card key={payment.id}>
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <div className="flex items-center gap-4 mb-2">
                        <h3 className="font-semibold text-lg">
                          Payment #{payment.id}
                        </h3>
                        <span
                          className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(
                            payment.status
                          )}`}
                        >
                          {payment.status}
                        </span>
                      </div>
                      <div className="grid grid-cols-2 gap-4 text-sm text-gray-600">
                        <div>
                          <span className="font-medium">Amount:</span> ${payment.amount}{' '}
                          {payment.currency}
                        </div>
                        <div>
                          <span className="font-medium">Method:</span> {payment.method}
                        </div>
                        {payment.transaction_id && (
                          <div>
                            <span className="font-medium">Transaction ID:</span>{' '}
                            {payment.transaction_id}
                          </div>
                        )}
                        <div>
                          <span className="font-medium">Date:</span>{' '}
                          {new Date(payment.created_at).toLocaleDateString()}
                        </div>
                        {payment.application && (
                          <div>
                            <span className="font-medium">Application:</span>{' '}
                            {payment.application.category?.name || 'Unknown'}
                          </div>
                        )}
                      </div>
                    </div>
                  </div>
                </Card>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

