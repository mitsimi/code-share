<template>
  <Button
    :variant="variant"
    @click.stop="authStore.isAuthenticated() && toggleLike()"
    :class="[
      { 'pointer-events-none': !authStore.isAuthenticated() },
      isLiked
        ? 'text-primary border-primary hover:bg-primary/10'
        : 'text-muted-foreground hover:bg-secondary/10',
    ]"
  >
    <span v-if="!hideCount">{{ likes }}</span>
    <template v-if="isLoading">
      <LoaderCircleIcon class="size-4 animate-spin" />
    </template>
    <template v-else>
      <Heart
        class="size-4 transition-transform duration-200"
        :class="{ 'scale-110': isLiked }"
        :fill="isLiked ? 'currentColor' : 'none'"
        :stroke="isLiked ? 'currentColor' : 'currentColor'"
        stroke-width="2"
      />
    </template>
  </Button>
</template>

<script setup lang="ts">
import { Heart, LoaderCircleIcon } from 'lucide-vue-next'
import { type Snippet } from '@/types'
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useSnippetsStore } from '@/stores/snippets'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { snippetsService } from '@/services/snippets'
import { toast } from 'vue-sonner'
import type { ButtonProps } from '@/components/ui/button'
import { Button } from '@/components/ui/button'

const authStore = useAuthStore()
const snippetsStore = useSnippetsStore()
const queryClient = useQueryClient()

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

const isLoading = ref(false)

const { mutate: updateLike } = useMutation<
  Snippet,
  Error,
  { snippetId: string; action: 'like' | 'unlike' }
>({
  mutationKey: ['likeMutation', props.snippetId],
  mutationFn: async ({ snippetId, action }) => {
    isLoading.value = true
    try {
      return await snippetsService.toggleLike(snippetId, action)
    } finally {
      isLoading.value = false
    }
  },
  onSuccess: (updatedSnippet) => {
    // Update the snippet in the details view
    queryClient.setQueryData(['snippet', updatedSnippet.id], updatedSnippet)

    // The Pinia store is already updated optimistically,
    // but we sync with the server response to ensure consistency
    snippetsStore.updateSnippet(updatedSnippet.id, {
      isLiked: updatedSnippet.isLiked,
      likes: updatedSnippet.likes,
    })

    // Invalidate related queries
    queryClient.invalidateQueries({ queryKey: ['liked-snippets'] })
    queryClient.invalidateQueries({ queryKey: ['my-snippets'] })
    queryClient.invalidateQueries({ queryKey: ['saved-snippets'] })
  },
  onError: (error, { snippetId, action }) => {
    // Revert optimistic update on error
    const revertAction = action === 'like' ? 'unlike' : 'like'
    snippetsStore.handleUserAction(snippetId, revertAction, action === 'unlike')

    console.error('Like mutation failed:', error)
    toast.error(error.message || 'Please try again')
  },
})

const toggleLike = () => {
  if (!props.snippetId) {
    console.error('Cannot toggle like: snippetId is missing')
    return
  }

  if (isLoading.value) {
    return
  }

  const action = props.isLiked ? 'unlike' : 'like'

  // Optimistic update - immediate UI feedback
  snippetsStore.handleUserAction(props.snippetId, action, action === 'like')

  // Then make the API call
  updateLike({ snippetId: props.snippetId, action })
}
</script>
