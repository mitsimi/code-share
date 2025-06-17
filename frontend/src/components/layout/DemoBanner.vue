<template>
  <div class="bg-primary/95 supports-[backdrop-filter]:bg-primary/60 border-b backdrop-blur">
    <div
      class="container mx-auto flex h-auto items-center justify-between px-4 py-2 sm:h-10 sm:py-0"
    >
      <div class="text-primary-foreground/80 flex items-center space-x-2">
        <CircleAlertIcon class="size-4" />
        <span class="text-sm font-semibold"> Demo Version</span>
        <span class="text-primary-foreground/60 hidden text-sm sm:inline"
          >This is a demonstration of the application</span
        >
      </div>
      <div class="flex items-center gap-x-4">
        <Button
          variant="link"
          href="https://github.com/mitsimi/code-share"
          target="_blank"
          rel="noopener noreferrer"
          class="text-primary-foreground/60 hover:text-primary-foreground/80 flex items-center space-x-1 p-0 text-sm transition-colors"
          as="a"
        >
          <span>View Source</span>
          <ExternalLinkIcon class="size-4" />
        </Button>
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger
              class="text-primary-foreground/60 hover:text-primary-foreground/80 flex hidden items-center space-x-1 p-0 text-sm transition-colors hover:underline hover:underline-offset-4 sm:inline"
            >
              <span>Demo Info</span>
              <InfoIcon class="size-4" />
            </TooltipTrigger>
            <TooltipContent class="space-y-2">
              <h3 class="text-sm font-medium">Limited functionality for demonstration purposes:</h3>
              <ul class="list-disc space-y-1 pl-4 text-sm">
                <li>Sample data included</li>
                <li>Rate limiting applied</li>
                <li>Signups disabled</li>
              </ul>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { CircleAlertIcon, ExternalLinkIcon, InfoIcon } from 'lucide-vue-next'
import Button from '@/components/ui/button/Button.vue'
import Tooltip from '@/components/ui/tooltip/Tooltip.vue'
</script>
