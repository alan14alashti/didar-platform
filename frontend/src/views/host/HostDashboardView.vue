<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { bookingApi, hostExperienceApi } from '../../api/endpoints'
import type { Booking, Experience } from '../../api/types'
import { useAuthStore } from '../../stores/auth'
import LoadingSpinner from '../../components/LoadingSpinner.vue'
import AppIcon from '../../components/AppIcon.vue'

const { t } = useI18n()
const authStore = useAuthStore()

const experiences = ref<Experience[]>([])
const bookings = ref<Booking[]>([])
const isLoading = ref(true)

const publishedCount = computed(() => experiences.value.filter((e) => e.status === 'published').length)
const pendingBookings = computed(() => bookings.value.filter((b) => b.status === 'pending'))

onMounted(async () => {
  try {
    const [exps, books] = await Promise.all([hostExperienceApi.list(), bookingApi.listForHost()])
    experiences.value = exps
    bookings.value = books
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-gray-900">
      {{ t('host.dashboardTitle') }} <span class="text-gray-400">— {{ authStore.user?.full_name || authStore.user?.phone }}</span>
    </h1>

    <LoadingSpinner v-if="isLoading" :label="t('common.loading')" />

    <div v-else class="mt-6 grid gap-4 sm:grid-cols-3">
      <router-link :to="{ name: 'host-experiences' }" class="surface-card group flex items-center gap-4 p-5 hover:-translate-y-1">
        <span class="grid h-11 w-11 shrink-0 place-items-center rounded-2xl bg-brand-primary/10 text-brand-primary">
          <AppIcon icon="heroicons-outline:map" size="lg" />
        </span>
        <div>
          <p class="text-2xl font-extrabold text-gray-900">{{ experiences.length }}</p>
          <p class="mt-0.5 text-sm text-gray-500">{{ t('host.myExperiences') }}</p>
        </div>
      </router-link>
      <div class="surface-card flex items-center gap-4 p-5">
        <span class="grid h-11 w-11 shrink-0 place-items-center rounded-2xl bg-green-100 text-green-700">
          <AppIcon icon="heroicons-outline:check-circle" size="lg" />
        </span>
        <div>
          <p class="text-2xl font-extrabold text-gray-900">{{ publishedCount }}</p>
          <p class="mt-0.5 text-sm text-gray-500">{{ t('host.statusPublished') }}</p>
        </div>
      </div>
      <router-link :to="{ name: 'host-bookings' }" class="surface-card group flex items-center gap-4 p-5 hover:-translate-y-1">
        <span class="grid h-11 w-11 shrink-0 place-items-center rounded-2xl bg-amber-100 text-amber-700">
          <AppIcon icon="heroicons-outline:envelope" size="lg" />
        </span>
        <div>
          <p class="text-2xl font-extrabold text-gray-900">{{ pendingBookings.length }}</p>
          <p class="mt-0.5 text-sm text-gray-500">{{ t('booking.statusPending') }}</p>
        </div>
      </router-link>
    </div>

    <div class="mt-8">
      <router-link :to="{ name: 'host-experience-new' }" class="btn-primary">
        + {{ t('host.newExperience') }}
      </router-link>
    </div>
  </div>
</template>
