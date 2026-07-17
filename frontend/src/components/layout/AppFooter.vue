<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useBrandStore } from '../../stores/brand'

const { t } = useI18n()
const brandStore = useBrandStore()
</script>

<template>
  <footer class="border-t border-gray-100 bg-white">
    <div class="mx-auto max-w-7xl px-4 py-10 text-center sm:px-6">
      <span class="mx-auto grid h-10 w-10 place-items-center rounded-2xl brand-gradient text-base font-extrabold text-white">
        {{ brandStore.brand?.name?.[0] ?? 'D' }}
      </span>
      <p class="mt-3 font-bold text-gray-800">{{ brandStore.brand?.name ?? 'Didar' }}</p>
      <p class="mt-1 text-sm text-gray-500">{{ brandStore.brand?.tagline ?? t('footer.tagline') }}</p>
      <p class="mt-5 text-xs text-gray-400">© {{ new Date().getFullYear() }} {{ brandStore.brand?.name ?? 'Didar' }} — {{ t('footer.rights') }}</p>
    </div>
  </footer>
</template>
