import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { snippetsService } from '@/services/snippets'
import { queryKeys } from '@/composables/queryKeys'
import { toast } from 'vue-sonner'
import type { Snippet } from '@/types'

interface SaveVariables {
  snippetId: string
  action: 'save' | 'unsave'
  currentIsSaved: boolean
}

interface SaveContext {
  previousSnippet: Snippet | undefined
}

export function useSaveSnippet() {
  const queryClient = useQueryClient()

  return useMutation<Snippet, Error, SaveVariables, SaveContext>({
    mutationKey: ['save'],
    mutationFn: ({ snippetId, action }) => snippetsService.toggleSave(snippetId, action),
    onMutate: async ({ snippetId, action }) => {
      await queryClient.cancelQueries({ queryKey: queryKeys.detail(snippetId) })

      const previousSnippet = queryClient.getQueryData<Snippet>(queryKeys.detail(snippetId))

      queryClient.setQueryData(queryKeys.detail(snippetId), (old: Snippet | undefined) => {
        if (!old) return old
        return {
          ...old,
          isSaved: action === 'save',
        }
      })

      return { previousSnippet }
    },
    onError: (error, { snippetId, currentIsSaved }, context) => {
      if (context?.previousSnippet) {
        queryClient.setQueryData(queryKeys.detail(snippetId), {
          ...context.previousSnippet,
          isSaved: currentIsSaved,
        })
      }
      toast.error(error.message || 'Failed to update save')
    },
    onSuccess: (data, { snippetId }) => {
      queryClient.setQueryData(queryKeys.detail(snippetId), data)
      queryClient.invalidateQueries({ queryKey: queryKeys.lists() })
      queryClient.invalidateQueries({ queryKey: queryKeys.my() })
      queryClient.invalidateQueries({ queryKey: queryKeys.saved() })
    },
  })
}
