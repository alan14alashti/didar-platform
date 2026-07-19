<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { useBrandStore } from '../../stores/brand'
import { useAuthStore } from '../../stores/auth'
import { setLocale } from '../../i18n'

const brandStore = useBrandStore()
const authStore = useAuthStore()
const { locale, t } = useI18n()
const router = useRouter()

const menuOpen = ref(false)
const scrolled = ref(false)

function onScroll() {
  scrolled.value = window.scrollY > 8
}

onMounted(() => window.addEventListener('scroll', onScroll, { passive: true }))
onUnmounted(() => window.removeEventListener('scroll', onScroll))

function toggleLocale() {
  setLocale(locale.value === 'fa' ? 'en' : 'fa')
}

async function handleLogout() {
  await authStore.logout()
  menuOpen.value = false
  router.push({ name: 'home' })
}
</script>

<template>
  <header
    class="glass-panel sticky top-0 z-40 border-b transition-shadow"
    :class="scrolled ? 'border-gray-200/70 shadow-soft' : 'border-transparent'"
  >
    <div class="mx-auto flex max-w-7xl items-center justify-between gap-4 px-4 py-3.5 sm:px-6">
      <router-link to="/" class="group flex shrink-0 items-center gap-2.5">
        <span
          class="grid h-10 w-10 place-items-center rounded-2xl brand-gradient text-lg font-extrabold text-white shadow-glow transition-transform duration-300 group-hover:scale-105 group-hover:rotate-3"
        >
          {{ brandStore.brand?.name?.[0] ?? 'D' }}
        </span>
        <span class="text-lg font-extrabold tracking-tight text-gray-900">{{ brandStore.brand?.name ?? 'Didar' }}</span>
      </router-link>

      <nav class="hidden items-center gap-1 rounded-full border border-gray-100 bg-white/70 p-1 text-sm font-medium text-gray-600 shadow-soft md:flex">
        <router-link
          to="/"
          class="rounded-full px-4 py-1.5 transition-colors hover:bg-gray-100"
          active-class="!bg-brand-primary/10 !text-brand-primary"
        >
          {{ t('nav.home') }}
        </router-link>
        <router-link
          v-if="authStore.isAuthenticated"
          to="/trips"
          class="rounded-full px-4 py-1.5 transition-colors hover:bg-gray-100"
          active-class="!bg-brand-primary/10 !text-brand-primary"
        >
          {{ t('nav.myTrips') }}
        </router-link>
        <router-link
          v-if="authStore.isHost"
          to="/host"
          class="rounded-full px-4 py-1.5 transition-colors hover:bg-gray-100"
          active-class="!bg-brand-primary/10 !text-brand-primary"
        >
          {{ t('nav.hostDashboard') }}
        </router-link>
        <router-link
          v-else-if="authStore.isAuthenticated"
          to="/profile"
          class="rounded-full px-4 py-1.5 transition-colors hover:bg-gray-100"
        >
          {{ t('nav.becomeHost') }}
        </router-link>
      </nav>

      <div class="flex items-center gap-2.5">
        <button
          type="button"
          class="rounded-full border border-gray-200 bg-white px-3 py-1.5 text-xs font-bold text-gray-600 shadow-soft transition hover:border-brand-primary hover:text-brand-primary"
          @click="toggleLocale"
        >
          {{ locale === 'fa' ? 'EN' : 'فا' }}
        </button>

        <router-link v-if="!authStore.isAuthenticated" to="/login" class="btn-primary !px-4 !py-2 text-sm">
          {{ t('nav.login') }}
        </router-link>

        <div v-else class="relative">
          <button
            type="button"
            class="grid h-10 w-10 place-items-center overflow-hidden rounded-full bg-gray-100 text-sm font-bold text-gray-700 ring-2 ring-transparent transition hover:ring-brand-primary/30"
            @click="menuOpen = !menuOpen"
          >
            <img v-if="authStore.user?.avatar_url" :src="authStore.user.avatar_url" class="h-full w-full object-cover" alt="" />
            <span v-else>{{ (authStore.user?.full_name || authStore.user?.email || '?')[0] }}</span>
          </button>

          <transition
            enter-active-class="transition duration-150 ease-out"
            enter-from-class="opacity-0 scale-95 translate-y-1"
            enter-to-class="opacity-100 scale-100 translate-y-0"
            leave-active-class="transition duration-100 ease-in"
            leave-from-class="opacity-100"
            leave-to-class="opacity-0"
          >
            <div
              v-if="menuOpen"
              class="absolute end-0 mt-2 w-52 overflow-hidden rounded-2xl border border-gray-100 bg-white py-2 shadow-card-hover"
              @click="menuOpen = false"
            >
              <router-link to="/profile" class="block px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50">
                {{ t('nav.profile') }}
              </router-link>
              <router-link to="/trips" class="block px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 md:hidden">
                {{ t('nav.myTrips') }}
              </router-link>
              <router-link v-if="authStore.isHost" to="/host" class="block px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 md:hidden">
                {{ t('nav.hostDashboard') }}
              </router-link>
              <div class="my-1 border-t border-gray-100" />
              <button type="button" class="block w-full px-4 py-2.5 text-start text-sm font-medium text-red-600 hover:bg-red-50" @click="handleLogout">
                {{ t('nav.logout') }}
              </button>
            </div>
          </transition>
        </div>
      </div>
    </div>
  </header>
</template>
