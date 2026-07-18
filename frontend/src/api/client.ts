import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'

import type { AuthResponse } from './types'
import { tokenStorage } from './tokenStorage'

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  // Skips the localtunnel.me interstitial "click to continue" HTML page that
  // would otherwise be returned instead of JSON when the API is exposed via
  // a `loca.lt` quick tunnel (used for the public GitHub Pages demo). Ignored
  // by any other host.
  headers: { 'Bypass-Tunnel-Reminder': 'true' },
})

apiClient.interceptors.request.use((config) => {
  const token = tokenStorage.getAccessToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

let refreshPromise: Promise<string | null> | null = null

async function refreshAccessToken(): Promise<string | null> {
  const refreshToken = tokenStorage.getRefreshToken()
  if (!refreshToken) return null

  try {
    const { data } = await axios.post<AuthResponse>(
      `${API_BASE_URL}/auth/refresh`,
      { refresh_token: refreshToken },
      { headers: { 'Bypass-Tunnel-Reminder': 'true' } },
    )
    tokenStorage.setTokens(data.access_token, data.refresh_token)
    return data.access_token
  } catch {
    tokenStorage.clear()
    return null
  }
}

interface RetryableConfig extends InternalAxiosRequestConfig {
  _retried?: boolean
}

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const config = error.config as RetryableConfig | undefined

    if (error.response?.status === 401 && config && !config._retried && !config.url?.includes('/auth/')) {
      config._retried = true
      refreshPromise ??= refreshAccessToken().finally(() => {
        refreshPromise = null
      })
      const newToken = await refreshPromise
      if (newToken) {
        config.headers = config.headers ?? {}
        config.headers.Authorization = `Bearer ${newToken}`
        return apiClient(config)
      }
      tokenStorage.clear()
      window.location.assign(`${import.meta.env.BASE_URL}login`.replace(/\/{2,}/g, '/'))
    }

    return Promise.reject(error)
  },
)

export function extractErrorMessage(error: unknown, fallback: string): string {
  if (axios.isAxiosError(error)) {
    const message = (error.response?.data as { error?: { message?: string } } | undefined)?.error?.message
    if (message) return message
  }
  return fallback
}

export function extractErrorCode(error: unknown): string | null {
  if (axios.isAxiosError(error)) {
    return (error.response?.data as { error?: { code?: string } } | undefined)?.error?.code ?? null
  }
  return null
}
