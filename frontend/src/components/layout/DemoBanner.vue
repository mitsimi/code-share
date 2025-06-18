<template>
  <div class="bg-primary/95 supports-[backdrop-filter]:bg-primary/60 border-b backdrop-blur">
    <div
      class="mx-auto flex h-10 max-w-7xl items-center justify-between px-4 py-0 md:h-auto md:px-8 md:py-2"
    >
      <div class="text-card-foreground/80 flex items-center space-x-2">
        <CircleAlertIcon class="size-4" />
        <span class="text-sm font-semibold"> Demo Version</span>
        <span class="text-card-foreground/60 hidden text-sm sm:inline"
          >This is a demonstration of the application</span
        >
      </div>
      <div class="flex items-center gap-x-4 px-2 md:px-4">
        <Button
          variant="link"
          href="https://github.com/mitsimi/code-share"
          target="_blank"
          rel="noopener noreferrer"
          class="text-card-foreground/60 hover:text-card-foreground/80 flex items-center gap-x-1 p-0 text-sm transition-colors"
          as="a"
        >
          <span>View Source</span>
          <ExternalLinkIcon class="size-4" />
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { CircleAlertIcon, ExternalLinkIcon } from 'lucide-vue-next'
import Button from '@/components/ui/button/Button.vue'
</script>
