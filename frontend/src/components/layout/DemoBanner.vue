<template>
  <div class="bg-primary/95 supports-[backdrop-filter]:bg-primary/60 border-b backdrop-blur">
    <div class="container mx-auto flex h-10 items-center justify-between px-4">
      <div class="flex items-center space-x-2">
        <CircleAlertIcon class="size-4" />
        <span class="text-foreground/80 text-sm font-semibold"> Demo Version</span>
        <span class="text-foreground/60 text-sm">This is a demonstration of the application</span>
      </div>
      <div class="flex items-center space-x-4">
        <Button
          variant="link"
          href="https://github.com/mitsimi/code-share"
          target="_blank"
          rel="noopener noreferrer"
          class="text-foreground/60 hover:text-foreground/80 flex items-center space-x-1 text-sm transition-colors"
          as="a"
        >
          <span>View Source</span>
          <ExternalLinkIcon class="size-4" />
        </Button>
        <Button
          variant="link"
          @click="showDemoInfo"
          class="text-foreground/60 hover:text-foreground/80 flex items-center space-x-1 text-sm transition-colors"
        >
          <span>Demo Info</span>
          <InfoIcon class="size-4" />
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { CircleAlertIcon, ExternalLinkIcon, InfoIcon } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import Button from '../ui/button/Button.vue'

const showDemoInfo = () => {
  toast.info('Demo Version Features', {
    description:
      'Limited functionality for demonstration purposes: Sample data included, Rate limiting applied, Session timeout after 30 minutes',
    duration: 5000,
  })
}
</script>
