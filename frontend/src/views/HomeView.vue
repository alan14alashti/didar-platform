<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { categoryApi, experienceApi } from '../api/endpoints'
import type { Category, Experience } from '../api/types'
import ExperienceCard from '../components/ExperienceCard.vue'
import ExperienceCardSkeleton from '../components/ExperienceCardSkeleton.vue'
import EmptyState from '../components/EmptyState.vue'
import AppIcon from '../components/AppIcon.vue'
import { useBrandStore } from '../stores/brand'

const { t } = useI18n()
const brandStore = useBrandStore()

const categoryIcons: Record<string, string> = {
  cooking: 'heroicons-outline:fire',
  handicraft: 'heroicons-outline:paint-brush',
  nature: 'heroicons-outline:sun',
  history: 'heroicons-outline:building-library',
  music: 'heroicons-outline:musical-note',
  adventure: 'heroicons-outline:map',
}

const categories = ref<Category[]>([])
const activeCategory = ref<string>('')
const query = ref('')

const experiences = ref<Experience[]>([])
const totalCount = ref(0)
const page = ref(1)
const pageSize = 9
const isLoading = ref(true)
const isLoadingMore = ref(false)

let searchDebounce: ReturnType<typeof setTimeout> | undefined

async function loadCategories() {
  categories.value = await categoryApi.list()
}

async function search(reset = true) {
  if (reset) {
    isLoading.value = true
    page.value = 1
  } else {
    isLoadingMore.value = true
  }

  try {
    const result = await experienceApi.search({
      query: query.value || undefined,
      category: activeCategory.value || undefined,
      page: page.value,
      page_size: pageSize,
    })
    experiences.value = reset ? result.items : [...experiences.value, ...result.items]
    totalCount.value = result.total_count
  } finally {
    isLoading.value = false
    isLoadingMore.value = false
  }
}

function loadMore() {
  page.value += 1
  search(false)
}

function selectCategory(slug: string) {
  activeCategory.value = activeCategory.value === slug ? '' : slug
}

const skeletons = computed(() => Array.from({ length: 6 }))

watch([activeCategory], () => search(true))

watch(query, () => {
  clearTimeout(searchDebounce)
  searchDebounce = setTimeout(() => search(true), 350)
})

onMounted(() => {
  loadCategories()
  search(true)
})
</script>

<template>
  <div>
    <!-- Hero -->
    <section class="relative overflow-hidden brand-gradient">
      <div class="blob animate-blob -top-24 -start-24 h-80 w-80 bg-white/30" />
      <div class="blob animate-blob animate-blob-delay-2 -bottom-32 end-0 h-96 w-96 bg-white/20" />
      <div class="blob animate-blob animate-blob-delay-4 top-1/3 start-1/2 h-64 w-64 bg-black/10" />

      <div class="relative mx-auto max-w-7xl px-4 py-20 text-center sm:px-6 sm:py-28">
        <span class="animate-fade-up mb-4 inline-flex items-center gap-1.5 rounded-full bg-white/15 px-4 py-1.5 text-xs font-semibold text-white ring-1 ring-white/30">
          <AppIcon icon="heroicons-outline:sparkles" size="sm" />
          {{ brandStore.brand?.name ?? 'Didar' }}
        </span>
        <h1 class="animate-fade-up text-4xl font-extrabold leading-tight text-white sm:text-5xl lg:text-6xl">
          {{ brandStore.brand?.tagline || t('home.heroTitle') }}
        </h1>
        <p class="animate-fade-up mx-auto mt-4 max-w-xl text-base text-white/90 sm:text-lg">{{ t('home.heroSubtitle') }}</p>

        <div class="animate-fade-up mx-auto mt-9 max-w-xl">
          <div class="flex items-center gap-2 rounded-full bg-white p-1.5 shadow-card-hover">
            <AppIcon icon="heroicons-outline:magnifying-glass" size="lg" icon-class="ms-3 shrink-0 text-gray-400" />
            <input
              v-model="query"
              type="search"
              :placeholder="t('home.searchPlaceholder')"
              class="w-full border-0 bg-transparent py-2 text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none focus:ring-0"
            />
          </div>
        </div>
      </div>
    </section>

    <!-- Category chips -->
    <section class="mx-auto max-w-7xl px-4 py-8 sm:px-6">
      <div class="flex flex-wrap gap-2.5">
        <button
          type="button"
          class="flex items-center gap-1.5 rounded-full border px-4 py-2 text-sm font-semibold transition-all duration-200"
          :class="!activeCategory ? 'border-brand-primary bg-brand-primary text-white shadow-glow' : 'border-gray-200 bg-white text-gray-600 hover:border-brand-primary hover:text-brand-primary'"
          @click="activeCategory = ''"
        >
          <AppIcon icon="heroicons-outline:squares-2x2" size="sm" />
          {{ t('home.allCategories') }}
        </button>
        <button
          v-for="c in categories"
          :key="c.id"
          type="button"
          class="flex items-center gap-1.5 rounded-full border px-4 py-2 text-sm font-semibold transition-all duration-200"
          :class="activeCategory === c.slug ? 'border-brand-primary bg-brand-primary text-white shadow-glow' : 'border-gray-200 bg-white text-gray-600 hover:border-brand-primary hover:text-brand-primary'"
          @click="selectCategory(c.slug)"
        >
          <AppIcon :icon="categoryIcons[c.slug] ?? 'heroicons-outline:map-pin'" size="sm" />
          {{ $i18n.locale === 'fa' ? c.title_fa : c.title_en }}
        </button>
      </div>
    </section>

    <!-- Results -->
    <section class="mx-auto max-w-7xl px-4 pb-20 sm:px-6">
      <div v-if="isLoading" class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <ExperienceCardSkeleton v-for="(_, i) in skeletons" :key="i" />
      </div>

      <template v-else>
        <p class="mb-5 text-sm font-medium text-gray-500">{{ t('home.resultsCount', { n: totalCount }) }}</p>

        <EmptyState v-if="experiences.length === 0" :title="t('home.noResults')" />

        <div v-else class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <ExperienceCard
            v-for="(exp, idx) in experiences"
            :key="exp.id"
            :experience="exp"
            :style="{ animationDelay: `${(idx % pageSize) * 60}ms` }"
            class="animate-fade-up"
          />
        </div>

        <div v-if="experiences.length < totalCount" class="mt-10 flex justify-center">
          <button type="button" class="btn-outline disabled:opacity-50" :disabled="isLoadingMore" @click="loadMore">
            {{ isLoadingMore ? t('common.loading') : t('home.loadMore') }}
          </button>
        </div>
      </template>
    </section>
  </div>
</template>
