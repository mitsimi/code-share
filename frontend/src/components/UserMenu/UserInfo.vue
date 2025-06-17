<template>
  <div class="flex items-center space-x-2">
    <Avatar class="size-8">
      <AvatarImage :src="authStore.user?.avatar" />
      <AvatarFallback>{{ authStore.user?.username[0].toUpperCase() }}</AvatarFallback>
    </Avatar>
    <div class="flex flex-col">
      <span class="text-sm font-medium">{{ authStore.user?.username }}</span>
      <span class="text-muted-foreground text-xs">{{ authStore.user?.email }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'

const authStore = useAuthStore()
</script>
