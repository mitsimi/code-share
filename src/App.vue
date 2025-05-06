<template>
  <div class="flex min-h-screen flex-col">
    <AppHeader />
    <main class="flex-1">
      <router-view />
    </main>
    <AppFooter />
    <ToastNotification />
  </div>
</template>

<script setup lang="ts">
import AppHeader from './components/AppHeader.vue'
import AppFooter from './components/AppFooter.vue'
import ToastNotification from './components/ToastNotification.vue'
</script>

<style>
/* Custom scrollbar styling */
::-webkit-scrollbar {
  width: 16px;
  height: 16px;
}

::-webkit-scrollbar-track {
  background: #fcf7de;
  border-left: 3px solid #000000;
}

::-webkit-scrollbar-thumb {
  background: var(--color-accent);
  border: 3px solid #000000;
  box-shadow: inset 2px 2px 0 rgba(255, 255, 255, 0.3);
}

::-webkit-scrollbar-corner {
  background: #fcf7de;
}
</style>
