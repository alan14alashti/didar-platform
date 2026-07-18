<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import type { Category } from '../../api/types'
import { setLocale } from '../../i18n'
import { useAuthStore } from '../../stores/auth'
import { useBrandStore } from '../../stores/brand'
import AppIcon from '../AppIcon.vue'

const props = defineProps<{
  categories: Category[]
  categoryIcons: Record<string, string>
  locationHint?: string
  showLocationAction?: boolean
}>()

void props

const query = defineModel<string>('query', { default: '' })
const activeCategory = defineModel<string>('category', { default: '' })
const viewMode = defineModel<'map' | 'list'>('viewMode', { default: 'map' })

const emit = defineEmits<{ search: []; requestLocation: [] }>()

const { locale, t } = useI18n()
const brandStore = useBrandStore()
const authStore = useAuthStore()
const router = useRouter()

const menuOpen = ref(false)
const filtersOpen = ref(false)

function toggleLocale() {
  setLocale(locale.value === 'fa' ? 'en' : 'fa')
}

function submitSearch() {
  emit('search')
}

async function handleLogout() {
  await authStore.logout()
  menuOpen.value = false
  router.push({ name: 'home' })
}
</script>

<template>
  <header class="sticky top-0 z-50 border-b border-gray-200 bg-white/95 backdrop-blur">
    <div class="mx-auto flex max-w-[1600px] items-center gap-3 px-4 py-3 sm:px-6">
      <!-- Logo -->
      <router-link to="/" class="group flex shrink-0 items-center gap-2">
        <span class="grid h-9 w-9 place-items-center rounded-xl brand-gradient text-base font-extrabold text-white shadow-glow">
          {{ brandStore.brand?.name?.[0] ?? 'D' }}
        </span>
        <span class="hidden font-extrabold text-gray-900 sm:inline">{{ brandStore.brand?.name ?? 'Didar' }}</span>
      </router-link>

      <!-- Airbnb-style search bar -->
      <div class="hidden min-w-0 flex-1 md:block">
        <form
          class="mx-auto flex max-w-3xl items-stretch overflow-hidden rounded-full border border-gray-200 bg-white shadow-soft transition hover:shadow-card"
          @submit.prevent="submitSearch"
        >
          <label class="flex min-w-0 flex-[1.2] flex-col border-e border-gray-200 px-4 py-2">
            <span class="text-[10px] font-bold uppercase tracking-wide text-gray-500">{{ t('home.searchWhere') }}</span>
            <input
              v-model="query"
              type="search"
              :placeholder="t('home.searchPlaceholder')"
              class="w-full border-0 bg-transparent p-0 text-sm text-gray-900 placeholder:text-gray-400 focus:outline-none focus:ring-0"
            />
          </label>

          <label class="relative flex min-w-[9rem] flex-1 flex-col px-4 py-2">
            <span class="text-[10px] font-bold uppercase tracking-wide text-gray-500">{{ t('home.searchCategory') }}</span>
            <select
              v-model="activeCategory"
              class="w-full cursor-pointer border-0 bg-transparent p-0 text-sm text-gray-900 focus:outline-none focus:ring-0"
            >
              <option value="">{{ t('home.allCategories') }}</option>
              <option v-for="c in categories" :key="c.id" :value="c.slug">
                {{ locale === 'fa' ? c.title_fa : c.title_en }}
              </option>
            </select>
          </label>

          <button type="submit" class="m-1.5 flex items-center gap-2 rounded-full bg-brand-primary px-5 text-sm font-bold text-white hover:opacity-95">
            <AppIcon icon="heroicons-outline:magnifying-glass" size="sm" />
            <span class="hidden lg:inline">{{ t('common.search') }}</span>
          </button>
        </form>
      </div>

      <!-- Mobile search toggle -->
      <button
        type="button"
        class="flex min-w-0 flex-1 items-center gap-2 rounded-full border border-gray-200 px-3 py-2 text-sm text-gray-500 md:hidden"
        @click="filtersOpen = !filtersOpen"
      >
        <AppIcon icon="heroicons-outline:magnifying-glass" size="sm" />
        <span class="truncate">{{ query || t('home.searchPlaceholder') }}</span>
      </button>

      <!-- Right controls -->
      <div class="flex shrink-0 items-center gap-2">
        <div class="hidden items-center rounded-full border border-gray-200 p-0.5 sm:flex">
          <button
            type="button"
            class="flex items-center gap-1 rounded-full px-3 py-1.5 text-xs font-semibold transition"
            :class="viewMode === 'map' ? 'bg-gray-900 text-white' : 'text-gray-600 hover:bg-gray-50'"
            @click="viewMode = 'map'"
          >
            <AppIcon icon="heroicons-outline:map" size="sm" />
            {{ t('home.mapView') }}
          </button>
          <button
            type="button"
            class="flex items-center gap-1 rounded-full px-3 py-1.5 text-xs font-semibold transition"
            :class="viewMode === 'list' ? 'bg-gray-900 text-white' : 'text-gray-600 hover:bg-gray-50'"
            @click="viewMode = 'list'"
          >
            <AppIcon icon="heroicons-outline:squares-2x2" size="sm" />
            {{ t('home.listView') }}
          </button>
        </div>

        <button
          type="button"
          class="rounded-full border border-gray-200 px-3 py-1.5 text-xs font-bold text-gray-600 hover:border-brand-primary hover:text-brand-primary"
          @click="toggleLocale"
        >
          {{ locale === 'fa' ? 'EN' : 'فا' }}
        </button>

        <router-link v-if="!authStore.isAuthenticated" to="/login" class="btn-primary hidden !px-3 !py-1.5 text-xs sm:inline-flex">
          {{ t('nav.login') }}
        </router-link>

        <div v-else class="relative">
          <button
            type="button"
            class="grid h-9 w-9 place-items-center overflow-hidden rounded-full bg-gray-100 text-sm font-bold text-gray-700"
            @click="menuOpen = !menuOpen"
          >
            <img v-if="authStore.user?.avatar_url" :src="authStore.user.avatar_url" class="h-full w-full object-cover" alt="" />
            <span v-else>{{ (authStore.user?.full_name || authStore.user?.phone || '?')[0] }}</span>
          </button>
          <div
            v-if="menuOpen"
            class="absolute end-0 mt-2 w-48 overflow-hidden rounded-2xl border border-gray-100 bg-white py-2 shadow-card-hover"
            @click="menuOpen = false"
          >
            <router-link to="/trips" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50">{{ t('nav.myTrips') }}</router-link>
            <router-link to="/profile" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50">{{ t('nav.profile') }}</router-link>
            <router-link v-if="authStore.isHost" to="/host" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50">{{ t('nav.hostDashboard') }}</router-link>
            <div class="my-1 border-t border-gray-100" />
            <button type="button" class="block w-full px-4 py-2 text-start text-sm text-red-600 hover:bg-red-50" @click="handleLogout">
              {{ t('nav.logout') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile filters panel -->
    <div v-if="filtersOpen" class="border-t border-gray-100 px-4 py-3 md:hidden">
      <input v-model="query" type="search" :placeholder="t('home.searchPlaceholder')" class="input-fancy mb-3" />
      <select v-model="activeCategory" class="input-fancy mb-3">
        <option value="">{{ t('home.allCategories') }}</option>
        <option v-for="c in categories" :key="c.id" :value="c.slug">
          {{ locale === 'fa' ? c.title_fa : c.title_en }}
        </option>
      </select>
      <div class="mb-3 flex rounded-full border border-gray-200 p-0.5">
        <button type="button" class="flex-1 rounded-full py-2 text-xs font-semibold" :class="viewMode === 'map' ? 'bg-gray-900 text-white' : ''" @click="viewMode = 'map'">{{ t('home.mapView') }}</button>
        <button type="button" class="flex-1 rounded-full py-2 text-xs font-semibold" :class="viewMode === 'list' ? 'bg-gray-900 text-white' : ''" @click="viewMode = 'list'">{{ t('home.listView') }}</button>
      </div>
      <button type="button" class="btn-primary w-full" @click="submitSearch">{{ t('common.search') }}</button>
    </div>

    <!-- Category filter chips -->
    <div class="border-t border-gray-100">
      <div class="mx-auto flex max-w-[1600px] gap-2 overflow-x-auto px-4 py-2.5 scrollbar-thin sm:px-6">
        <button
          type="button"
          class="flex shrink-0 items-center gap-1.5 rounded-full border px-3.5 py-1.5 text-xs font-semibold transition"
          :class="!activeCategory ? 'border-gray-900 bg-gray-900 text-white' : 'border-gray-200 bg-white text-gray-600'"
          @click="activeCategory = ''"
        >
          {{ t('home.allCategories') }}
        </button>
        <button
          v-for="c in categories"
          :key="c.id"
          type="button"
          class="flex shrink-0 items-center gap-1.5 rounded-full border px-3.5 py-1.5 text-xs font-semibold transition"
          :class="activeCategory === c.slug ? 'border-gray-900 bg-gray-900 text-white' : 'border-gray-200 bg-white text-gray-600'"
          @click="activeCategory = activeCategory === c.slug ? '' : c.slug"
        >
          <AppIcon :icon="categoryIcons[c.slug] ?? 'heroicons-outline:map-pin'" size="sm" />
          {{ locale === 'fa' ? c.title_fa : c.title_en }}
        </button>

        <span
          v-if="locationHint"
          class="ms-auto hidden shrink-0 items-center gap-1 rounded-full bg-blue-50 px-3 py-1.5 text-xs font-medium text-blue-700 lg:inline-flex"
        >
          <AppIcon icon="heroicons-outline:map-pin" size="xs" />
          {{ locationHint }}
        </span>
        <button
          v-if="showLocationAction"
          type="button"
          class="hidden shrink-0 rounded-full border border-gray-200 px-3 py-1.5 text-xs font-semibold text-gray-600 lg:inline-flex"
          @click="emit('requestLocation')"
        >
          {{ t('home.enableLocation') }}
        </button>
      </div>
    </div>
  </header>
</template>
