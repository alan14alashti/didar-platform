<script setup lang="ts">
import AppIcon from './AppIcon.vue'

defineProps<{ rating: number; size?: 'sm' | 'md' }>()
</script>

<template>
  <span class="inline-flex items-center gap-1" :class="size === 'md' ? 'text-base' : 'text-xs'">
    <AppIcon icon="heroicons-outline:star" :size="size === 'md' ? 'md' : 'sm'" icon-class="text-amber-500" />
    <span class="font-semibold text-gray-800">{{ rating > 0 ? rating.toFixed(1) : '-' }}</span>
  </span>
</template>
