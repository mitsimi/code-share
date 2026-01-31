import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { snippetsService } from '@/services/snippets'
import { queryKeys } from '@/composables/queryKeys'
import { toast } from 'vue-sonner'
import type { Snippet } from '@/types'

interface LikeVariables {
  snippetId: string
  action: 'like' | 'unlike'
  currentLikes: number
  currentIsLiked: boolean
}

interface LikeContext {
  previousSnippet: Snippet | undefined
}

export function useLikeSnippet() {
  const queryClient = useQueryClient()

  return useMutation<Snippet, Error, LikeVariables, LikeContext>({
    mutationKey: ['like'],
    mutationFn: ({ snippetId, action }) => snippetsService.toggleLike(snippetId, action),
    onMutate: async ({ snippetId, action }) => {
      await queryClient.cancelQueries({ queryKey: queryKeys.detail(snippetId) })

      const previousSnippet = queryClient.getQueryData<Snippet>(queryKeys.detail(snippetId))

      queryClient.setQueryData(queryKeys.detail(snippetId), (old: Snippet | undefined) => {
        if (!old) return old
        const currentCount = old.likes
        return {
          ...old,
          isLiked: action === 'like',
          likes: action === 'like' ? currentCount + 1 : Math.max(0, currentCount - 1),
        }
      })

      return { previousSnippet }
    },
    onError: (error, { snippetId }, context) => {
      if (context?.previousSnippet) {
        queryClient.setQueryData(queryKeys.detail(snippetId), context.previousSnippet)
      }
      toast.error(error.message || 'Failed to update like')
    },
    onSuccess: (data, { snippetId }) => {
      queryClient.setQueryData(queryKeys.detail(snippetId), data)
      queryClient.invalidateQueries({ queryKey: queryKeys.lists() })
      queryClient.invalidateQueries({ queryKey: queryKeys.my() })
      queryClient.invalidateQueries({ queryKey: queryKeys.liked() })
    },
  })
}
