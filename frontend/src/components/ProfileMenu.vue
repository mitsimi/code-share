<template>
  <DropdownMenu>
    <DropdownMenuTrigger asChild>
      <Button variant="ghost" size="sm">
        <span>{{ authStore.user?.username }}</span>
        <Avatar class="size-8">
          <AvatarImage :src="authStore.user?.avatar" />
          <AvatarFallback>{{ authStore.user?.username[0].toUpperCase() }}</AvatarFallback>
        </Avatar>
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent class="w-56" align="end">
      <DropdownMenuLabel>
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
      </DropdownMenuLabel>
      <DropdownMenuItem @click="$router.push('/profile')">
        <UserIcon class="mr-2 size-4" />
        Profile
      </DropdownMenuItem>
      <DropdownMenuSeparator />
      <DropdownMenuItem @click="emit('logout')">
        <LogOutIcon class="mr-2 size-4" />
        Logout
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>

<script setup lang="ts">
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu'
import { useAuthStore } from '@/stores/auth'
import { LogOutIcon, UserIcon } from 'lucide-vue-next'

const authStore = useAuthStore()

const emit = defineEmits<{
  (e: 'logout'): void
}>()
</script>
