<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Experience } from '../api/types'
import { useFormatters } from '../composables/useFormatters'
import AppIcon from './AppIcon.vue'
import RatingStars from './RatingStars.vue'

defineProps<{
  experience: Experience
  style?: Record<string, string>
}>()

const emit = defineEmits<{ close: []; open: [] }>()
const { t } = useI18n()
const { formatPrice, formatDistance } = useFormatters()
</script>

<template>
  <div
    class="absolute z-20 w-[min(18rem,calc(100vw-2rem))] -translate-x-1/2 -translate-y-full overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-card-hover"
    :style="style"
  >
    <button
      type="button"
      class="absolute end-2 top-2 z-10 grid h-7 w-7 place-items-center rounded-full bg-white/90 text-gray-500 shadow-soft hover:text-gray-800"
      @click.stop="emit('close')"
    >
      <AppIcon icon="heroicons-outline:x-mark" size="sm" />
    </button>

    <button type="button" class="block w-full text-start" @click="emit('open')">
      <div class="aspect-[4/3] w-full overflow-hidden bg-gray-100">
        <img
          v-if="experience.cover_photo_url"
          :src="experience.cover_photo_url"
          :alt="experience.title"
          class="h-full w-full object-cover"
        />
      </div>
      <div class="space-y-1.5 p-3">
        <div class="flex items-start justify-between gap-2 pe-6">
          <h3 class="line-clamp-2 text-sm font-bold leading-snug text-gray-900">{{ experience.title }}</h3>
          <RatingStars v-if="experience.rating_count > 0" :rating="experience.rating_avg" size="sm" />
        </div>
        <p class="flex items-center gap-1 text-xs text-gray-500">
          <AppIcon icon="heroicons-outline:map-pin" size="xs" />
          {{ experience.city }}
          <span v-if="experience.distance_meters != null" class="text-brand-primary">
            · {{ formatDistance(experience.distance_meters) }}
          </span>
        </p>
        <p class="line-clamp-2 text-xs leading-5 text-gray-500">{{ experience.summary }}</p>
        <div class="flex items-center justify-between pt-1">
          <p class="text-sm font-extrabold text-gray-900 tabular-nums">
            {{ formatPrice(experience.price_amount) }}
            <span class="text-xs font-normal text-gray-400">/ {{ t('experience.perPerson') }}</span>
          </p>
          <span class="text-xs font-semibold text-brand-primary">{{ t('home.viewDetails') }} →</span>
        </div>
      </div>
    </button>

    <div class="absolute -bottom-2 start-1/2 h-4 w-4 -translate-x-1/2 rotate-45 border-b border-e border-gray-200 bg-white" />
  </div>
</template>
