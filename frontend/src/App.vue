<template>
  <div class="bg-background flex min-h-screen flex-col">
    <AppHeader />
    <main class="my-12 grow">
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
import AppFooter from './components/layout/AppFooter.vue'
import AppHeader from './components/layout/AppHeader.vue'
import { useUserActions } from '@/composables/useSnippetSubscription'

// Initialize user actions subscription for syncing likes/saves across tabs
useUserActions()
</script>
