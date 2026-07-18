<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import AppFooter from './components/layout/AppFooter.vue'
import AppHeader from './components/layout/AppHeader.vue'

const route = useRoute()
const isExploreHome = computed(() => route.name === 'home')
</script>

<template>
  <div class="flex min-h-full flex-col bg-brand-surface">
    <AppHeader v-if="!isExploreHome" />
    <main class="flex-1">
      <router-view v-slot="{ Component, route: childRoute }">
        <transition name="page-fade" mode="out-in">
          <component :is="Component" :key="childRoute.fullPath" />
        </transition>
      </router-view>
    </main>
    <AppFooter v-if="!isExploreHome" />
  </div>
</template>
