'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { logger } from '@boltvisa/utils'
import { apiRequest } from '@boltvisa/utils'

import type { AppError } from '@boltvisa/types'
import { Notification } from '@boltvisa/types'
import { Button, Card } from '@boltvisa/ui'
import { getUserFriendlyError } from '@/lib/errorHelpers'

export default function NotificationsPage() {
  const router = useRouter()
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [loading, setLoading] = useState(true)
  const [unreadCount, setUnreadCount] = useState(0)

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }

    loadNotifications()
  }, [router])

  const loadNotifications = async () => {
    try {
      const data = await apiRequest<Notification[]>('/api/v1/notifications')
      setNotifications(data)
      setUnreadCount(data.filter((n) => !n.read).length)
    } catch (err) {
      const e = err as AppError
      logger.error('Failed to load notifications:', e)
      // Silently fail for notifications - user can retry by refreshing
    } finally {
      setLoading(false)
    }
  }

  const markAsRead = async (id: number) => {
    try {
      await apiRequest(`/api/v1/notifications/${id}/read`, {
        method: 'PUT',
      })
      loadNotifications()
    } catch (err) {
      const e = err as AppError
      logger.error('Failed to mark notification as read:', e)
      alert(getUserFriendlyError(e))
    }
  }

  const markAllAsRead = async () => {
    try {
      await apiRequest('/api/v1/notifications/read-all', {
        method: 'PUT',
      })
      loadNotifications()
    } catch (err) {
      const e = err as AppError
      logger.error('Failed to mark all as read:', e)
      alert(getUserFriendlyError(e))
    }
  }

  const getTypeIcon = (type: string) => {
    const icons: Record<string, string> = {
      application_update: '📋',
      document_request: '📄',
      payment: '💳',
      system: '🔔',
    }
    return icons[type] || '🔔'
  }

  const getTypeColor = (type: string) => {
    const colors: Record<string, string> = {
      application_update: 'bg-blue-100 text-blue-800',
      document_request: 'bg-yellow-100 text-yellow-800',
      payment: 'bg-green-100 text-green-800',
      system: 'bg-gray-100 text-gray-800',
    }
    return colors[type] || 'bg-gray-100 text-gray-800'
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
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-6 flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold">Notifications</h1>
            {unreadCount > 0 && (
              <p className="text-sm text-gray-600 mt-1">
                {unreadCount} unread notification{unreadCount !== 1 ? 's' : ''}
              </p>
            )}
          </div>
          {unreadCount > 0 && (
            <Button variant="secondary" onClick={markAllAsRead}>
              Mark All as Read
            </Button>
          )}
        </div>

        {notifications.length === 0 ? (
          <Card>
            <div className="text-center py-12">
              <div className="text-6xl mb-4">🔔</div>
              <h2 className="text-xl font-semibold mb-2">No notifications</h2>
              <p className="text-gray-600">You&apos;re all caught up!</p>
            </div>
          </Card>
        ) : (
          <div className="space-y-4">
            {notifications.map((notification) => (
              <Card
                key={notification.id}
                className={`cursor-pointer transition-all hover:shadow-lg ${
                  !notification.read ? 'border-l-4 border-primary-500 bg-primary-50' : ''
                }`}
                onClick={() => {
                  if (!notification.read) {
                    markAsRead(notification.id)
                  }
                }}
              >
                <div className="flex items-start gap-4">
                  <div className="text-3xl">{getTypeIcon(notification.type)}</div>
                  <div className="flex-1">
                    <div className="flex items-start justify-between mb-2">
                      <h3 className="font-semibold text-lg">{notification.title}</h3>
                      {!notification.read && (
                        <span className="px-2 py-1 bg-primary-600 text-white text-xs rounded-full">
                          New
                        </span>
                      )}
                    </div>
                    <p className="text-gray-700 mb-2">{notification.message}</p>
                    <div className="flex items-center gap-4 text-sm text-gray-500">
                      <span
                        className={`px-2 py-1 rounded ${getTypeColor(notification.type)}`}
                      >
                        {notification.type.replace('_', ' ')}
                      </span>
                      <span>
                        {new Date(notification.created_at).toLocaleDateString('en-US', {
                          year: 'numeric',
                          month: 'long',
                          day: 'numeric',
                          hour: '2-digit',
                          minute: '2-digit',
                        })}
                      </span>
                    </div>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

