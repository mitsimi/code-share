<template>
  <slot v-if="!authStore.isAuthenticated()" />
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
</script>
