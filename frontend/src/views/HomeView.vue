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
const experiences = ref<Experience[]>([])
const totalCount = ref(0)
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
    totalCount.value = result.total_count
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

onMounted(async () => {
  await loadCategories()
  await search()
})
</script>

<template>
  <div class="flex min-h-screen flex-col bg-brand-surface">
    <ExploreSearchHeader
      v-model:query="query"
      v-model:category="activeCategory"
      v-model:view-mode="viewMode"
      :categories="categories"
      :category-icons="categoryIcons"
      :location-hint="locationHint"
      :show-location-action="Boolean(geoError)"
      @search="search"
      @request-location="requestLocation"
    />

    <!-- Map view -->
    <section v-if="viewMode === 'map'" class="relative flex min-h-0 flex-1 flex-col lg:flex-row">
      <div class="relative min-h-[50vh] flex-1 lg:min-h-[calc(100vh-8.5rem)]">
        <ExperienceExploreMap
          :experiences="experiences"
          :user-location="userLocation"
          :selected-slug="selectedSlug"
          @select="selectExperience"
          @open-detail="openDetail"
        />
        <div v-if="isLoading" class="absolute inset-0 grid place-items-center bg-white/40 backdrop-blur-[1px]">
          <p class="rounded-full bg-white px-4 py-2 text-sm font-medium text-gray-600 shadow-soft">{{ t('common.loading') }}</p>
        </div>
      </div>

      <aside class="hidden w-[400px] shrink-0 flex-col border-s border-gray-200 bg-white xl:flex">
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
              :active="selectedSlug === exp.slug"
            />
          </div>
        </div>
      </aside>
    </section>

    <!-- List view (Airbnb-style categorized scroll) -->
    <section v-else class="mx-auto w-full max-w-7xl flex-1 px-4 py-6 sm:px-6">
      <div class="mb-5 flex items-center justify-between gap-3">
        <h1 class="text-xl font-extrabold text-gray-900">{{ t('home.listExploreTitle') }}</h1>
        <p class="text-sm text-gray-500">{{ t('home.resultsCount', { n: totalCount }) }}</p>
      </div>

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
