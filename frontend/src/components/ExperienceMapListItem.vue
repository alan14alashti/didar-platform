<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Experience } from '../api/types'
import { useFormatters } from '../composables/useFormatters'
import AppIcon from './AppIcon.vue'
import RatingStars from './RatingStars.vue'

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
    class="flex gap-3 rounded-2xl border p-2.5 transition-all duration-200"
    :class="active ? 'border-brand-primary bg-brand-primary/5 shadow-soft' : 'border-gray-100 bg-white hover:border-brand-primary/40 hover:shadow-soft'"
  >
    <div class="relative h-24 w-28 shrink-0 overflow-hidden rounded-xl bg-gray-100">
      <img
        v-if="props.experience.cover_photo_url"
        :src="props.experience.cover_photo_url"
        :alt="props.experience.title"
        class="h-full w-full object-cover"
        loading="lazy"
      />
    </div>

    <div class="min-w-0 flex-1 py-0.5">
      <div class="flex items-start justify-between gap-2">
        <h3 class="line-clamp-2 text-sm font-bold leading-snug text-gray-900">{{ props.experience.title }}</h3>
        <RatingStars v-if="props.experience.rating_count > 0" :rating="props.experience.rating_avg" size="sm" />
      </div>
      <p class="mt-1 flex items-center gap-1 text-xs text-gray-500">
        <AppIcon icon="heroicons-outline:map-pin" size="xs" />
        {{ props.experience.city }}
        <span v-if="props.experience.distance_meters != null" class="text-brand-primary">
          · {{ formatDistance(props.experience.distance_meters) }}
        </span>
      </p>
      <p class="mt-1 line-clamp-1 text-xs text-gray-400">{{ props.experience.summary }}</p>
      <p class="mt-1.5 text-sm font-extrabold text-gray-900 tabular-nums">
        {{ formatPrice(props.experience.price_amount) }}
        <span class="text-xs font-normal text-gray-400">/ {{ t('experience.perPerson') }}</span>
      </p>
    </div>
  </router-link>
</template>
