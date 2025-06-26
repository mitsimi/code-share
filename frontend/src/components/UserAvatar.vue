<template>
  <div :class="cn('flex items-center gap-4', props.class)">
    <Avatar :class="cn('ring-border ring-2', avatarClass)">
      <AvatarImage :src="user.avatar || ''" :alt="user.username" />
      <AvatarFallback class="bg-primary text-primary-foreground font-semibold">
        {{ user.username[0].toUpperCase() }}
      </AvatarFallback>
    </Avatar>
    <div class="flex flex-col">
      <span :class="cn(usernameClass, 'text-foreground font-semibold')">{{ user.username }}</span>
      <span v-if="subtitle" :class="cn(subtitleClass, 'text-muted-foreground text-sm')">
        {{ subtitle }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { cn } from '@/lib/utils'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'

interface User {
  username: string
  avatar?: string
}

interface Props {
  user: User
  subtitle?: string
  class?: string
  avatarClass?: string
  usernameClass?: string
  subtitleClass?: string
}

const props = withDefaults(defineProps<Props>(), {
  class: '',
  avatarClass: '',
  usernameClass: '',
  subtitleClass: '',
})
</script>
