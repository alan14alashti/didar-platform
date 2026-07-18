<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    icon: string
    size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl'
    iconClass?: string
  }>(),
  { size: 'md' },
)

const sizeClass = computed(
  () =>
    ({
      xs: 'h-3 w-3',
      sm: 'h-3.5 w-3.5',
      md: 'h-4 w-4',
      lg: 'h-5 w-5',
      xl: 'h-6 w-6',
      '2xl': 'h-8 w-8',
    })[props.size],
)
</script>

<template>
  <Icon :icon="icon" :class="[sizeClass, props.iconClass]" aria-hidden="true" />
</template>
