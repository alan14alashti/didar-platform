<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Experience } from '../api/types'
import { useFormatters } from '../composables/useFormatters'
import AppIcon from './AppIcon.vue'

const props = defineProps<{
  experience: Experience
  active?: boolean
}>()

const { t } = useI18n()
const { formatPrice, formatDistance } = useFormatters()
</script>

<template>
  <router-link
    :to="{ name: 'experience-detail', params: { slug: props.experience.slug } }"
    class="group flex gap-3 rounded-xl p-2 transition-colors duration-200"
    :class="active ? 'bg-[#ebebeb]' : 'hover:bg-[#f7f7f7]'"
  >
    <div class="relative h-[4.5rem] w-[4.5rem] shrink-0 overflow-hidden rounded-lg bg-[#f7f7f7]">
      <img
        v-if="props.experience.cover_photo_url"
        :src="props.experience.cover_photo_url"
        :alt="props.experience.title"
        class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
        loading="lazy"
      />
    </div>

    <div class="min-w-0 flex-1 py-0.5">
      <div class="flex items-start justify-between gap-2">
        <h3 class="line-clamp-2 text-[15px] font-semibold leading-snug text-[#222222]">
          {{ props.experience.title }}
        </h3>
        <span
          v-if="props.experience.rating_count > 0"
          class="inline-flex shrink-0 items-center gap-0.5 text-sm text-[#222222]"
        >
          <AppIcon icon="heroicons-outline:star" size="xs" icon-class="text-[#222222]" />
          <span class="font-normal">{{ props.experience.rating_avg.toFixed(2) }}</span>
        </span>
      </div>

      <p class="mt-0.5 truncate text-sm text-[#717171]">
        {{ props.experience.city }}
        <span v-if="props.experience.distance_meters != null">
          · {{ formatDistance(props.experience.distance_meters) }}
        </span>
      </p>

      <p class="mt-0.5 truncate text-sm text-[#717171]">
        {{ t('experience.durationMinutes', { n: props.experience.duration_minutes }) }}
      </p>

      <p class="flex flex-wrap items-baseline gap-x-1 mt-1 text-sm text-[#222222]">
        <span class="font-semibold tabular-nums">{{ formatPrice(props.experience.price_amount) }}</span>
        <span class="font-normal text-[#717171]">{{ t('experience.perPerson') }}</span>
      </p>
    </div>
  </router-link>
</template>
