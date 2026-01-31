<template>
  <Button
    :variant="variant"
    @click.stop="authStore.isAuthenticated() && handleLike()"
    :class="[
      { 'pointer-events-none': !authStore.isAuthenticated() },
      isLiked
        ? 'text-primary border-primary hover:bg-primary/10'
        : 'text-muted-foreground hover:bg-secondary/10',
    ]"
  >
    <span v-if="!hideCount">{{ likes }}</span>
    <template v-if="isPending">
      <LoaderCircleIcon class="size-4 animate-spin" />
    </template>
    <template v-else>
      <Heart
        class="size-4 transition-transform duration-200"
        :class="{ 'scale-110': isLiked }"
        :fill="isLiked ? 'currentColor' : 'none'"
        stroke-width="2"
      />
    </template>
  </Button>
</template>

<script setup lang="ts">
import { Heart, LoaderCircleIcon } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useLikeSnippet } from '@/composables/useLikeSnippet'
import type { ButtonProps } from '@/components/ui/button'
import { Button } from '@/components/ui/button'

const authStore = useAuthStore()
const { mutate: toggleLike, isPending } = useLikeSnippet()

const props = withDefaults(
  defineProps<{
    snippetId: string
    likes: number
    isLiked: boolean
    hideCount?: boolean
    variant?: ButtonProps['variant']
  }>(),
  {
    variant: 'ghost',
    hideCount: false,
  },
)

const handleLike = () => {
  if (!props.snippetId || isPending.value) return

  const action = props.isLiked ? 'unlike' : 'like'
  toggleLike({
    snippetId: props.snippetId,
    action,
    currentLikes: props.likes,
    currentIsLiked: props.isLiked,
  })
}
</script>
