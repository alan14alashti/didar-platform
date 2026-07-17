import { ref } from 'vue'
import { defineStore } from 'pinia'

import { brandApi } from '../api/endpoints'
import type { Brand } from '../api/types'

// Applies the resolved white-label brand's palette as CSS variables and
// exposes its metadata (name, logo, locales, currency) to the rest of the
// app. This is the entire mechanism behind "re-skin without code changes".
function applyTheme(brand: Brand) {
  const root = document.documentElement
  root.style.setProperty('--color-primary', brand.primary_color)
  root.style.setProperty('--color-secondary', brand.secondary_color)
  root.style.setProperty('--color-accent', brand.accent_color)
  document.title = brand.name
  if (brand.favicon_url) {
    const link = document.querySelector<HTMLLinkElement>('link[rel="icon"]')
    if (link) link.href = brand.favicon_url
  }
}

export const useBrandStore = defineStore('brand', () => {
  const brand = ref<Brand | null>(null)
  const isLoaded = ref(false)

  async function load() {
    try {
      brand.value = await brandApi.current()
      applyTheme(brand.value)
    } finally {
      isLoaded.value = true
    }
  }

  return { brand, isLoaded, load }
})
