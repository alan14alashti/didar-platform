import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

// Centralizes locale-aware money/date formatting so every view renders
// Persian (Jalali calendar, Toman) or English (Gregorian, as-is currency)
// consistently based on the active UI language.
export function useFormatters() {
  const { locale, t } = useI18n()

  const intlLocale = computed(() => (locale.value === 'fa' ? 'fa-IR' : 'en-US'))

  // Amounts are stored in Rial in the database; Iranians conventionally
  // talk about prices in Toman (1 Toman = 10 Rial).
  function formatPrice(amountRial: number): string {
    const toman = Math.round(amountRial / 10)
    const formatted = new Intl.NumberFormat(intlLocale.value).format(toman)
    return `${formatted} ${t('common.toman')}`
  }

  function formatDateTime(iso: string): string {
    if (!iso) return ''
    const date = new Date(iso)
    return new Intl.DateTimeFormat(intlLocale.value, {
      weekday: 'short',
      day: 'numeric',
      month: 'long',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date)
  }

  function formatDate(iso: string): string {
    if (!iso) return ''
    const date = new Date(iso)
    return new Intl.DateTimeFormat(intlLocale.value, {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
    }).format(date)
  }

  function formatNumber(n: number): string {
    return new Intl.NumberFormat(intlLocale.value).format(n)
  }

  return { formatPrice, formatDateTime, formatDate, formatNumber }
}
