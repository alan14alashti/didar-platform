<script setup lang="ts">
withDefaults(defineProps<{ title: string; body?: string }>(), { body: '' })
</script>

<template>
  <div class="flex flex-col items-center justify-center gap-2 rounded-2xl border border-dashed border-gray-300 bg-white py-16 text-center">
    <p class="text-base font-semibold text-gray-700">{{ title }}</p>
    <p v-if="body" class="max-w-sm text-sm text-gray-500">{{ body }}</p>
    <slot />
  </div>
</template>
