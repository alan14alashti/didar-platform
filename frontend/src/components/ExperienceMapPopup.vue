<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Experience } from '../api/types'
import { useFormatters } from '../composables/useFormatters'
import AppIcon from './AppIcon.vue'

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
    class="exp-airbnb-popup absolute z-20 w-[min(16.5rem,calc(100vw-2rem))] -translate-x-1/2 -translate-y-full"
    :style="style"
  >
    <div class="relative overflow-hidden rounded-xl bg-white shadow-[0_6px_16px_rgba(0,0,0,0.12)]">
      <button
        type="button"
        class="absolute end-3 top-3 z-10 grid h-7 w-7 place-items-center rounded-full bg-white text-[#222222] shadow-[0_2px_4px_rgba(0,0,0,0.18)] transition hover:scale-105"
        @click.stop="emit('close')"
      >
        <AppIcon icon="heroicons-outline:x-mark" size="sm" />
      </button>

      <button type="button" class="block w-full text-start" @click="emit('open')">
        <div class="relative aspect-square w-full overflow-hidden bg-[#f7f7f7]">
          <img
            v-if="experience.cover_photo_url"
            :src="experience.cover_photo_url"
            :alt="experience.title"
            class="h-full w-full object-cover transition-transform duration-300 hover:scale-105"
          />

          <span
            v-if="experience.rating_count > 0"
            class="absolute start-3 top-3 flex items-center gap-1 rounded-md bg-white/95 px-2 py-1 text-xs font-semibold text-[#222222] shadow-[0_2px_4px_rgba(0,0,0,0.12)]"
          >
            <AppIcon icon="heroicons-outline:star" size="xs" icon-class="text-[#222222]" />
            {{ experience.rating_avg.toFixed(2) }}
          </span>
        </div>

        <div class="space-y-0.5 p-3">
          <div class="flex items-start justify-between gap-2 pe-5">
            <h3 class="line-clamp-2 text-[15px] font-semibold leading-snug text-[#222222]">
              {{ experience.title }}
            </h3>
          </div>

          <p class="truncate text-sm text-[#717171]">
            {{ experience.city }}
            <span v-if="experience.distance_meters != null">
              · {{ formatDistance(experience.distance_meters) }}
            </span>
          </p>

          <p class="truncate text-sm text-[#717171]">
            {{ t('experience.durationMinutes', { n: experience.duration_minutes }) }}
          </p>

          <p class="flex flex-wrap items-baseline gap-x-1 pt-1 text-sm text-[#222222]">
            <span class="font-semibold tabular-nums">{{ formatPrice(experience.price_amount) }}</span>
            <span class="font-normal text-[#717171]">{{ t('experience.perPerson') }}</span>
          </p>
        </div>
      </button>
    </div>

    <div
      class="absolute -bottom-[7px] start-1/2 h-3.5 w-3.5 -translate-x-1/2 rotate-45 bg-white shadow-[2px_2px_4px_rgba(0,0,0,0.08)]"
      aria-hidden="true"
    />
  </div>
</template>
