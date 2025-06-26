import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { useSnippetsStore } from '@/stores/snippets'
import { snippetsService } from '@/services/snippets'
import { toast } from 'vue-sonner'
import { computed } from 'vue'

export function useSnippets() {
  const snippetsStore = useSnippetsStore()
  const queryClient = useQueryClient()

  // Use TanStack Query for fetching, but sync data to Pinia store
  const { isPending, isError, error, refetch } = useQuery({
    queryKey: ['snippets'],
    queryFn: async () => {
      snippetsStore.setLoading(true)
      try {
        const data = await snippetsService.getSnippets()
        snippetsStore.setSnippets(data)
        snippetsStore.setLoading(false)
        return data
      } catch (err) {
        snippetsStore.setLoading(false)
        snippetsStore.setError(err instanceof Error ? err.message : 'Failed to fetch snippets')
        throw err
      }
    },
  })

  // Create snippet mutation
  const { mutate: createSnippet, isPending: isCreating } = useMutation({
    mutationFn: snippetsService.createSnippet,
    onSuccess: (newSnippet) => {
      // Add to Pinia store
      snippetsStore.addSnippet(newSnippet)

      // Also update TanStack Query cache for consistency
      queryClient.setQueryData(['snippets'], (old: any) => {
        if (!old) return [newSnippet]
        return [newSnippet, ...old]
      })

      toast.success('Snippet added successfully', {
        description: `"${newSnippet.title}" has been added successfully!`,
      })
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to create snippet. Please try again.',
      )
    },
  })

  return {
    // State from Pinia store (reactive)
    snippets: computed(() => snippetsStore.snippets),
    isLoading: computed(() => snippetsStore.isLoading || isPending.value),
    isError: computed(() => isError.value || !!snippetsStore.error),
    error: computed(() => error.value || snippetsStore.error),

    // Getters
    getSnippetById: snippetsStore.getSnippetById,
    snippetsCount: snippetsStore.snippetsCount,

    // Actions
    createSnippet,
    refetch,

    // Loading states
    isCreating,
  }
}
