<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { hostExperienceApi } from '../../api/endpoints'
import type { Experience } from '../../api/types'
import { useFormatters } from '../../composables/useFormatters'
import LoadingSpinner from '../../components/LoadingSpinner.vue'
import EmptyState from '../../components/EmptyState.vue'

const { t } = useI18n()
const { formatPrice } = useFormatters()

const experiences = ref<Experience[]>([])
const isLoading = ref(true)

const statusStyles: Record<string, string> = {
  draft: 'bg-gray-100 text-gray-600',
  published: 'bg-green-100 text-green-700',
  archived: 'bg-red-50 text-red-500',
}

async function load() {
  isLoading.value = true
  try {
    experiences.value = await hostExperienceApi.list()
  } finally {
    isLoading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div>
    <div class="flex items-center justify-between">
      <h1 class="text-xl font-bold text-gray-900">{{ t('host.myExperiences') }}</h1>
      <router-link :to="{ name: 'host-experience-new' }" class="btn-primary !px-4 !py-2 text-sm">
        + {{ t('host.newExperience') }}
      </router-link>
    </div>

    <LoadingSpinner v-if="isLoading" :label="t('common.loading')" />

    <EmptyState v-else-if="experiences.length === 0" :title="t('host.noExperiencesYet')" class="mt-6">
      <router-link :to="{ name: 'host-experience-new' }" class="btn-primary mt-2">
        {{ t('host.createFirstExperience') }}
      </router-link>
    </EmptyState>

    <div v-else class="mt-6 space-y-3">
      <router-link
        v-for="exp in experiences"
        :key="exp.id"
        :to="{ name: 'host-experience-edit', params: { id: exp.id } }"
        class="surface-card flex items-center gap-4 p-4 hover:-translate-y-0.5"
      >
        <div class="h-16 w-16 shrink-0 overflow-hidden rounded-xl bg-gray-100">
          <img v-if="exp.cover_photo_url" :src="exp.cover_photo_url" class="h-full w-full object-cover" alt="" />
        </div>
        <div class="flex-1">
          <div class="flex items-center gap-2">
            <p class="font-semibold text-gray-900">{{ exp.title }}</p>
            <span class="rounded-full px-2 py-0.5 text-xs font-medium" :class="statusStyles[exp.status]">
              {{ t(`host.status${exp.status.charAt(0).toUpperCase()}${exp.status.slice(1)}`) }}
            </span>
          </div>
          <p class="mt-0.5 text-sm text-gray-500">{{ exp.city }} · {{ formatPrice(exp.price_amount) }}</p>
        </div>
        <span class="text-sm font-medium text-brand-primary">{{ t('host.manage') }}</span>
      </router-link>
    </div>
  </div>
</template>
