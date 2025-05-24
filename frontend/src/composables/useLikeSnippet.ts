import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useCustomFetch } from './useCustomFetch'
import { toast } from 'vue-sonner'
import type { Card } from '@/models'

export function useLikeSnippet() {
  const queryClient = useQueryClient()

  const { mutate: updateLike } = useMutation({
    mutationKey: ['likeMutation'],
    mutationFn: async ({ snippetId, action }: { snippetId: string; action: 'like' | 'unlike' }) => {
      console.log(`Starting ${action} mutation for snippet:`, snippetId)
      try {
        const { data, error } = await useCustomFetch<Card>(
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
      }
    },
    onSuccess: (updatedSnippet) => {
      // Update the snippet in the details view
      queryClient.setQueryData(['snippet', updatedSnippet.id], updatedSnippet)

      // Update the snippet in the list view
      queryClient.setQueryData(['snippets'], (oldData: Card[] | undefined) => {
        if (!oldData) return [updatedSnippet]
        return oldData.map((snippet) =>
          snippet.id === updatedSnippet.id
            ? { ...updatedSnippet, isLiked: !snippet.isLiked }
            : snippet,
        )
      })

      /*toast.success(updatedSnippet.isLiked ? 'Added to favorites' : 'Removed from favorites', {
        description: `"${updatedSnippet.title}" ${updatedSnippet.isLiked ? 'added to' : 'removed from'} your favorites`,
      })*/
    },
    onError: (error) => {
      console.error('Like mutation failed:', error)
      toast.error(error.message || 'Please try again')
    },
  })

  return {
    updateLike,
  }
}
