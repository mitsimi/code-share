<template>
  <slot v-if="authStore.user?.id === props.authorId" />
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const props = defineProps<{
  authorId: string
}>()
</script>
