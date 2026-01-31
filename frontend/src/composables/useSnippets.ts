import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { snippetsService } from '@/services/snippets'
import { queryKeys } from '@/composables/queryKeys'
import { toast } from 'vue-sonner'
import type { Snippet } from '@/types'

export function useSnippets() {
  const queryClient = useQueryClient()

  const getList = useQuery({
    queryKey: queryKeys.lists(),
    queryFn: snippetsService.getSnippets,
  })

  const getSnippet = (snippetId: string) =>
    useQuery({
      queryKey: queryKeys.detail(snippetId),
      queryFn: () => snippetsService.getSnippet(snippetId),
    })

  const createSnippet = useMutation({
    mutationFn: snippetsService.createSnippet,
    onSuccess: (newSnippet) => {
      queryClient.setQueryData(queryKeys.detail(newSnippet.id), newSnippet)
      queryClient.invalidateQueries({ queryKey: queryKeys.lists() })
      toast.success(`"${newSnippet.title}" has been added successfully!`)
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : 'Failed to create snippet')
    },
  })

  const updateSnippet = useMutation({
    mutationFn: ({
      snippetId,
      formData,
    }: {
      snippetId: string
      formData: { title: string; content: string; language?: string }
    }) => snippetsService.updateSnippet(snippetId, formData),
    onMutate: async ({ snippetId, formData }) => {
      await queryClient.cancelQueries({ queryKey: queryKeys.detail(snippetId) })
      const previousSnippet = queryClient.getQueryData(queryKeys.detail(snippetId))

      queryClient.setQueryData(queryKeys.detail(snippetId), (old: Snippet | undefined) => {
        if (!old) return old
        return { ...old, ...formData }
      })

      return { previousSnippet }
    },
    onError: (error, vars, context) => {
      if (context?.previousSnippet) {
        queryClient.setQueryData(queryKeys.detail(vars.snippetId), context.previousSnippet)
      }
      toast.error(error instanceof Error ? error.message : 'Failed to update snippet')
    },
    onSuccess: (updatedSnippet) => {
      queryClient.setQueryData(queryKeys.detail(updatedSnippet.id), updatedSnippet)
      queryClient.invalidateQueries({ queryKey: queryKeys.lists() })
      queryClient.invalidateQueries({ queryKey: queryKeys.my() })
      toast.success('Snippet updated successfully')
    },
  })

  const deleteSnippet = useMutation({
    mutationFn: ({ snippetId }: { snippetId: string }) => snippetsService.deleteSnippet(snippetId),
    onSuccess: (_, vars) => {
      queryClient.removeQueries({ queryKey: queryKeys.detail(vars.snippetId) })
      queryClient.invalidateQueries({ queryKey: queryKeys.lists() })
      queryClient.invalidateQueries({ queryKey: queryKeys.my() })
      toast.success('Snippet deleted successfully')
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : 'Failed to delete snippet')
    },
  })

  return {
    getList,
    getSnippet,
    createSnippet,
    updateSnippet,
    deleteSnippet,
  }
}
