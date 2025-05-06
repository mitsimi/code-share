<template>
  <RouterView />
  <ToastNotification />
</template>

<script setup lang="ts">
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
