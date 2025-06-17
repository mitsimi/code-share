<template>
  <div>
    <template v-if="!isMobile">
      <DropdownMenuLabel>
        <UserInfo />
      </DropdownMenuLabel>
      <DropdownMenuItem @click="$router.push('/profile?tab=settings')">
        <SettingsIcon class="mr-2 size-4" />
        Settings
      </DropdownMenuItem>
      <DropdownMenuItem @click="$router.push('/profile?tab=mine')">
        <CodeIcon class="mr-2 size-4" />
        My Snippets
      </DropdownMenuItem>
      <DropdownMenuItem @click="$router.push('/profile?tab=liked')">
        <HeartIcon class="mr-2 size-4" />
        Liked Snippets
      </DropdownMenuItem>
      <DropdownMenuItem @click="$router.push('/profile?tab=saved')">
        <BookmarkIcon class="mr-2 size-4" />
        Saved Snippets
      </DropdownMenuItem>
      <DropdownMenuSeparator />
      <DropdownMenuItem @click="emit('logout')">
        <LogOutIcon class="mr-2 size-4" />
        Logout
      </DropdownMenuItem>
    </template>
    <template v-else>
      <Button
        variant="ghost"
        class="w-full justify-start"
        @click="handleClick('/profile?tab=settings')"
      >
        <SettingsIcon class="mr-2 h-4 w-4" />
        Settings
      </Button>
      <Button
        variant="ghost"
        class="w-full justify-start"
        @click="handleClick('/profile?tab=mine')"
      >
        <CodeIcon class="mr-2 h-4 w-4" />
        My Snippets
      </Button>
      <Button
        variant="ghost"
        class="w-full justify-start"
        @click="handleClick('/profile?tab=liked')"
      >
        <HeartIcon class="mr-2 h-4 w-4" />
        Liked Snippets
      </Button>
      <Button
        variant="ghost"
        class="w-full justify-start"
        @click="handleClick('/profile?tab=saved')"
      >
        <BookmarkIcon class="mr-2 h-4 w-4" />
        Saved Snippets
      </Button>
      <Button variant="ghost" class="w-full justify-start" @click="handleLogout">
        <LogOutIcon class="mr-2 h-4 w-4" />
        Logout
      </Button>
    </template>
  </div>
</template>

<script setup lang="ts">
import { SettingsIcon, CodeIcon, HeartIcon, BookmarkIcon, LogOutIcon } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import {
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuItem,
} from '@/components/ui/dropdown-menu'
import UserInfo from './UserInfo.vue'
import { useRouter } from 'vue-router'

const props = defineProps<{
  isMobile?: boolean
}>()

const emit = defineEmits<{
  (e: 'logout'): void
  (e: 'close'): void
}>()

const router = useRouter()

const handleClick = (route: string) => {
  router.push(route)
  if (props.isMobile) {
    emit('close')
  }
}

const handleLogout = () => {
  emit('logout')
  if (props.isMobile) {
    emit('close')
  }
}
</script>
