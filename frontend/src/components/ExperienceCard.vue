<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Experience } from '../api/types'
import { useFormatters } from '../composables/useFormatters'
import AppIcon from './AppIcon.vue'

const props = defineProps<{ experience: Experience }>()
const { t } = useI18n()
const { formatPrice } = useFormatters()
</script>

<template>
  <router-link
    :to="{ name: 'experience-detail', params: { slug: props.experience.slug } }"
    class="group block"
  >
    <div class="relative mb-3 aspect-square w-full overflow-hidden rounded-xl bg-[#f7f7f7]">
      <img
        v-if="props.experience.cover_photo_url"
        :src="props.experience.cover_photo_url"
        :alt="props.experience.title"
        class="h-full w-full object-cover transition-transform duration-300 ease-out group-hover:scale-105"
        loading="lazy"
      />

      <button
        type="button"
        class="absolute end-3 top-3 grid h-8 w-8 place-items-center rounded-full text-[#222222] transition hover:scale-105"
        @click.prevent
      >
        <AppIcon icon="heroicons-outline:heart" size="md" />
      </button>
    </div>

    <div class="space-y-0.5 px-0.5">
      <div class="flex items-start justify-between gap-2">
        <h3 class="line-clamp-1 text-[15px] font-semibold text-[#222222]">
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

      <p class="truncate text-sm text-[#717171]">{{ props.experience.city }}</p>

      <p class="truncate text-sm text-[#717171]">
        {{ t('experience.durationMinutes', { n: props.experience.duration_minutes }) }}
      </p>

      <p class="flex flex-wrap items-baseline gap-x-1 pt-0.5 text-sm text-[#222222]">
        <span class="font-semibold tabular-nums">{{ formatPrice(props.experience.price_amount) }}</span>
        <span class="font-normal text-[#717171]">{{ t('experience.perPerson') }}</span>
      </p>
    </div>
  </router-link>
</template>
