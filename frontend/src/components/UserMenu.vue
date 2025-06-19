<template>
  <div>
    <!-- Desktop Version (Dropdown) -->
    <DropdownMenu v-if="!isMobile">
      <DropdownMenuTrigger asChild>
        <Button variant="ghost">
          <UserAvatar :user="authStore.user" class="flex-row-reverse gap-2" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent class="w-56" align="end">
        <MenuItems :is-mobile="false" @logout="emit('logout')" />
      </DropdownMenuContent>
    </DropdownMenu>

    <!-- Mobile Version (Buttons) -->
    <div v-else>
      <UserInfo class="px-2 py-2" />
      <MenuItems :is-mobile="true" @logout="emit('logout')" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { DropdownMenu, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { Button } from '@/components/ui/button'
import UserInfo from './UserMenu/UserInfo.vue'
import MenuItems from './UserMenu/MenuItems.vue'

const props = defineProps<{
  isMobile?: boolean
}>()

const authStore = useAuthStore()

const emit = defineEmits<{
  (e: 'logout'): void
}>()
</script>
