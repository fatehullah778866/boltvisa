'use client'

import useSWR from 'swr'
import { authedFetch } from '@/lib/authClient'

export interface DashboardData {
  applications: number
  pending: number
  approved: number
  rejected: number
  notifications: number
}

interface DashboardResponse {
  data: DashboardData
}

export function useDashboard() {
  const { data, error, isLoading, mutate } = useSWR<DashboardResponse>(
    '/api/v1/dashboard',
    (key: string) => authedFetch<DashboardResponse>(key)
  )

  return {
    data: data?.data,
    isLoading,
    error,
    reload: mutate,
  }
}

