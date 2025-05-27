<template>
  <Button
    variant="outline"
    @click.stop="authStore.isAuthenticated() && toggleLike()"
    :class="{ 'pointer-events-none': !authStore.isAuthenticated() }"
  >
    <span>{{ likes }}</span>
    <Heart v-if="!isLoading" class="size-5" :fill="is_liked ? 'red' : 'none'" />
    <div v-else class="flex items-center gap-2">
      <svg
        class="text-primary h-5 w-5 animate-spin"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
      >
        <circle
          class="opacity-25"
          cx="12"
          cy="12"
          r="10"
          stroke="currentColor"
          stroke-width="4"
        ></circle>
        <path
          class="opacity-75"
          fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
        ></path>
      </svg>
    </div>
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

const { mutate: updateLike } = useMutation<
  Snippet,
  Error,
  { snippetId: string; action: 'like' | 'unlike' }
>({
  mutationKey: ['likeMutation', props.snippetId],
  mutationFn: async ({ snippetId, action }) => {
    isLoading.value = true
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
      return oldData.map((snippet) =>
        snippet.id === updatedSnippet.id
          ? { ...updatedSnippet, is_liked: !snippet.is_liked }
          : snippet,
      )
    })
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
