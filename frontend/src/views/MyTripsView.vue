<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { bookingApi } from '../api/endpoints'
import type { Booking } from '../api/types'
import { useFormatters } from '../composables/useFormatters'
import LoadingSpinner from '../components/LoadingSpinner.vue'
import EmptyState from '../components/EmptyState.vue'

const { t } = useI18n()
const { formatPrice, formatDateTime } = useFormatters()

const bookings = ref<Booking[]>([])
const isLoading = ref(true)
const cancellingId = ref<string | null>(null)

const statusStyles: Record<string, string> = {
  pending: 'bg-amber-100 text-amber-700',
  confirmed: 'bg-green-100 text-green-700',
  cancelled: 'bg-gray-100 text-gray-500',
  completed: 'bg-blue-100 text-blue-700',
}

function statusLabel(status: string) {
  return t(`booking.status${status.charAt(0).toUpperCase()}${status.slice(1)}`)
}

async function load() {
  isLoading.value = true
  try {
    bookings.value = await bookingApi.listMine()
  } finally {
    isLoading.value = false
  }
}

async function cancelBooking(id: string) {
  if (!confirm(t('booking.cancelConfirm'))) return
  cancellingId.value = id
  try {
    await bookingApi.cancelMine(id)
    await load()
  } finally {
    cancellingId.value = null
  }
}

onMounted(load)
</script>

<template>
  <div class="mx-auto max-w-4xl px-4 py-10 sm:px-6">
    <h1 class="text-2xl font-bold text-gray-900">{{ t('booking.myTripsTitle') }}</h1>

    <LoadingSpinner v-if="isLoading" :label="t('common.loading')" />

    <EmptyState v-else-if="bookings.length === 0" :title="t('booking.noBookings')" class="mt-6">
      <router-link to="/" class="btn-primary mt-2">
        {{ t('booking.exploreExperiences') }}
      </router-link>
    </EmptyState>

    <div v-else class="mt-6 space-y-4">
      <div v-for="(b, idx) in bookings" :key="b.id" class="surface-card animate-fade-up flex gap-4 p-4" :style="{ animationDelay: `${idx * 60}ms` }">
        <router-link
          :to="{ name: 'experience-detail', params: { slug: b.experience_slug } }"
          class="h-24 w-24 shrink-0 overflow-hidden rounded-xl bg-gray-100"
        >
          <img v-if="b.cover_photo_url" :src="b.cover_photo_url" class="h-full w-full object-cover transition-transform duration-300 hover:scale-105" alt="" />
        </router-link>

        <div class="flex-1">
          <div class="flex items-start justify-between gap-2">
            <router-link :to="{ name: 'experience-detail', params: { slug: b.experience_slug } }" class="font-semibold text-gray-900 hover:text-brand-primary">
              {{ b.experience_title }}
            </router-link>
            <span class="shrink-0 rounded-full px-2.5 py-1 text-xs font-medium" :class="statusStyles[b.status]">
              {{ statusLabel(b.status) }}
            </span>
          </div>
          <p class="mt-1 text-sm text-gray-500">{{ formatDateTime(b.slot_start_at || '') }}</p>
          <p class="mt-1 text-sm text-gray-500">{{ t('common.guestsCount', { n: b.guests_count }) }}</p>
          <div class="mt-2 flex items-center justify-between">
            <span class="text-sm font-bold text-gray-900 tabular-nums">{{ formatPrice(b.total_amount) }}</span>
            <button
              v-if="b.status === 'pending' || b.status === 'confirmed'"
              type="button"
              class="text-sm font-medium text-red-600 hover:underline disabled:opacity-50"
              :disabled="cancellingId === b.id"
              @click="cancelBooking(b.id)"
            >
              {{ t('booking.cancelBooking') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
