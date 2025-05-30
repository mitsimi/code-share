<!-- components/ui/FeatureCard.vue -->
<template>
  <div
    class="group bg-card relative overflow-hidden rounded-lg border p-6 transition-all duration-300 hover:shadow-lg"
  >
    <div class="mb-4 flex items-center gap-4">
      <div
        class="bg-primary/10 group-hover:bg-primary/20 flex h-12 w-12 items-center justify-center rounded-lg transition-colors"
      >
        <component :is="icon" class="text-primary h-6 w-6" />
      </div>
    </div>
    <h3 class="text-card-foreground mb-2 text-lg font-semibold">{{ title }}</h3>
    <p class="text-muted-foreground text-sm leading-relaxed">{{ description }}</p>
  </div>
</template>

<script setup lang="ts">
interface Props {
  icon: any
  title: string
  description: string
}

defineProps<Props>()
</script>
