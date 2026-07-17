import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import { authApi } from '../api/endpoints'
import { tokenStorage } from '../api/tokenStorage'
import type { User } from '../api/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isBooting = ref(true)

  const isAuthenticated = computed(() => user.value !== null)
  const isHost = computed(() => user.value?.is_host ?? false)

  function setSession(newUser: User, accessToken: string, refreshToken: string) {
    user.value = newUser
    tokenStorage.setTokens(accessToken, refreshToken)
  }

  async function restoreSession() {
    isBooting.value = true
    if (tokenStorage.getAccessToken()) {
      try {
        user.value = await authApi.me()
      } catch {
        tokenStorage.clear()
        user.value = null
      }
    }
    isBooting.value = false
  }

  async function requestOtp(phone: string) {
    await authApi.requestOtp(phone)
  }

  async function verifyOtp(phone: string, code: string) {
    const res = await authApi.verifyOtp(phone, code)
    setSession(res.user, res.access_token, res.refresh_token)
    return res
  }

  async function becomeHost(payload: { display_name?: string; bio?: string; city?: string }) {
    const res = await authApi.becomeHost(payload)
    setSession(res.user, res.access_token, res.refresh_token)
    return res
  }

  async function updateProfile(payload: { full_name?: string; locale?: string }) {
    user.value = await authApi.updateMe(payload)
  }

  async function logout() {
    const refreshToken = tokenStorage.getRefreshToken()
    tokenStorage.clear()
    user.value = null
    if (refreshToken) {
      try {
        await authApi.logout(refreshToken)
      } catch {
        // best-effort; local session is already cleared
      }
    }
  }

  return {
    user,
    isBooting,
    isAuthenticated,
    isHost,
    restoreSession,
    requestOtp,
    verifyOtp,
    becomeHost,
    updateProfile,
    logout,
  }
})
