// Mirrors the JSON shapes returned by the Go backend (internal/*/model.go).

export interface Brand {
  id: string
  slug: string
  name: string
  tagline: string | null
  description: string | null
  logo_url: string | null
  favicon_url: string | null
  primary_color: string
  secondary_color: string
  accent_color: string
  default_locale: string
  supported_locales: string[]
  default_currency: string
  contact_phone: string | null
  contact_email: string | null
  domain: string | null
}

export interface Category {
  id: string
  slug: string
  title_fa: string
  title_en: string
  icon: string | null
  sort_order: number
}

export interface User {
  id: string
  email: string | null
  phone: string | null
  full_name: string | null
  avatar_url: string | null
  locale: string
  is_host: boolean
  is_admin: boolean
  status: string
  last_login_at: string | null
  created_at: string
  updated_at: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  expires_at: string
  user: User
  is_new_user?: boolean
}

export interface Photo {
  id: string
  url: string
  alt: string | null
  sort_order: number
  is_cover: boolean
}

export interface Slot {
  id: string
  experience_id: string
  start_at: string
  end_at: string
  capacity: number
  booked_count: number
  status: string
}

export interface HostSummary {
  id: string
  full_name: string | null
  avatar_url: string | null
  display_name: string | null
  bio: string | null
  is_verified: boolean
  rating_avg: number
  rating_count: number
}

export interface Experience {
  id: string
  brand_id: string
  host_id: string
  category_id: string | null
  category_slug?: string | null
  slug: string
  title: string
  summary: string | null
  description: string | null
  language: string
  price_amount: number
  currency: string
  duration_minutes: number
  max_guests: number
  min_guests: number
  city: string
  address: string | null
  latitude: number
  longitude: number
  status: 'draft' | 'published' | 'archived'
  rating_avg: number
  rating_count: number
  distance_meters?: number | null
  created_at: string
  updated_at: string
  cover_photo_url?: string | null
  photos?: Photo[]
  slots?: Slot[]
  host?: HostSummary
}

export interface SearchResult {
  items: Experience[]
  page: number
  page_size: number
  total_count: number
}

export interface Booking {
  id: string
  experience_id: string
  experience_title?: string
  experience_slug?: string
  cover_photo_url?: string | null
  slot_id: string
  slot_start_at?: string
  slot_end_at?: string
  guest_id: string
  guest_name?: string | null
  guest_phone?: string
  guests_count: number
  total_amount: number
  currency: string
  status: 'pending' | 'confirmed' | 'cancelled' | 'completed'
  guest_note: string | null
  host_note: string | null
  created_at: string
  updated_at: string
}

export interface ApiErrorPayload {
  error: {
    code: string
    message: string
  }
}
