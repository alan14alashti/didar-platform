<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>

<template>
  <div class="relative flex min-h-[70vh] max-w-lg flex-col items-center justify-center overflow-hidden px-4 text-center mx-auto">
    <div class="blob animate-blob h-64 w-64 bg-brand-primary/15 -top-10" />
    <p class="text-gradient relative text-7xl font-black">404</p>
    <p class="relative mt-3 text-lg font-semibold text-gray-700">{{ t('experience.notFound') }}</p>
    <router-link to="/" class="btn-primary relative mt-6">
      {{ t('nav.home') }}
    </router-link>
  </div>
</template>
