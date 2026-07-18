<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { Category, Experience } from '../api/types'
import ExperienceCard from './ExperienceCard.vue'
import AppIcon from './AppIcon.vue'

const props = defineProps<{
  category: Category
  experiences: Experience[]
  icon?: string
}>()

const { locale } = useI18n()
</script>

<template>
  <section class="space-y-3">
    <div class="flex items-center justify-between gap-3 px-1">
      <div class="flex items-center gap-2">
        <AppIcon v-if="icon" :icon="icon" size="md" icon-class="text-brand-primary" />
        <h2 class="text-lg font-bold text-gray-900">
          {{ locale === 'fa' ? category.title_fa : category.title_en }}
        </h2>
      </div>
      <span class="text-xs font-medium text-gray-400">{{ experiences.length }}</span>
    </div>

    <div class="flex gap-4 overflow-x-auto pb-2 scrollbar-thin">
      <ExperienceCard
        v-for="exp in experiences"
        :key="exp.id"
        :experience="exp"
        class="w-[17rem] shrink-0 sm:w-[19rem]"
      />
    </div>
  </section>
</template>
