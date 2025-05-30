<!-- components/ui/TestimonialCard.vue -->
<template>
  <div class="bg-card flex flex-col rounded-lg border p-6 transition-shadow hover:shadow-lg">
    <div class="mb-4 grow">
      <p class="text-card-foreground leading-relaxed italic">"{{ quote }}"</p>
    </div>
    <div class="flex items-center gap-3">
      <div
        class="bg-primary text-primary-foreground flex h-10 w-10 items-center justify-center rounded-full text-sm font-medium"
      >
        {{ avatar }}
      </div>
      <div>
        <div class="text-card-foreground text-sm font-medium">{{ name }}</div>
        <div class="text-muted-foreground text-xs">{{ role }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  quote: string
  name: string
  role: string
  avatar: string
}

defineProps<Props>()
</script>
