'use client'

const apiBase = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'

let isRefreshing = false
let waiters: Array<() => void> = []

async function refreshToken(): Promise<void> {
  if (isRefreshing) {
    await new Promise<void>((resolve) => {
      waiters.push(resolve)
    })
    return
  }

  isRefreshing = true
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      throw new Error('refresh_failed')
    }

    const res = await fetch(`${apiBase}/api/v1/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      credentials: 'include',
    })

    if (!res.ok) {
      throw new Error('refresh_failed')
    }

    const data = await res.json()
    if (data.token) {
      localStorage.setItem('token', data.token)
      if (data.user) {
        localStorage.setItem('user', JSON.stringify(data.user))
      }
    }
  } finally {
    isRefreshing = false
    waiters.forEach((w) => w())
    waiters = []
  }
}

export async function authedFetch<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = localStorage.getItem('token')

  const attempt = async (): Promise<T> => {
    const url = path.startsWith('/api') ? `${apiBase}${path}` : path
    const res = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...(options.headers || {}),
      },
      credentials: 'include',
    })

    if (res.status === 401) {
      throw new Error('unauthorized')
    }

    if (!res.ok) {
      const text = await res.text()
      throw new Error(`${res.status}:${text}`)
    }

    return res.json()
  }

  try {
    return await attempt()
  } catch (e: any) {
    if (e?.message === 'unauthorized') {
      await refreshToken()
      const newToken = localStorage.getItem('token')
      
      const url = path.startsWith('/api') ? `${apiBase}${path}` : path
      const res = await fetch(url, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...(newToken ? { Authorization: `Bearer ${newToken}` } : {}),
          ...(options.headers || {}),
        },
        credentials: 'include',
      })

      if (res.status === 401) {
        throw new Error('reauth_required')
      }

      if (!res.ok) {
        const text = await res.text()
        throw new Error(`${res.status}:${text}`)
      }

      return res.json()
    }
    throw e
  }
}

