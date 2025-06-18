<template>
  <Button
    variant="ghost"
    @click.stop="authStore.isAuthenticated() && toggleLike()"
    :class="[
      { 'pointer-events-none': !authStore.isAuthenticated() },
      isLiked
        ? 'text-primary border-primary hover:bg-primary/10'
        : 'text-muted-foreground border-secondary-foreground hover:bg-secondary/10',
    ]"
  >
    <span>{{ likes }}</span>
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
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { snippetsService } from '@/services/snippets'
import { toast } from 'vue-sonner'

const authStore = useAuthStore()
const queryClient = useQueryClient()

const props = defineProps<{
  snippetId: string
  likes: number
  isLiked: boolean
}>()

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

    // Update the snippet in the list view
    queryClient.setQueryData(['snippets'], (oldData: Snippet[] | undefined) => {
      if (!oldData) return [updatedSnippet]
      return oldData.map((snippet) => (snippet.id === updatedSnippet.id ? updatedSnippet : snippet))
    })

    // Invalidate liked snippets query to trigger a refetch
    queryClient.invalidateQueries({ queryKey: ['liked-snippets'] })
    queryClient.invalidateQueries({ queryKey: ['my-snippets'] })
    queryClient.invalidateQueries({ queryKey: ['saved-snippets'] })
  },
  onError: (error) => {
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
  updateLike({ snippetId: props.snippetId, action })
}
</script>
