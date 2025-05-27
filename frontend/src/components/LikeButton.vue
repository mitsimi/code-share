<template>
  <Button
    variant="outline"
    @click.stop="authStore.isAuthenticated() && toggleLike()"
    :class="{ 'pointer-events-none': !authStore.isAuthenticated() }"
  >
    <span>{{ likes }}</span>
    <Heart class="size-5" :fill="is_liked ? 'red' : 'none'" />
  </Button>
</template>

<script setup lang="ts">
import { Heart } from 'lucide-vue-next'
import { Button } from './ui/button'
import { type Snippet } from '@/types'
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useCustomFetch } from '@/composables/useCustomFetch'
import { toast } from 'vue-sonner'

const authStore = useAuthStore()
const queryClient = useQueryClient()

const props = defineProps<{
  snippetId: string
  likes: number
  is_liked: boolean
}>()

const isLoading = ref(false)

// Inline implementation of like functionality instead of using useLikeSnippet
const { mutate: updateLike } = useMutation({
  mutationKey: ['likeMutation', props.snippetId],
  mutationFn: async ({ snippetId, action }: { snippetId: string; action: 'like' | 'unlike' }) => {
    console.log(`Starting ${action} mutation for snippet:`, snippetId)
    try {
      const { data, error } = await useCustomFetch<Snippet>(
        `/snippets/${snippetId}/like?action=${action}`,
        {
          method: 'PATCH',
        },
      ).json()

      console.log('Response data:', data.value)
      console.log('Response error:', error.value)

      if (error.value) {
        console.error('Error in mutation:', error.value)
        throw new Error(`Failed to ${action}: ${error.value.message || 'Unknown error'}`)
      }

      if (!data.value) {
        console.error('No data received from server')
        throw new Error('No data received from server')
      }

      return data.value
    } catch (err) {
      console.error('Caught error in mutation:', err)
      throw err
    }
  },
  onSuccess: (updatedSnippet) => {
    // Update the snippet in the details view
    queryClient.setQueryData(['snippet', updatedSnippet.id], updatedSnippet)

    // Update the snippet in the list view
    queryClient.setQueryData(['snippets'], (oldData: Snippet[] | undefined) => {
      if (!oldData) return [updatedSnippet]
      return oldData.map((snippet) =>
        snippet.id === updatedSnippet.id
          ? { ...updatedSnippet, isLiked: !snippet.is_liked }
          : snippet,
      )
    })

    /*toast.success(updatedSnippet.isLiked ? 'Added to favorites' : 'Removed from favorites', {
      description: `"${updatedSnippet.title}" ${updatedSnippet.isLiked ? 'added to' : 'removed from'} your favorites`,
    })*/
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

  const action = props.is_liked ? 'unlike' : 'like'
  updateLike({ snippetId: props.snippetId, action })
}
</script>
