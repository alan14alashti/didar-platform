<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Experience } from '../api/types'
import { useFormatters } from '../composables/useFormatters'
import AppIcon from './AppIcon.vue'
import RatingStars from './RatingStars.vue'

const props = defineProps<{ experience: Experience }>()
const { t } = useI18n()
const { formatPrice } = useFormatters()
</script>

<template>
  <router-link
    :to="{ name: 'experience-detail', params: { slug: props.experience.slug } }"
    class="group block overflow-hidden rounded-3xl border border-gray-100 bg-white shadow-card transition-all duration-300 hover:-translate-y-1.5 hover:shadow-card-hover"
  >
    <div class="relative aspect-[4/3] w-full overflow-hidden bg-gray-100">
      <img
        v-if="props.experience.cover_photo_url"
        :src="props.experience.cover_photo_url"
        :alt="props.experience.title"
        class="h-full w-full object-cover transition-transform duration-500 ease-out group-hover:scale-110"
        loading="lazy"
      />
      <div class="absolute inset-0 bg-gradient-to-t from-black/35 via-transparent to-transparent opacity-0 transition-opacity duration-300 group-hover:opacity-100" />

      <span
        v-if="props.experience.rating_count > 0"
        class="absolute start-3 top-3 flex items-center gap-1 rounded-full bg-white/90 px-2.5 py-1 text-xs font-bold text-gray-800 shadow-soft backdrop-blur"
      >
        <AppIcon icon="heroicons-outline:star" size="xs" icon-class="text-amber-500" />
        {{ props.experience.rating_avg.toFixed(1) }}
      </span>

      <span class="absolute end-3 top-3 grid h-8 w-8 place-items-center rounded-full bg-white/85 text-gray-500 shadow-soft backdrop-blur transition hover:text-brand-primary">
        <AppIcon icon="heroicons-outline:heart" size="md" />
      </span>
    </div>

    <div class="space-y-1.5 p-4">
      <div class="flex items-center justify-between gap-2 text-xs font-medium text-gray-400">
        <span class="flex items-center gap-1">
          <AppIcon icon="heroicons-outline:map-pin" size="sm" />
          {{ props.experience.city }}
        </span>
        <RatingStars v-if="props.experience.rating_count === 0" :rating="0" />
      </div>
      <h3 class="line-clamp-1 font-bold text-gray-900 transition-colors group-hover:text-brand-primary">{{ props.experience.title }}</h3>
      <p class="line-clamp-2 text-sm text-gray-500">{{ props.experience.summary }}</p>
      <p class="flex items-baseline gap-1 pt-1.5 text-sm font-extrabold text-gray-900 tabular-nums">
        {{ formatPrice(props.experience.price_amount) }}
        <span class="text-xs font-normal text-gray-400">/ {{ t('experience.perPerson') }}</span>
      </p>
    </div>
  </router-link>
</template>
