<template>
  <div class="bg-background flex min-h-screen flex-col">
    <AppHeader />
    <main class="grow">
      <RouterView />
    </main>
    <AppFooter />
  </div>
  <Toaster richColors position="bottom-left" />
  <VueQueryDevtools />
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import { Toaster } from 'vue-sonner'
import { VueQueryDevtools } from '@tanstack/vue-query-devtools'
import AppFooter from './components/AppFooter.vue'
import AppHeader from './components/AppHeader.vue'
</script>

<style>
/* Sonner toast customization */
:root {
  --toast-bg: #fcf7de;
  --toast-border: #000000;
  --toast-shadow: 4px 4px 0 0 #000000;
}

[data-sonner-toast] {
  border: 4px solid var(--toast-border) !important;
  background: var(--toast-bg) !important;
  box-shadow: var(--toast-shadow) !important;
  border-radius: 0.5rem !important;
  font-weight: 600 !important;
}

[data-sonner-toast][data-type='success'] {
  background: var(--secondary) !important;
  color: var(--secondary-foreground) !important;
}

[data-sonner-toast][data-type='error'] {
  background: var(--destructive) !important;
  color: var(--destructive-foreground) !important;
}

[data-sonner-toast][data-type='info'] {
  background: var(--primary) !important;
  color: var(--primary-foreground) !important;
}

[data-sonner-toast] [data-close-button] {
  border: 2px solid var(--toast-border) !important;
  background: var(--toast-bg) !important;
  color: var(--toast-border) !important;
  border-radius: 0.25rem !important;
  padding: 0.25rem !important;
  margin-left: 0.5rem !important;
}

[data-sonner-toast] [data-close-button]:hover {
  background: var(--toast-border) !important;
  color: var(--toast-bg) !important;
}

[data-sonner-toast] [data-description] {
  color: inherit !important;
  opacity: 0.9 !important;
}
</style>
