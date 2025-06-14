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
import { useFetch } from '@/composables/useCustomFetch'
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
    console.log(`Starting ${action} mutation for snippet:`, snippetId)
    try {
      const { data, error } = await useFetch<Snippet>(
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
          ? { ...updatedSnippet, isLiked: !snippet.isLiked }
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

  const action = props.isLiked ? 'unlike' : 'like'
  updateLike({ snippetId: props.snippetId, action })
}
</script>
