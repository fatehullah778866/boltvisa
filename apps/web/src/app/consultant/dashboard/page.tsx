'use client'

import { useEffect, useState, useCallback } from 'react'
import { useRouter } from 'next/navigation'
import { logger } from '@boltvisa/utils'
import { apiRequest } from '@boltvisa/utils'

import type { AppError } from '@boltvisa/types'
import { User, VisaApplication } from '@boltvisa/types'
import { Button, Card } from '@boltvisa/ui'
import { getUserFriendlyError } from '@/lib/errorHelpers'

interface PaginatedResponse<T> {
  data: T[]
  page: number
  page_size: number
  total: number
  total_pages: number
}

export default function ConsultantDashboardPage() {
  const router = useRouter()
  const [user, setUser] = useState<User | null>(null)
  const [applications, setApplications] = useState<VisaApplication[]>([])
  const [clients, setClients] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedClient, setSelectedClient] = useState<string>('')
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [pagination, setPagination] = useState({ page: 1, total: 0, totalPages: 0 })

  const loadData = useCallback(async () => {
    try {
      const [userData, appsData, clientsData] = await Promise.all([
        apiRequest<User>('/api/v1/users/me'),
        apiRequest<PaginatedResponse<VisaApplication>>(
          `/api/v1/consultant/applications?page=${pagination.page}&status=${statusFilter}&user_id=${selectedClient}`
        ),
        apiRequest<User[]>('/api/v1/consultant/clients'),
      ])

      setUser(userData)
      if (userData.role !== 'consultant' && userData.role !== 'admin') {
        router.push('/dashboard')
        return
      }

      setApplications(appsData.data || [])
      setPagination({
        page: appsData.page,
        total: appsData.total,
        totalPages: appsData.total_pages,
      })
      setClients(clientsData)
    } catch (err) {
      const e = err as AppError
      logger.error('Failed to load data:', e)
      if (e.status === 401 || e.status === 403) {
        localStorage.removeItem('token')
        router.push('/login')
      } else {
        alert(getUserFriendlyError(e))
      }
    } finally {
      setLoading(false)
    }
  }, [router, pagination.page, statusFilter, selectedClient])

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }

    loadData()
  }, [router, loadData])

  useEffect(() => {
    if (user) {
      loadData()
    }
  }, [user, loadData, selectedClient, statusFilter, pagination.page])

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    router.push('/')
  }

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      draft: 'bg-gray-100 text-gray-800',
      submitted: 'bg-blue-100 text-blue-800',
      in_review: 'bg-yellow-100 text-yellow-800',
      approved: 'bg-green-100 text-green-800',
      rejected: 'bg-red-100 text-red-800',
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
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16 items-center">
            <h1 className="text-xl font-bold">Consultant Dashboard</h1>
            <div className="flex items-center gap-4">
              <span className="text-sm text-gray-600">
                {user?.first_name} {user?.last_name}
              </span>
              <Button variant="secondary" size="sm" onClick={handleLogout}>
                Logout
              </Button>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h2 className="text-2xl font-bold mb-2">Client Applications</h2>
          <p className="text-gray-600">Manage applications for your clients</p>
        </div>

        {/* Filters */}
        <Card className="mb-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Filter by Client
              </label>
              <select
                className="w-full px-4 py-2 border border-gray-300 rounded-lg"
                value={selectedClient}
                onChange={(e) => {
                  setSelectedClient(e.target.value)
                  setPagination({ ...pagination, page: 1 })
                }}
              >
                <option value="">All Clients</option>
                {clients.map((client) => (
                  <option key={client.id} value={client.id}>
                    {client.first_name} {client.last_name} ({client.email})
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Filter by Status
              </label>
              <select
                className="w-full px-4 py-2 border border-gray-300 rounded-lg"
                value={statusFilter}
                onChange={(e) => {
                  setStatusFilter(e.target.value)
                  setPagination({ ...pagination, page: 1 })
                }}
              >
                <option value="">All Statuses</option>
                <option value="draft">Draft</option>
                <option value="submitted">Submitted</option>
                <option value="in_review">In Review</option>
                <option value="approved">Approved</option>
                <option value="rejected">Rejected</option>
              </select>
            </div>

            <div className="flex items-end">
              <Button
                variant="secondary"
                onClick={() => {
                  setSelectedClient('')
                  setStatusFilter('')
                  setPagination({ ...pagination, page: 1 })
                }}
              >
                Clear Filters
              </Button>
            </div>
          </div>
        </Card>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <Card>
            <h3 className="text-lg font-semibold mb-2">Total Clients</h3>
            <p className="text-3xl font-bold text-primary-600">{clients.length}</p>
          </Card>
          <Card>
            <h3 className="text-lg font-semibold mb-2">Total Applications</h3>
            <p className="text-3xl font-bold text-blue-600">{pagination.total}</p>
          </Card>
          <Card>
            <h3 className="text-lg font-semibold mb-2">In Review</h3>
            <p className="text-3xl font-bold text-yellow-600">
              {applications.filter((app) => app.status === 'in_review').length}
            </p>
          </Card>
          <Card>
            <h3 className="text-lg font-semibold mb-2">Approved</h3>
            <p className="text-3xl font-bold text-green-600">
              {applications.filter((app) => app.status === 'approved').length}
            </p>
          </Card>
        </div>

        {/* Applications List */}
        <div>
          <h3 className="text-xl font-semibold mb-4">Applications</h3>
          {applications.length === 0 ? (
            <Card>
              <p className="text-gray-600 text-center py-8">No applications found.</p>
            </Card>
          ) : (
            <div className="space-y-4">
              {applications.map((app) => (
                <Card key={app.id} className="hover:shadow-lg transition-shadow">
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <div className="flex items-center gap-4 mb-2">
                        <h4 className="font-semibold text-lg">
                          {app.category?.name || 'Unknown Category'}
                        </h4>
                        <span
                          className={`px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(
                            app.status
                          )}`}
                        >
                          {app.status}
                        </span>
                      </div>
                      <div className="grid grid-cols-2 gap-4 text-sm text-gray-600">
                        <div>
                          <span className="font-medium">Client:</span>{' '}
                          {app.user?.first_name} {app.user?.last_name} ({app.user?.email})
                        </div>
                        {app.passport_number && (
                          <div>
                            <span className="font-medium">Passport:</span> {app.passport_number}
                          </div>
                        )}
                        {app.nationality && (
                          <div>
                            <span className="font-medium">Nationality:</span> {app.nationality}
                          </div>
                        )}
                        <div>
                          <span className="font-medium">Created:</span>{' '}
                          {new Date(app.created_at).toLocaleDateString()}
                        </div>
                      </div>
                    </div>
                    <Button
                      size="sm"
                      onClick={() => router.push(`/consultant/applications/${app.id}`)}
                    >
                      View Details
                    </Button>
                  </div>
                </Card>
              ))}
            </div>
          )}

          {/* Pagination */}
          {pagination.totalPages > 1 && (
            <div className="mt-6 flex justify-center gap-2">
              <Button
                variant="secondary"
                size="sm"
                disabled={pagination.page === 1}
                onClick={() => setPagination({ ...pagination, page: pagination.page - 1 })}
              >
                Previous
              </Button>
              <span className="px-4 py-2 text-sm">
                Page {pagination.page} of {pagination.totalPages}
              </span>
              <Button
                variant="secondary"
                size="sm"
                disabled={pagination.page === pagination.totalPages}
                onClick={() => setPagination({ ...pagination, page: pagination.page + 1 })}
              >
                Next
              </Button>
            </div>
          )}
        </div>
      </main>
    </div>
  )
}

