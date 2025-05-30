<!-- components/ui/StatCard.vue -->
<template>
  <div class="group text-center">
    <div class="mb-2">
      <div
        class="text-foreground group-hover:text-primary text-3xl font-bold transition-colors sm:text-4xl"
      >
        {{ value }}
      </div>
    </div>
    <div class="text-foreground mb-1 text-sm font-medium">{{ label }}</div>
    <div class="text-muted-foreground text-xs">{{ description }}</div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  value: string
  label: string
  description: string
}

defineProps<Props>()
</script>
