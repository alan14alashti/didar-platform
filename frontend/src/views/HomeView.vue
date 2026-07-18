<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { categoryApi, experienceApi } from '../api/endpoints'
import type { Category, Experience } from '../api/types'
import AppIcon from '../components/AppIcon.vue'
import EmptyState from '../components/EmptyState.vue'
import ExperienceExploreMap from '../components/ExperienceExploreMap.vue'
import ExperienceMapListItem from '../components/ExperienceMapListItem.vue'
import { useGeolocation } from '../composables/useGeolocation'
import { useBrandStore } from '../stores/brand'

const { t, locale } = useI18n()
const brandStore = useBrandStore()
const { location: userLocation, isLoading: geoLoading, error: geoError, requestLocation } = useGeolocation()

const categoryIcons: Record<string, string> = {
  cooking: 'heroicons-outline:fire',
  handicraft: 'heroicons-outline:paint-brush',
  nature: 'heroicons-outline:sun',
  history: 'heroicons-outline:building-library',
  music: 'heroicons-outline:musical-note',
  adventure: 'heroicons-outline:map',
}

const categories = ref<Category[]>([])
const activeCategory = ref('')
const query = ref('')
const experiences = ref<Experience[]>([])
const totalCount = ref(0)
const isLoading = ref(true)
const activeSlug = ref<string | null>(null)

let searchDebounce: ReturnType<typeof setTimeout> | undefined

async function loadCategories() {
  categories.value = await categoryApi.list()
}

async function search() {
  isLoading.value = true
  try {
    const result = await experienceApi.search({
      query: query.value || undefined,
      category: activeCategory.value || undefined,
      lat: userLocation.value?.lat,
      lng: userLocation.value?.lng,
      radius_km: userLocation.value ? 1200 : undefined,
      page: 1,
      page_size: 50,
    })
    experiences.value = result.items
    totalCount.value = result.total_count
    if (activeSlug.value && !experiences.value.some((e) => e.slug === activeSlug.value)) {
      activeSlug.value = experiences.value[0]?.slug ?? null
    }
  } finally {
    isLoading.value = false
  }
}

function selectCategory(slug: string) {
  activeCategory.value = activeCategory.value === slug ? '' : slug
}

function selectExperience(slug: string) {
  activeSlug.value = slug
}

const locationHint = computed(() => {
  if (geoLoading.value) return t('home.locatingYou')
  if (geoError.value) return t('home.locationUnavailable')
  if (userLocation.value) return t('home.nearYou')
  return ''
})

watch([activeCategory], () => search())

watch(query, () => {
  clearTimeout(searchDebounce)
  searchDebounce = setTimeout(() => search(), 350)
})

watch(userLocation, () => search())

onMounted(async () => {
  await loadCategories()
  await search()
})
</script>

