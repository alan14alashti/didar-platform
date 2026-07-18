<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { categoryApi, experienceApi } from '../api/endpoints'
import type { Category, Experience } from '../api/types'
import EmptyState from '../components/EmptyState.vue'
import ExperienceCategorySection from '../components/ExperienceCategorySection.vue'
import ExperienceExploreMap from '../components/ExperienceExploreMap.vue'
import ExperienceMapListItem from '../components/ExperienceMapListItem.vue'
import AppIcon from '../components/AppIcon.vue'
import ExploreSearchHeader from '../components/layout/ExploreSearchHeader.vue'
import { useGeolocation } from '../composables/useGeolocation'

const { t } = useI18n()
const router = useRouter()
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
const viewMode = ref<'map' | 'list'>('map')
const mobileMapPanel = ref<'map' | 'list'>('map')
const experiences = ref<Experience[]>([])
const isLoading = ref(true)
const selectedSlug = ref<string | null>(null)

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
  } finally {
    isLoading.value = false
  }
}

function selectExperience(slug: string) {
  selectedSlug.value = slug || null
}

function openDetail(slug: string) {
  router.push({ name: 'experience-detail', params: { slug } })
}

const locationHint = computed(() => {
  if (geoLoading.value) return t('home.locatingYou')
  if (geoError.value) return t('home.locationUnavailable')
  if (userLocation.value) return t('home.nearYou')
  return ''
})

const groupedSections = computed(() => {
  if (activeCategory.value) {
    const cat = categories.value.find((c) => c.slug === activeCategory.value)
    if (!cat) return []
    return [{ category: cat, items: experiences.value }]
  }

  const bySlug = new Map<string, Experience[]>()
  for (const exp of experiences.value) {
    const slug = exp.category_slug ?? 'other'
    if (!bySlug.has(slug)) bySlug.set(slug, [])
    bySlug.get(slug)!.push(exp)
  }

  return categories.value
    .filter((c) => (bySlug.get(c.slug)?.length ?? 0) > 0)
    .map((c) => ({ category: c, items: bySlug.get(c.slug)! }))
})

watch([activeCategory], () => search())

watch(query, () => {
  clearTimeout(searchDebounce)
  searchDebounce = setTimeout(() => search(), 350)
})

watch(userLocation, () => search())

watch(viewMode, (mode) => {
  if (mode === 'map') mobileMapPanel.value = 'map'
})

onMounted(async () => {
  await loadCategories()
  await search()
})
</script>

<template>
  <div
    class="flex flex-col bg-brand-surface"
    :class="viewMode === 'map' ? 'h-[100dvh] overflow-hidden' : 'min-h-screen'"
  >
    <ExploreSearchHeader
      v-model:query="query"
      v-model:category="activeCategory"
      v-model:view-mode="viewMode"
      :categories="categories"
      :location-hint="locationHint"
      :show-location-action="Boolean(geoError)"
      @search="search"
      @request-location="requestLocation"
    />

    <!-- Map view -->
    <section
      v-if="viewMode === 'map'"
      class="relative h-[calc(100dvh-4.25rem)] overflow-hidden xl:flex xl:min-h-[calc(100vh-5.5rem)] xl:flex-1 xl:flex-row"
    >
      <div class="relative h-full w-full xl:min-h-0 xl:flex-1">
        <ExperienceExploreMap
          :experiences="experiences"
          :user-location="userLocation"
          :selected-slug="selectedSlug"
          @select="selectExperience"
          @open-detail="openDetail"
        />
        <div v-if="isLoading" class="absolute inset-0 z-10 grid place-items-center bg-white/40 backdrop-blur-[1px]">
          <p class="rounded-full bg-white px-4 py-2 text-sm font-medium text-gray-600 shadow-soft">{{ t('common.loading') }}</p>
        </div>
      </div>

      <aside class="hidden w-[400px] shrink-0 flex-col border-s border-gray-200 bg-white xl:flex">
        <div class="flex-1 overflow-y-auto px-2 py-2 scrollbar-thin">
          <EmptyState v-if="!isLoading && experiences.length === 0" :title="t('home.noResults')" />
          <div v-else class="space-y-1">
            <ExperienceMapListItem
              v-for="exp in experiences"
              :key="exp.id"
              :experience="exp"
              :active="selectedSlug === exp.slug"
            />
          </div>
        </div>
      </aside>

      <div
        v-if="mobileMapPanel === 'list'"
        class="absolute inset-0 z-30 flex flex-col bg-white xl:hidden"
      >
        <div class="flex-1 overflow-y-auto px-2 py-2 pt-3 scrollbar-thin">
          <EmptyState v-if="!isLoading && experiences.length === 0" :title="t('home.noResults')" />
          <div v-else class="space-y-1">
            <ExperienceMapListItem
              v-for="exp in experiences"
              :key="`mobile-${exp.id}`"
              :experience="exp"
              :active="selectedSlug === exp.slug"
            />
          </div>
        </div>
      </div>

      <button
        type="button"
        class="fixed bottom-5 start-1/2 z-40 flex -translate-x-1/2 items-center gap-2 rounded-full bg-[#222222] px-5 py-3 text-sm font-semibold text-white shadow-[0_6px_16px_rgba(0,0,0,0.2)] xl:hidden"
        @click="mobileMapPanel = mobileMapPanel === 'map' ? 'list' : 'map'"
      >
        <AppIcon
          :icon="mobileMapPanel === 'map' ? 'heroicons-outline:squares-2x2' : 'heroicons-outline:map'"
          size="sm"
        />
        {{ mobileMapPanel === 'map' ? t('home.showList') : t('home.showMap') }}
      </button>
    </section>

    <!-- List view (Airbnb-style categorized scroll) -->
    <section v-else class="mx-auto w-full max-w-7xl flex-1 px-4 py-6 sm:px-6">
      <h1 class="mb-5 text-xl font-extrabold text-gray-900">{{ t('home.listExploreTitle') }}</h1>

      <div v-if="isLoading" class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <div v-for="i in 6" :key="i" class="h-64 rounded-3xl shimmer-bg" />
      </div>

      <EmptyState v-else-if="experiences.length === 0" :title="t('home.noResults')" />

      <div v-else class="space-y-10">
        <ExperienceCategorySection
          v-for="section in groupedSections"
          :key="section.category.id"
          :category="section.category"
          :experiences="section.items"
          :icon="categoryIcons[section.category.slug]"
        />
      </div>
    </section>
  </div>
</template>
