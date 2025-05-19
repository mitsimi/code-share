import { useMutation, useQueryClient } from '@tanstack/vue-query'
import { useCustomFetch } from './useCustomFetch'
import { useToast } from './useToast'
import type { Card } from '@/models'

export function useLikeSnippet() {
  const queryClient = useQueryClient()
  const { showToast } = useToast()

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
          snippet.id === updatedSnippet.id ? updatedSnippet : snippet,
        )
      })

      showToast(
        updatedSnippet.likes > 0 ? 'Added to favorites' : 'Removed from favorites',
        updatedSnippet.likes > 0 ? 'success' : 'info',
      )
    },
    onError: (error) => {
      console.error('Like mutation failed:', error)
      showToast(`${error.message || 'Please try again'}`, 'error')
    },
  })

  return {
    updateLike,
  }
}
