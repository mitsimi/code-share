import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { useSnippetsStore } from '@/stores/snippets'
import { snippetsService } from '@/services/snippets'
import { queryKeys } from '@/composables/queryKeys'
import { toast } from 'vue-sonner'
import { computed } from 'vue'

export function useSnippets() {
  const snippetsStore = useSnippetsStore()
  const queryClient = useQueryClient()
  const snippetsQueryKey = queryKeys.snippets()

  const { isPending, isError, error, refetch } = useQuery({
    queryKey: snippetsQueryKey,
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

  const { mutate: createSnippet, isPending: isCreating } = useMutation({
    mutationFn: snippetsService.createSnippet,
    onSuccess: (newSnippet) => {
      snippetsStore.addSnippet(newSnippet)

      queryClient.setQueryData(snippetsQueryKey, (old: Array<typeof newSnippet> | undefined) => {
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
    snippets: computed(() => snippetsStore.snippets),
    isLoading: computed(() => snippetsStore.isLoading || isPending.value),
    isError: computed(() => isError.value || !!snippetsStore.error),
    error: computed(() => error.value || snippetsStore.error),
    getSnippetById: snippetsStore.getSnippetById,
    snippetsCount: snippetsStore.snippetsCount,
    createSnippet,
    refetch,
    isCreating,
  }
}
