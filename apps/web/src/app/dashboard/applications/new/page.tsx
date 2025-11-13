'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { logger } from '@boltvisa/utils'
import { apiRequest } from '@boltvisa/utils'

import type { AppError } from '@boltvisa/types'
import { VisaCategory, VisaApplication } from '@boltvisa/types'
import { Button, Input, Card } from '@boltvisa/ui'
import { getUserFriendlyError } from '@/lib/errorHelpers'

export default function NewApplicationPage() {
  const router = useRouter()
  const [categories, setCategories] = useState<VisaCategory[]>([])
  const [formData, setFormData] = useState({
    category_id: '',
    passport_number: '',
    date_of_birth: '',
    nationality: '',
    travel_date: '',
    notes: '',
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }
    loadCategories()
  }, [router])

  const loadCategories = async () => {
    try {
      const data = await apiRequest<VisaCategory[]>('/api/v1/visa-categories')
      setCategories(data)
    } catch (err) {
      logger.error('Failed to load categories:', err)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      await apiRequest<VisaApplication>('/api/v1/applications', {
        method: 'POST',
        body: {
          category_id: parseInt(formData.category_id),
          passport_number: formData.passport_number || undefined,
          date_of_birth: formData.date_of_birth || undefined,
          nationality: formData.nationality || undefined,
          travel_date: formData.travel_date || undefined,
          notes: formData.notes || undefined,
        },
      })

      router.push('/dashboard')
    } catch (err) {
      const e = err as AppError
      logger.error('Create application error', e)
      setError(getUserFriendlyError(e))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-6">
          <Button variant="secondary" onClick={() => router.back()}>
            ← Back
          </Button>
        </div>

        <Card>
          <h1 className="text-2xl font-bold mb-6">New Visa Application</h1>

          {error && (
            <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Visa Category *
              </label>
              <select
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                value={formData.category_id}
                onChange={(e) => setFormData({ ...formData, category_id: e.target.value })}
                required
                disabled={loading}
              >
                <option value="">Select a category</option>
                {categories.map((cat) => (
                  <option key={cat.id} value={cat.id}>
                    {cat.name} - {cat.country} (${cat.price})
                  </option>
                ))}
              </select>
            </div>

            <Input
              type="text"
              label="Passport Number"
              value={formData.passport_number}
              onChange={(e) => setFormData({ ...formData, passport_number: e.target.value })}
              disabled={loading}
            />

            <Input
              type="date"
              label="Date of Birth"
              value={formData.date_of_birth}
              onChange={(e) => setFormData({ ...formData, date_of_birth: e.target.value })}
              disabled={loading}
            />

            <Input
              type="text"
              label="Nationality"
              value={formData.nationality}
              onChange={(e) => setFormData({ ...formData, nationality: e.target.value })}
              disabled={loading}
            />

            <Input
              type="date"
              label="Intended Travel Date"
              value={formData.travel_date}
              onChange={(e) => setFormData({ ...formData, travel_date: e.target.value })}
              disabled={loading}
            />

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Notes
              </label>
              <textarea
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                rows={4}
                value={formData.notes}
                onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                disabled={loading}
              />
            </div>

            <div className="flex gap-4">
              <Button type="submit" disabled={loading} className="flex-1">
                {loading ? 'Creating...' : 'Create Application'}
              </Button>
              <Button
                type="button"
                variant="secondary"
                onClick={() => router.back()}
                disabled={loading}
              >
                Cancel
              </Button>
            </div>
          </form>
        </Card>
      </div>
    </div>
  )
}

