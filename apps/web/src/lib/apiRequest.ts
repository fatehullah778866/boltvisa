// Re-export AppError from shared types for convenience
// Prefer importing directly from @boltvisa/types in new code
import type { AppError } from '@boltvisa/types'
export type { AppError } from '@boltvisa/types'

const asAppError = (status: number, message: string, details?: unknown, cause?: unknown): AppError =>
  ({ status, message, details, cause })

// Best-effort JSON parse (returns raw text if not JSON)
async function safeJson<T>(res: Response): Promise<T> {
  const text = await res.text()
  try { return JSON.parse(text) as T } catch { return text as unknown as T }
}

// Type guard to extract error message from unknown response body
function extractErrorMessage(body: unknown): string | undefined {
  if (body && typeof body === 'object') {
    if ('error' in body && typeof body.error === 'string') {
      return body.error
    }
    if ('message' in body && typeof body.message === 'string') {
      return body.message
    }
    if (Array.isArray(body) && body.length > 0 && body[0] && typeof body[0] === 'object' && 'message' in body[0]) {
      return String(body[0].message)
    }
  }
  if (typeof body === 'string') {
    return body
  }
  return undefined
}

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export async function apiRequest<T>(input: string, init?: RequestInit): Promise<T> {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(init?.headers as Record<string, string>),
  }

  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  // URL construction: Use Next.js proxy for /api paths, direct URL for absolute URLs, or construct from API_URL
  // Next.js rewrites /api/* to ${API_URL}/api/* automatically (configured in next.config.js)
  const url = input.startsWith('/api') 
    ? input  // Use relative path - Next.js will proxy it
    : input.startsWith('http') 
      ? input  // Use absolute URL as-is
      : `${API_URL}${input}`  // Construct full URL for other paths

  // Log request for debugging (only in development)
  if (typeof window !== 'undefined' && process.env.NODE_ENV === 'development') {
    console.debug('[apiRequest]', { method: init?.method || 'GET', url, input })
  }

  try {
    const controller = new AbortController()
    // 10 second timeout for faster feedback and better UX
    const timeoutId = setTimeout(() => controller.abort(), 10000)

    const res = await fetch(url, {
      ...init,
      headers: headers as HeadersInit,
      cache: 'no-store',
      signal: controller.signal,
    })

    clearTimeout(timeoutId)

    if (!res.ok) {
      const body = await safeJson<unknown>(res)
      // Accept common shapes: {error}, {message}, array of errors, plain string
      const message = extractErrorMessage(body) || `HTTP ${res.status} ${res.statusText}`
      throw asAppError(res.status, message, body)
    }

    return safeJson<T>(res)
  } catch (error) {
    // Handle network errors (backend not running, CORS, etc.)
    if (error instanceof TypeError && error.message.includes('fetch')) {
      // Check if it's a CORS error vs network error
      const isCorsError = error.message.includes('CORS') || 
                         error.message.includes('Failed to fetch') ||
                         (error.message.includes('NetworkError') && url.startsWith('http'))
      
      const errorMessage = isCorsError
        ? `CORS error: Cannot connect to backend API. This may indicate a CORS configuration issue. Please check that the backend allows requests from ${typeof window !== 'undefined' ? window.location.origin : 'the frontend'}.`
        : `Cannot connect to backend API at ${API_URL}. Please ensure the backend server is running.`
      
      throw asAppError(
        0,
        errorMessage,
        { 
          url, 
          apiUrl: API_URL,
          isCorsError,
          originalError: error instanceof Error ? error.message : String(error)
        },
        error
      )
    }
    if (error instanceof Error && error.name === 'AbortError') {
      // Distinguish between timeout and other abort reasons
      const isTimeout = error.message.includes('aborted') || error.message.includes('timeout')
      const errorMessage = isTimeout
        ? `Request timed out after 10 seconds. The backend server may not be running or is unreachable. Please check: 1) Backend is running on ${API_URL}, 2) Network connection is working, 3) No firewall blocking the connection.`
        : `Request was aborted. Please try again.`
      
      throw asAppError(
        0,
        errorMessage,
        { input, url, timeout: 10000, isTimeout },
        error
      )
    }
    // Re-throw AppError as-is
    if (error && typeof error === 'object' && 'status' in error) {
      throw error
    }
    // Wrap unknown errors
    throw asAppError(0, error instanceof Error ? error.message : 'Unknown error occurred', null, error)
  }
}