<template>
  <div class="flex min-h-[calc(100vh-4.5rem)] flex-col">
    <!-- Compact explore header -->
    <section class="border-b border-gray-200/80 bg-white/90 px-4 py-4 backdrop-blur sm:px-6">
      <div class="mx-auto flex max-w-[1600px] flex-col gap-4">
        <div class="flex flex-wrap items-end justify-between gap-3">
          <div>
            <p class="text-xs font-semibold uppercase tracking-wide text-brand-primary">
              {{ brandStore.brand?.name ?? 'Didar' }}
            </p>
            <h1 class="text-xl font-extrabold text-gray-900 sm:text-2xl">
              {{ t('home.mapExploreTitle') }}
            </h1>
            <p class="mt-1 text-sm text-gray-500">{{ t('home.mapExploreSubtitle') }}</p>
          </div>
          <div v-if="locationHint" class="flex items-center gap-2 text-xs font-medium text-gray-500">
            <span
              class="inline-flex items-center gap-1 rounded-full px-2.5 py-1"
              :class="userLocation ? 'bg-blue-50 text-blue-700' : 'bg-gray-100 text-gray-600'"
            >
              <AppIcon icon="heroicons-outline:map-pin" size="xs" />
              {{ locationHint }}
            </span>
            <button
              v-if="geoError"
              type="button"
              class="rounded-full border border-gray-200 px-2.5 py-1 text-xs font-semibold text-gray-600 hover:border-brand-primary hover:text-brand-primary"
              @click="requestLocation"
            >
              {{ t('home.enableLocation') }}
            </button>
          </div>
        </div>

        <div class="flex flex-col gap-3 lg:flex-row lg:items-center">
          <div class="flex flex-1 items-center gap-2 rounded-full border border-gray-200 bg-white px-3 py-2 shadow-soft">
            <AppIcon icon="heroicons-outline:magnifying-glass" size="md" icon-class="shrink-0 text-gray-400" />
            <input
              v-model="query"
              type="search"
              :placeholder="t('home.searchPlaceholder')"
              class="w-full border-0 bg-transparent py-1 text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none focus:ring-0"
            />
          </div>
          <div class="flex gap-2 overflow-x-auto pb-1 scrollbar-thin">
            <button
              type="button"
              class="flex shrink-0 items-center gap-1.5 rounded-full border px-3.5 py-2 text-xs font-semibold transition-all"
              :class="!activeCategory ? 'border-brand-primary bg-brand-primary text-white' : 'border-gray-200 bg-white text-gray-600 hover:border-brand-primary hover:text-brand-primary'"
              @click="activeCategory = ''"
            >
              <AppIcon icon="heroicons-outline:squares-2x2" size="sm" />
              {{ t('home.allCategories') }}
            </button>
            <button
              v-for="c in categories"
              :key="c.id"
              type="button"
              class="flex shrink-0 items-center gap-1.5 rounded-full border px-3.5 py-2 text-xs font-semibold transition-all"
              :class="activeCategory === c.slug ? 'border-brand-primary bg-brand-primary text-white' : 'border-gray-200 bg-white text-gray-600 hover:border-brand-primary hover:text-brand-primary'"
              @click="selectCategory(c.slug)"
            >
              <AppIcon :icon="categoryIcons[c.slug] ?? 'heroicons-outline:map-pin'" size="sm" />
              {{ locale === 'fa' ? c.title_fa : c.title_en }}
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- Map-first layout -->
    <section class="mx-auto flex w-full max-w-[1600px] flex-1 flex-col lg:flex-row">
      <div class="relative min-h-[42vh] flex-1 lg:min-h-[calc(100vh-12rem)] lg:order-1">
        <ExperienceExploreMap
          :experiences="experiences"
          :user-location="userLocation"
          :active-slug="activeSlug"
          @select="selectExperience"
        />
        <div
          v-if="isLoading"
          class="absolute inset-0 grid place-items-center bg-white/50 backdrop-blur-sm"
        >
          <p class="rounded-full bg-white px-4 py-2 text-sm font-medium text-gray-600 shadow-soft">
            {{ t('common.loading') }}
          </p>
        </div>
      </div>

      <aside class="flex w-full flex-col border-t border-gray-200 bg-brand-surface lg:w-[420px] lg:border-s lg:border-t-0 xl:w-[460px] lg:order-2">
        <div class="border-b border-gray-200 px-4 py-3">
          <p class="text-sm font-semibold text-gray-800">{{ t('home.resultsCount', { n: totalCount }) }}</p>
        </div>

        <div class="flex-1 overflow-y-auto px-3 py-3 scrollbar-thin">
          <EmptyState v-if="!isLoading && experiences.length === 0" :title="t('home.noResults')" />

          <div v-else class="space-y-2.5">
            <ExperienceMapListItem
              v-for="exp in experiences"
              :key="exp.id"
              :experience="exp"
              :active="activeSlug === exp.slug"
              @mouseenter="activeSlug = exp.slug"
              @focusin="activeSlug = exp.slug"
            />
          </div>
        </div>
      </aside>
    </section>
  </div>
</template>
