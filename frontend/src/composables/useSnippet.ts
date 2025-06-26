import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { useSnippetsStore } from '@/stores/snippets'
import { snippetsService } from '@/services/snippets'
import { toast } from 'vue-sonner'
import { computed } from 'vue'

export function useSnippet(snippetId: string) {
  const snippetsStore = useSnippetsStore()
  const queryClient = useQueryClient()

  const { isPending, isError, error, refetch } = useQuery({
    queryKey: ['snippet', snippetId],
    queryFn: async () => {
      try {
        const data = await snippetsService.getSnippet(snippetId)
        snippetsStore.addSnippet(data)
        return data
      } catch (err) {
        throw err
      }
    },
  })

  // Get snippet from store (reactive)
  const snippet = computed(() => snippetsStore.getSnippetById(snippetId))

  // Update snippet mutation
  const { mutate: updateSnippet, isPending: isUpdating } = useMutation({
    mutationFn: (formData: { title: string; content: string; language?: string }) =>
      snippetsService.updateSnippet(snippetId, formData),
    onSuccess: (updatedSnippet) => {
      // Update in Pinia store
      snippetsStore.updateSnippet(updatedSnippet.id, updatedSnippet)

      // Also update TanStack Query cache for consistency
      queryClient.setQueryData(['snippet', updatedSnippet.id], updatedSnippet)

      queryClient.invalidateQueries({ queryKey: ['my-snippets'] })
      queryClient.invalidateQueries({ queryKey: ['liked-snippets'] })
      queryClient.invalidateQueries({ queryKey: ['saved-snippets'] })

      toast.success('Snippet updated successfully')
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : 'Failed to update snippet')
    },
  })

  // Delete snippet mutation
  const { mutate: deleteSnippet, isPending: isDeleting } = useMutation({
    mutationFn: () => snippetsService.deleteSnippet(snippetId),
    onSuccess: () => {
      // Remove from Pinia store
      snippetsStore.removeSnippet(snippetId)

      // Remove from TanStack Query cache
      queryClient.removeQueries({ queryKey: ['snippet', snippetId] })

      queryClient.invalidateQueries({ queryKey: ['my-snippets'] })
      queryClient.invalidateQueries({ queryKey: ['liked-snippets'] })
      queryClient.invalidateQueries({ queryKey: ['saved-snippets'] })

      toast.success('Snippet deleted successfully')
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : 'Failed to delete snippet')
    },
  })

  return {
    // State from Pinia store (reactive)
    snippet,
    isLoading: isPending,
    isError,
    error,

    // Actions
    updateSnippet,
    deleteSnippet,
    refetch,

    // Loading states
    isUpdating,
    isDeleting,
  }
}
