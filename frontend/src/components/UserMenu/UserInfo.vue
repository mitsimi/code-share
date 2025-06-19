<template>
  <UserAvatar
    v-if="authStore.user"
    :user="authStore.user"
    :subtitle="authStore.user.email"
    avatar-class="size-8"
  />
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import UserAvatar from '@/components/UserAvatar.vue'

const authStore = useAuthStore()
</script>
