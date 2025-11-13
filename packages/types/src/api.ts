// Shared API response types
export type ApiData<T> = { data: T }
export type ApiMessage = { message: string }

// Structured error type for consistent error handling
export type AppError = {
  status: number
  message: string
  details?: unknown
  cause?: unknown
}

// User API types
export type UserMe = { id: number; email: string; role: 'applicant' | 'consultant' | 'admin' }

// Auth API types
export type AuthLoginReq = { email: string; password: string }
export type AuthLoginRes = { token: string; user: UserMe } | ApiMessage

export type AuthRegisterReq = { email: string; password: string; first_name: string; last_name: string; phone_number?: string }
export type AuthRegisterRes = ApiMessage | { token: string; user: UserMe }

// Paginated response type
export interface PaginatedResponse<T> {
  data: T[]
  page: number
  page_size: number
  total: number
  total_pages: number
}

