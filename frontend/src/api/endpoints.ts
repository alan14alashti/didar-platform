import { apiClient } from './client'
import type {
  AuthResponse,
  Booking,
  Brand,
  Category,
  Experience,
  SearchResult,
  User,
} from './types'

export const brandApi = {
  current: () => apiClient.get<Brand>('/brand').then((r) => r.data),
}

export const categoryApi = {
  list: () => apiClient.get<{ items: Category[] }>('/categories').then((r) => r.data.items),
}

export const authApi = {
  requestOtp: (phone: string) => apiClient.post('/auth/otp/request', { phone }),
  verifyOtp: (phone: string, code: string) =>
    apiClient.post<AuthResponse>('/auth/otp/verify', { phone, code }).then((r) => r.data),
  logout: (refreshToken: string) => apiClient.post('/auth/logout', { refresh_token: refreshToken }),
  me: () => apiClient.get<User>('/me').then((r) => r.data),
  updateMe: (payload: { full_name?: string; locale?: string }) =>
    apiClient.patch<User>('/me', payload).then((r) => r.data),
  becomeHost: (payload: { display_name?: string; bio?: string; city?: string }) =>
    apiClient.post<AuthResponse>('/me/become-host', payload).then((r) => r.data),
}

export interface SearchExperiencesParams {
  query?: string
  category?: string
  city?: string
  lat?: number
  lng?: number
  radius_km?: number
  page?: number
  page_size?: number
}

export const experienceApi = {
  search: (params: SearchExperiencesParams) =>
    apiClient.get<SearchResult>('/experiences', { params }).then((r) => r.data),
  getBySlug: (slug: string) => apiClient.get<Experience>(`/experiences/${slug}`).then((r) => r.data),
}

export interface ExperienceFormPayload {
  category_id?: string | null
  title: string
  summary?: string
  description?: string
  price_amount: number
  duration_minutes?: number
  max_guests?: number
  min_guests?: number
  city: string
  address?: string
  latitude: number
  longitude: number
  status?: 'draft' | 'published' | 'archived'
}

export const hostExperienceApi = {
  list: () => apiClient.get<{ items: Experience[] }>('/host/experiences').then((r) => r.data.items),
  get: (id: string) => apiClient.get<Experience>(`/host/experiences/${id}`).then((r) => r.data),
  create: (payload: ExperienceFormPayload) =>
    apiClient.post<Experience>('/host/experiences', payload).then((r) => r.data),
  update: (id: string, payload: ExperienceFormPayload) =>
    apiClient.patch<Experience>(`/host/experiences/${id}`, payload).then((r) => r.data),
  archive: (id: string) => apiClient.delete(`/host/experiences/${id}`),
  addPhoto: (id: string, payload: { url: string; alt?: string; is_cover?: boolean }) =>
    apiClient.post(`/host/experiences/${id}/photos`, payload),
  deletePhoto: (id: string, photoId: string) => apiClient.delete(`/host/experiences/${id}/photos/${photoId}`),
  addSlot: (id: string, payload: { start_at: string; end_at: string; capacity: number }) =>
    apiClient.post(`/host/experiences/${id}/slots`, payload),
  deleteSlot: (id: string, slotId: string) => apiClient.delete(`/host/experiences/${id}/slots/${slotId}`),
}

export const bookingApi = {
  create: (payload: { slot_id: string; guests_count: number; guest_note?: string }) =>
    apiClient.post<Booking>('/bookings', payload).then((r) => r.data),
  listMine: () => apiClient.get<{ items: Booking[] }>('/me/bookings').then((r) => r.data.items),
  cancelMine: (id: string) => apiClient.post(`/me/bookings/${id}/cancel`),
  listForHost: () => apiClient.get<{ items: Booking[] }>('/host/bookings').then((r) => r.data.items),
  updateStatusForHost: (id: string, payload: { status: string; host_note?: string }) =>
    apiClient.patch<Booking>(`/host/bookings/${id}`, payload).then((r) => r.data),
}
