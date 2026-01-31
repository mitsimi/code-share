import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { useSnippetsStore } from '@/stores/snippets'
import { snippetsService } from '@/services/snippets'
import { queryKeys } from '@/composables/queryKeys'
import { toast } from 'vue-sonner'
import { computed } from 'vue'

export function useSnippet(snippetId: string) {
  const snippetsStore = useSnippetsStore()
  const queryClient = useQueryClient()
  const snippetQueryKey = queryKeys.snippet(snippetId)
  const mySnippetsQueryKey = queryKeys.mySnippets()
  const likedSnippetsQueryKey = queryKeys.likedSnippets()
  const savedSnippetsQueryKey = queryKeys.savedSnippets()

  const { isPending, isError, error, refetch } = useQuery({
    queryKey: snippetQueryKey,
    queryFn: async () => {
      const data = await snippetsService.getSnippet(snippetId)
      snippetsStore.addSnippet(data)
      return data
    },
  })

  const snippet = computed(() => snippetsStore.getSnippetById(snippetId))

  const { mutate: updateSnippet, isPending: isUpdating } = useMutation({
    mutationFn: (formData: { title: string; content: string; language?: string }) =>
      snippetsService.updateSnippet(snippetId, formData),
    onSuccess: (updatedSnippet) => {
      snippetsStore.updateSnippet(updatedSnippet.id, updatedSnippet)
      queryClient.setQueryData(queryKeys.snippet(updatedSnippet.id), updatedSnippet)

      queryClient.invalidateQueries({ queryKey: mySnippetsQueryKey })
      queryClient.invalidateQueries({ queryKey: likedSnippetsQueryKey })
      queryClient.invalidateQueries({ queryKey: savedSnippetsQueryKey })

      toast.success('Snippet updated successfully')
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : 'Failed to update snippet')
    },
  })

  const { mutate: deleteSnippet, isPending: isDeleting } = useMutation({
    mutationFn: () => snippetsService.deleteSnippet(snippetId),
    onSuccess: () => {
      snippetsStore.removeSnippet(snippetId)
      queryClient.removeQueries({ queryKey: snippetQueryKey })

      queryClient.invalidateQueries({ queryKey: mySnippetsQueryKey })
      queryClient.invalidateQueries({ queryKey: likedSnippetsQueryKey })
      queryClient.invalidateQueries({ queryKey: savedSnippetsQueryKey })

      toast.success('Snippet deleted successfully')
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : 'Failed to delete snippet')
    },
  })

  return {
    snippet,
    isLoading: isPending,
    isError,
    error,
    updateSnippet,
    deleteSnippet,
    refetch,
    isUpdating,
    isDeleting,
  }
}
