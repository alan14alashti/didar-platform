// Small localStorage-backed token holder shared between the axios client
// (which needs to read/refresh tokens without importing the Pinia store,
// to avoid a circular dependency) and the auth store (source of truth for
// the rest of the app / UI reactivity).

const ACCESS_KEY = 'didar_access_token'
const REFRESH_KEY = 'didar_refresh_token'

export const tokenStorage = {
  getAccessToken(): string | null {
    return localStorage.getItem(ACCESS_KEY)
  },
  getRefreshToken(): string | null {
    return localStorage.getItem(REFRESH_KEY)
  },
  setTokens(accessToken: string, refreshToken: string) {
    localStorage.setItem(ACCESS_KEY, accessToken)
    localStorage.setItem(REFRESH_KEY, refreshToken)
  },
  clear() {
    localStorage.removeItem(ACCESS_KEY)
    localStorage.removeItem(REFRESH_KEY)
  },
}
