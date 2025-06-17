<template>
  <Button
    variant="ghost"
    @click.stop="authStore.isAuthenticated() && toggleSave()"
    :class="[
      { 'pointer-events-none': !authStore.isAuthenticated() },
      isSaved
        ? 'text-primary border-primary hover:bg-primary/10'
        : 'text-muted-foreground border-secondary-foreground hover:bg-secondary/10',
    ]"
  >
    <template v-if="isLoading">
      <LoaderCircleIcon class="size-4 animate-spin" />
    </template>
    <template v-else>
      <BookmarkIcon
        class="size-4 transition-transform duration-200"
        :class="{ 'scale-110': isSaved }"
        :fill="isSaved ? 'currentColor' : 'none'"
        :stroke="isSaved ? 'currentColor' : 'currentColor'"
        stroke-width="2"
      />
    </template>
  </Button>
</template>

<script setup lang="ts">
import { BookmarkIcon, LoaderCircleIcon } from 'lucide-vue-next'
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
  isSaved: boolean
}>()

const isLoading = ref(false)

const { mutate: updateSave } = useMutation<
  Snippet,
  Error,
  { snippetId: string; action: 'save' | 'unsave' }
>({
  mutationKey: ['saveMutation', props.snippetId],
  mutationFn: async ({ snippetId, action }) => {
    isLoading.value = true
    console.log(`Starting ${action} mutation for snippet:`, snippetId)
    try {
      const { data, error } = await useFetch<Snippet>(
        `/snippets/${snippetId}/save?action=${action}`,
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
      return oldData.map((snippet) => (snippet.id === updatedSnippet.id ? updatedSnippet : snippet))
    })

    // Invalidate saved snippets query to trigger a refetch
    queryClient.invalidateQueries({ queryKey: ['saved-snippets'] })
    queryClient.invalidateQueries({ queryKey: ['liked-snippets'] })
    queryClient.invalidateQueries({ queryKey: ['my-snippets'] })
  },
  onError: (error) => {
    console.error('Save mutation failed:', error)
    toast.error(error.message || 'Please try again')
  },
})

const toggleSave = () => {
  if (!props.snippetId) {
    console.error('Cannot toggle save: snippetId is missing')
    return
  }

  if (isLoading.value) {
    return
  }

  const action = props.isSaved ? 'unsave' : 'save'
  updateSave({ snippetId: props.snippetId, action })
}
</script>
