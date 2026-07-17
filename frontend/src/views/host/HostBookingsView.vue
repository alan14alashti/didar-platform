<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { bookingApi } from '../../api/endpoints'
import type { Booking } from '../../api/types'
import { useFormatters } from '../../composables/useFormatters'
import LoadingSpinner from '../../components/LoadingSpinner.vue'
import EmptyState from '../../components/EmptyState.vue'

const { t } = useI18n()
const { formatPrice, formatDateTime } = useFormatters()

const bookings = ref<Booking[]>([])
const isLoading = ref(true)
const updatingId = ref<string | null>(null)

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
    bookings.value = await bookingApi.listForHost()
  } finally {
    isLoading.value = false
  }
}

async function updateStatus(id: string, status: string) {
  updatingId.value = id
  try {
    await bookingApi.updateStatusForHost(id, { status })
    await load()
  } finally {
    updatingId.value = null
  }
}

onMounted(load)
</script>

<template>
  <div>
    <h1 class="text-xl font-bold text-gray-900">{{ t('host.bookingsTitle') }}</h1>

    <LoadingSpinner v-if="isLoading" :label="t('common.loading')" />

    <EmptyState v-else-if="bookings.length === 0" :title="t('host.noBookingsYet')" class="mt-6" />

    <div v-else class="mt-6 space-y-3">
      <div v-for="b in bookings" :key="b.id" class="rounded-2xl border border-gray-100 bg-white p-4 shadow-sm">
        <div class="flex items-start justify-between gap-2">
          <div>
            <p class="font-semibold text-gray-900">{{ b.experience_title }}</p>
            <p class="mt-0.5 text-sm text-gray-500">{{ formatDateTime(b.slot_start_at || '') }}</p>
            <p class="mt-0.5 text-sm text-gray-500">
              {{ t('host.guestLabel') }}: {{ b.guest_name || b.guest_phone }} · {{ t('common.guestsCount', { n: b.guests_count }) }}
            </p>
          </div>
          <span class="shrink-0 rounded-full px-2.5 py-1 text-xs font-medium" :class="statusStyles[b.status]">
            {{ statusLabel(b.status) }}
          </span>
        </div>

        <div class="mt-3 flex items-center justify-between border-t border-gray-100 pt-3">
          <span class="text-sm font-bold text-gray-900 tabular-nums">{{ formatPrice(b.total_amount) }}</span>
          <div v-if="b.status === 'pending' || b.status === 'confirmed'" class="flex gap-2">
            <button
              v-if="b.status === 'pending'"
              type="button"
              class="rounded-full bg-green-600 px-3 py-1.5 text-xs font-semibold text-white disabled:opacity-50"
              :disabled="updatingId === b.id"
              @click="updateStatus(b.id, 'confirmed')"
            >
              {{ t('host.confirmBooking') }}
            </button>
            <button
              v-if="b.status === 'confirmed'"
              type="button"
              class="rounded-full bg-blue-600 px-3 py-1.5 text-xs font-semibold text-white disabled:opacity-50"
              :disabled="updatingId === b.id"
              @click="updateStatus(b.id, 'completed')"
            >
              {{ t('host.markCompleted') }}
            </button>
            <button
              type="button"
              class="rounded-full bg-red-50 px-3 py-1.5 text-xs font-semibold text-red-600 disabled:opacity-50"
              :disabled="updatingId === b.id"
              @click="updateStatus(b.id, 'cancelled')"
            >
              {{ t('host.cancelBooking') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
