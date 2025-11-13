// packages/types/src/index.ts

// ===== Core users/roles =====
export type UserRole = 'admin' | 'consultant' | 'applicant';

export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  role: UserRole;
  active: boolean;
  created_at: string; // ISO 8601
  updated_at: string; // ISO 8601
}

// ===== Auth DTOs =====
export interface LoginRequest { email: string; password: string; }
export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
}
export interface AuthResponse {
  token: string;
  user: User;
  expires_at: string; // ISO 8601
}
export interface ApiError { error: string; }

// Some files import this shape from @boltvisa/types
export interface AppError {
  status: number;
  message: string;
  details?: unknown;
  cause?: unknown;
}

// ===== Visa domain =====
export type VisaStatus =
  | 'draft'
  | 'submitted'
  | 'in_review'
  | 'approved'
  | 'rejected'
  | 'cancelled';

export interface VisaCategory {
  id: number;
  name: string;
  description: string;
  country: string;
  duration: string;
  price: number;
  active: boolean;
  created_at: string;
  updated_at: string;
}

// Include populated relations that UI expects
export interface VisaApplication {
  id: number;
  user_id: number;
  consultant_id?: number | null;
  category_id: number;
  status: VisaStatus;
  passport_number?: string;
  date_of_birth?: string;
  nationality?: string;
  travel_date?: string;
  notes?: string;
  created_at: string;
  updated_at: string;

  // populated relations for UI
  user?: User;
  consultant?: User | null;
  category?: VisaCategory;
}

// ===== Documents =====
export interface Document {
  id: number;
  application_id: number;
  type: string;
  name: string;
  file_name: string;
  gcs_url: string;
  size: number;
  mime_type: string;
  uploaded_by: number;
  created_at: string;
  updated_at: string;
}

// ===== Payments =====
export type PaymentStatus = 'pending' | 'succeeded' | 'failed' | 'refunded';

export interface Payment {
  id: number;
  user_id: number;
  amount: number;
  status: PaymentStatus;
  method: 'stripe' | 'razorpay' | 'test';
  transaction: string;
  created_at: string;
  updated_at: string;
}

// ===== Notifications =====
export type NotificationType = 'email' | 'sms' | 'app';

export interface Notification {
  id: number;
  user_id: number;
  type: NotificationType;
  title: string;
  message: string;
  read: boolean;
  read_at?: string | null;
  created_at: string;
}

// ===== Audit Logs =====
export interface AuditLog {
  id: number;
  user_id?: number | null;
  action: string;
  resource: string; // e.g., "application:123"
  ip_address?: string | null;
  user_agent?: string | null;
  metadata?: Record<string, unknown> | null;
  created_at: string;
}
// …existing types above…

export interface CreateApplicationRequest {
  category_id: number;
  passport_number?: string;
  date_of_birth?: string; // ISO-8601
  nationality?: string;
  travel_date?: string;   // ISO-8601
  notes?: string;
  submit?: boolean;       // true -> submitted, false/undefined -> draft
}
