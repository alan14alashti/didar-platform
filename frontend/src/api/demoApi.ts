import { demoBrand, demoCategories, demoExperiences } from './demoData'
import type { Booking, Experience, SearchResult, User } from './types'

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

const delay = (ms = 120) => new Promise((resolve) => setTimeout(resolve, ms))

function haversineMeters(lat1: number, lng1: number, lat2: number, lng2: number) {
  const toRad = (deg: number) => (deg * Math.PI) / 180
  const dLat = toRad(lat2 - lat1)
  const dLng = toRad(lng2 - lng1)
  const a =
    Math.sin(dLat / 2) ** 2 +
    Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) * Math.sin(dLng / 2) ** 2
  return 6371000 * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))
}

function listItem(exp: Experience): Experience {
  const { photos, slots, host, ...rest } = exp
  void photos
  void slots
  void host
  return rest
}

export const demoBrandApi = {
  current: async () => {
    await delay()
    return demoBrand
  },
}

export const demoCategoryApi = {
  list: async () => {
    await delay()
    return demoCategories
  },
}

export const demoExperienceApi = {
  search: async (params: SearchExperiencesParams = {}): Promise<SearchResult> => {
    await delay()
    const page = params.page ?? 1
    const pageSize = params.page_size ?? 9
    let items = demoExperiences.map(listItem)

    if (params.category) {
      items = items.filter((item) => item.category_slug === params.category)
    }
    if (params.query) {
      const q = params.query.trim().toLowerCase()
      items = items.filter(
        (item) =>
          item.title.toLowerCase().includes(q) ||
          (item.summary?.toLowerCase().includes(q) ?? false) ||
          item.city.toLowerCase().includes(q),
      )
    }

    if (params.lat != null && params.lng != null) {
      items = items.map((item) => ({
        ...item,
        distance_meters: haversineMeters(params.lat!, params.lng!, item.latitude, item.longitude),
      }))
      if (params.radius_km && params.radius_km > 0) {
        items = items.filter((item) => (item.distance_meters ?? 0) <= params.radius_km! * 1000)
      }
      items.sort((a, b) => (a.distance_meters ?? 0) - (b.distance_meters ?? 0))
    }

    const total_count = items.length
    const start = (page - 1) * pageSize
    return {
      items: items.slice(start, start + pageSize),
      page,
      page_size: pageSize,
      total_count,
    }
  },
  getBySlug: async (slug: string) => {
    await delay()
    const exp = demoExperiences.find((item) => item.slug === slug)
    if (!exp) {
      const err = new Error('Not found') as Error & { response?: { status: number } }
      err.response = { status: 404 }
      throw err
    }
    return structuredClone(exp)
  },
}

export const demoAuthApi = {
  requestOtp: async () => {
    await delay()
  },
  verifyOtp: async () => {
    await delay()
    const user: User = {
      id: 'demo-user',
      phone: '+989000000000',
      full_name: 'کاربر دمو',
      avatar_url: null,
      locale: 'fa',
      is_host: false,
      is_admin: false,
      status: 'active',
      last_login_at: new Date().toISOString(),
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }
    return {
      access_token: 'demo-access-token',
      refresh_token: 'demo-refresh-token',
      expires_at: new Date(Date.now() + 86_400_000).toISOString(),
      user,
      is_new_user: false,
    }
  },
  logout: async () => {
    await delay()
  },
  me: async () => {
    await delay()
    throw new Error('Unauthorized')
  },
  updateMe: async () => {
    await delay()
    throw new Error('Unauthorized')
  },
  becomeHost: async () => {
    await delay()
    throw new Error('Demo mode')
  },
}

export const demoHostExperienceApi = {
  list: async () => [],
  get: async () => {
    throw new Error('Demo mode')
  },
  create: async () => {
    throw new Error('Demo mode')
  },
  update: async () => {
    throw new Error('Demo mode')
  },
  archive: async () => {
    throw new Error('Demo mode')
  },
  addPhoto: async () => {
    throw new Error('Demo mode')
  },
  deletePhoto: async () => {
    throw new Error('Demo mode')
  },
  addSlot: async () => {
    throw new Error('Demo mode')
  },
  deleteSlot: async () => {
    throw new Error('Demo mode')
  },
}

export const demoBookingApi = {
  create: async () => {
    throw new Error('Demo mode')
  },
  listMine: async (): Promise<Booking[]> => [],
  cancelMine: async () => {
    throw new Error('Demo mode')
  },
  listForHost: async (): Promise<Booking[]> => [],
  updateStatusForHost: async () => {
    throw new Error('Demo mode')
  },
}
