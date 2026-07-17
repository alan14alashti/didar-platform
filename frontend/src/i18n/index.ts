import { createI18n } from 'vue-i18n'

import fa from './locales/fa.json'
import en from './locales/en.json'

export const RTL_LOCALES = new Set(['fa'])

const STORAGE_KEY = 'didar_locale'

export function getInitialLocale(): string {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored === 'fa' || stored === 'en') return stored
  return navigator.language?.toLowerCase().startsWith('en') ? 'en' : 'fa'
}

export const i18n = createI18n({
  legacy: false,
  locale: getInitialLocale(),
  fallbackLocale: 'fa',
  messages: { fa, en },
})

export function applyDocumentDirection(locale: string) {
  const dir = RTL_LOCALES.has(locale) ? 'rtl' : 'ltr'
  document.documentElement.dir = dir
  document.documentElement.lang = locale
}

export function setLocale(locale: string) {
  ;(i18n.global.locale as unknown as { value: string }).value = locale
  localStorage.setItem(STORAGE_KEY, locale)
  applyDocumentDirection(locale)
}
