<script setup lang="ts">
defineProps<{ rating: number; size?: 'sm' | 'md' }>()
</script>

<template>
  <span class="inline-flex items-center gap-1" :class="size === 'md' ? 'text-base' : 'text-xs'">
    <svg viewBox="0 0 20 20" fill="currentColor" class="h-4 w-4 text-amber-500">
      <path d="M10 1.5l2.7 5.9 6.3.7-4.7 4.5 1.2 6.4L10 15.9l-5.5 3.1 1.2-6.4L1 8.1l6.3-.7L10 1.5z" />
    </svg>
    <span class="font-semibold text-gray-800">{{ rating > 0 ? rating.toFixed(1) : '-' }}</span>
  </span>
</template>
