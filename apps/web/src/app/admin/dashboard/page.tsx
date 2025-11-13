'use client'

import { useEffect, useState, useCallback } from 'react'
import { useRouter } from 'next/navigation'
import { logger } from '@boltvisa/utils'
import { apiRequest } from '@boltvisa/utils'

import type { AppError } from '@boltvisa/types'
import { User, VisaCategory, VisaApplication } from '@boltvisa/types'
import { Button, Card, Input } from '@boltvisa/ui'
import { getUserFriendlyError } from '@/lib/errorHelpers'

interface PaginatedResponse<T> {
  data: T[]
  page: number
  page_size: number
  total: number
  total_pages: number
}

export default function AdminDashboardPage() {
  const router = useRouter()
  const [user, setUser] = useState<User | null>(null)
  const [users, setUsers] = useState<User[]>([])
  const [categories, setCategories] = useState<VisaCategory[]>([])
  const [applications, setApplications] = useState<VisaApplication[]>([])
  const [activeTab, setActiveTab] = useState<'users' | 'categories' | 'applications'>('users')
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [userPagination, setUserPagination] = useState({ page: 1, total: 0, totalPages: 0 })
  const [appPagination, setAppPagination] = useState({ page: 1, total: 0, totalPages: 0 })

  const loadUser = useCallback(async () => {
    try {
      const userData = await apiRequest<User>('/api/v1/users/me')
      setUser(userData)
    } catch (err) {
      const e = err as AppError
      logger.error('Failed to load user:', e)
      if (e.status === 401 || e.status === 403) {
        router.push('/login')
      } else {
        alert(getUserFriendlyError(e))
      }
    }
  }, [router])

  const loadData = useCallback(async () => {
    setLoading(true)
    try {
      if (activeTab === 'users') {
        const usersData = await apiRequest<PaginatedResponse<User>>(
          `/api/v1/admin/users?page=${userPagination.page}&search=${searchQuery}`
        )
        setUsers(usersData.data || [])
        setUserPagination({
          page: usersData.page,
          total: usersData.total,
          totalPages: usersData.total_pages,
        })
      } else if (activeTab === 'categories') {
        const catsData = await apiRequest<VisaCategory[]>('/api/v1/visa-categories')
        setCategories(catsData)
      } else if (activeTab === 'applications') {
        const appsData = await apiRequest<PaginatedResponse<VisaApplication>>(
          `/api/v1/applications?page=${appPagination.page}&search=${searchQuery}`
        )
        setApplications(appsData.data || [])
        setAppPagination({
          page: appsData.page,
          total: appsData.total,
          totalPages: appsData.total_pages,
        })
      }
    } catch (err) {
      const e = err as AppError
      logger.error('Failed to load data:', e)
      alert(getUserFriendlyError(e))
    } finally {
      setLoading(false)
    }
  }, [activeTab, userPagination.page, appPagination.page, searchQuery])

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }

    loadUser()
  }, [router, loadUser])

  useEffect(() => {
    if (user) {
      if (user.role !== 'admin') {
        router.push('/dashboard')
        return
      }
      loadData()
    }
  }, [user, router, loadData])

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    router.push('/')
  }

  const handleUserRoleChange = async (userId: number, newRole: string) => {
    try {
      await apiRequest(`/api/v1/admin/users/${userId}`, {
        method: 'PUT',
        body: { role: newRole },
      })
      loadData()
    } catch (err) {
      const e = err as AppError
      logger.error('Failed to update user:', e)
      alert(getUserFriendlyError(e))
    }
  }

  const handleUserActiveToggle = async (userId: number, active: boolean) => {
    try {
      await apiRequest(`/api/v1/admin/users/${userId}`, {
        method: 'PUT',
        body: { active: !active },
      })
      loadData()
    } catch (err) {
      const e = err as AppError
      logger.error('Failed to update user:', e)
      alert(getUserFriendlyError(e))
    }
  }

  if (loading && !user) {
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
            <h1 className="text-xl font-bold">Admin Console</h1>
            <div className="flex items-center gap-4">
              <span className="text-sm text-gray-600">
                {user?.first_name} {user?.last_name} (Admin)
              </span>
              <Button variant="secondary" size="sm" onClick={handleLogout}>
                Logout
              </Button>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Tabs */}
        <div className="border-b border-gray-200 mb-6">
          <nav className="-mb-px flex space-x-8">
            {(['users', 'categories', 'applications'] as const).map((tab) => (
              <button
                key={tab}
                onClick={() => {
                  setActiveTab(tab)
                  setSearchQuery('')
                }}
                className={`py-4 px-1 border-b-2 font-medium text-sm capitalize ${
                  activeTab === tab
                    ? 'border-primary-500 text-primary-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                {tab}
              </button>
            ))}
          </nav>
        </div>

        {/* Search */}
        {activeTab !== 'categories' && (
          <div className="mb-6">
            <Input
              type="text"
              placeholder="Search..."
              value={searchQuery}
              onChange={(e) => {
                setSearchQuery(e.target.value)
                if (activeTab === 'users') {
                  setUserPagination({ ...userPagination, page: 1 })
                } else {
                  setAppPagination({ ...appPagination, page: 1 })
                }
              }}
              className="max-w-md"
            />
          </div>
        )}

        {/* Content */}
        {loading ? (
          <div className="text-center py-8">Loading...</div>
        ) : (
          <>
            {activeTab === 'users' && (
              <div>
                <div className="mb-4 flex justify-between items-center">
                  <h2 className="text-xl font-semibold">Users ({userPagination.total})</h2>
                </div>
                <div className="space-y-4">
                  {users.map((u) => (
                    <Card key={u.id}>
                      <div className="flex justify-between items-center">
                        <div>
                          <h3 className="font-semibold">
                            {u.first_name} {u.last_name}
                          </h3>
                          <p className="text-sm text-gray-600">{u.email}</p>
                          <p className="text-sm text-gray-500">Role: {u.role}</p>
                        </div>
                        <div className="flex gap-2">
                          <select
                            className="px-3 py-1 border rounded-lg text-sm"
                            value={u.role}
                            onChange={(e) => handleUserRoleChange(u.id, e.target.value)}
                          >
                            <option value="applicant">Applicant</option>
                            <option value="consultant">Consultant</option>
                            <option value="admin">Admin</option>
                          </select>
                          <Button
                            size="sm"
                            variant={u.active ? 'danger' : 'primary'}
                            onClick={() => handleUserActiveToggle(u.id, u.active)}
                          >
                            {u.active ? 'Deactivate' : 'Activate'}
                          </Button>
                        </div>
                      </div>
                    </Card>
                  ))}
                </div>
                {userPagination.totalPages > 1 && (
                  <div className="mt-6 flex justify-center gap-2">
                    <Button
                      variant="secondary"
                      size="sm"
                      disabled={userPagination.page === 1}
                      onClick={() =>
                        setUserPagination({ ...userPagination, page: userPagination.page - 1 })
                      }
                    >
                      Previous
                    </Button>
                    <span className="px-4 py-2 text-sm">
                      Page {userPagination.page} of {userPagination.totalPages}
                    </span>
                    <Button
                      variant="secondary"
                      size="sm"
                      disabled={userPagination.page === userPagination.totalPages}
                      onClick={() =>
                        setUserPagination({ ...userPagination, page: userPagination.page + 1 })
                      }
                    >
                      Next
                    </Button>
                  </div>
                )}
              </div>
            )}

            {activeTab === 'categories' && (
              <div>
                <div className="mb-4 flex justify-between items-center">
                  <h2 className="text-xl font-semibold">Visa Categories</h2>
                  <Button onClick={() => router.push('/admin/categories/new')}>
                    Add Category
                  </Button>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {categories.map((cat) => (
                    <Card key={cat.id}>
                      <h3 className="font-semibold text-lg mb-2">{cat.name}</h3>
                      <p className="text-sm text-gray-600 mb-2">{cat.description}</p>
                      <div className="flex justify-between items-center mt-4">
                        <span className="text-sm font-medium">{cat.country}</span>
                        <span className="text-sm font-semibold text-primary-600">
                          ${cat.price}
                        </span>
                      </div>
                      <div className="mt-2">
                        <Button
                          size="sm"
                          variant="secondary"
                          onClick={() => router.push(`/admin/categories/${cat.id}/edit`)}
                        >
                          Edit
                        </Button>
                      </div>
                    </Card>
                  ))}
                </div>
              </div>
            )}

            {activeTab === 'applications' && (
              <div>
                <div className="mb-4">
                  <h2 className="text-xl font-semibold">All Applications ({appPagination.total})</h2>
                </div>
                <div className="space-y-4">
                  {applications.map((app) => (
                    <Card key={app.id}>
                      <div className="flex justify-between items-start">
                        <div>
                          <h3 className="font-semibold">
                            {app.category?.name || 'Unknown Category'}
                          </h3>
                          <p className="text-sm text-gray-600">
                            User: {app.user?.first_name} {app.user?.last_name} ({app.user?.email})
                          </p>
                          <p className="text-sm text-gray-500">Status: {app.status}</p>
                        </div>
                        <Button size="sm" onClick={() => router.push(`/admin/applications/${app.id}`)}>
                          View
                        </Button>
                      </div>
                    </Card>
                  ))}
                </div>
                {appPagination.totalPages > 1 && (
                  <div className="mt-6 flex justify-center gap-2">
                    <Button
                      variant="secondary"
                      size="sm"
                      disabled={appPagination.page === 1}
                      onClick={() =>
                        setAppPagination({ ...appPagination, page: appPagination.page - 1 })
                      }
                    >
                      Previous
                    </Button>
                    <span className="px-4 py-2 text-sm">
                      Page {appPagination.page} of {appPagination.totalPages}
                    </span>
                    <Button
                      variant="secondary"
                      size="sm"
                      disabled={appPagination.page === appPagination.totalPages}
                      onClick={() =>
                        setAppPagination({ ...appPagination, page: appPagination.page + 1 })
                      }
                    >
                      Next
                    </Button>
                  </div>
                )}
              </div>
            )}
          </>
        )}
      </main>
    </div>
  )
}

