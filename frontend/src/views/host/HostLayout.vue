<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>

<template>
  <div class="mx-auto max-w-7xl px-4 py-8 sm:px-6">
    <div class="grid gap-8 lg:grid-cols-[220px_1fr]">
      <aside class="lg:sticky lg:top-24 lg:self-start">
        <nav class="surface-card flex gap-2 overflow-x-auto p-2 text-sm font-semibold text-gray-600 lg:flex-col lg:overflow-visible">
          <router-link
            :to="{ name: 'host-dashboard' }"
            class="whitespace-nowrap rounded-xl px-4 py-2.5 transition-colors hover:bg-gray-50"
            active-class="!bg-brand-primary/10 !text-brand-primary"
          >
            📊 {{ t('host.dashboardTitle') }}
          </router-link>
          <router-link
            :to="{ name: 'host-experiences' }"
            class="whitespace-nowrap rounded-xl px-4 py-2.5 transition-colors hover:bg-gray-50"
            active-class="!bg-brand-primary/10 !text-brand-primary"
          >
            🧭 {{ t('host.myExperiences') }}
          </router-link>
          <router-link
            :to="{ name: 'host-bookings' }"
            class="whitespace-nowrap rounded-xl px-4 py-2.5 transition-colors hover:bg-gray-50"
            active-class="!bg-brand-primary/10 !text-brand-primary"
          >
            📩 {{ t('host.bookings') }}
          </router-link>
        </nav>
      </aside>

      <section>
        <router-view />
      </section>
    </div>
  </div>
</template>
