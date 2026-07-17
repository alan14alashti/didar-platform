<script setup lang="ts">
withDefaults(defineProps<{ label?: string }>(), { label: '' })
</script>

<template>
  <div class="flex flex-col items-center justify-center gap-3 py-16 text-gray-500">
    <span class="h-8 w-8 animate-spin rounded-full border-2 border-gray-200 border-t-brand-primary" />
    <p v-if="label" class="text-sm">{{ label }}</p>
  </div>
</template>
