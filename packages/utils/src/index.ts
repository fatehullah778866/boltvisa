// Re-export the shared API client so apps can import from
// "@boltvisa/utils" or "@boltvisa/utils/apiRequest"
export * from './apiRequest';


export function formatDate(date: string | Date): string {
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

export function formatCurrency(amount: number, currency = 'USD'): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(amount)
}

export function getInitials(firstName: string, lastName: string): string {
  return `${firstName.charAt(0)}${lastName.charAt(0)}`.toUpperCase()
}

// Export logger
export { logger, type Logger } from './logger'

// Export HTTP utilities
export { safeJson, getJson } from './http'

