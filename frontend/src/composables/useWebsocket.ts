import { useWebSocketStore } from '@/stores/websocket'
import { queryKeys } from '@/composables/queryKeys'
import {
  MessageType,
  type WebSocketMessage,
  type UserActionData,
  type ListUpdateData,
  type SnippetUpdateData,
} from '@/services/websocket'
import { useQueryClient } from '@tanstack/vue-query'
import { onUnmounted } from 'vue'
import type { Snippet } from '@/types'

export function useWebSocket() {
  const wsStore = useWebSocketStore()
  const queryClient = useQueryClient()

  const updateSnippetInCache = (snippetId: string, updates: Partial<Snippet>) => {
    const detailKey = queryKeys.detail(snippetId)
    const current = queryClient.getQueryData<Snippet>(detailKey)

    if (current) {
      queryClient.setQueryData(detailKey, { ...current, ...updates })
    }
  }

  const cleanupUserActionsHandler = wsStore.onMessage(
    MessageType.USER_ACTIONS,
    (message: WebSocketMessage) => {
      const actionData = message.data as UserActionData

      if (actionData.action === 'like' || actionData.action === 'unlike') {
        updateSnippetInCache(actionData.snippet_id, {
          isLiked: actionData.value,
          likes: actionData.like_count,
        })
      } else if (actionData.action === 'save' || actionData.action === 'unsave') {
        updateSnippetInCache(actionData.snippet_id, {
          isSaved: actionData.value,
        })
      }

      queryClient.invalidateQueries({ queryKey: queryKeys.lists() })
      queryClient.invalidateQueries({ queryKey: queryKeys.my() })
      queryClient.invalidateQueries({ queryKey: queryKeys.liked() })
      queryClient.invalidateQueries({ queryKey: queryKeys.saved() })
    },
  )

  const cleanupListUpdatesHandler = wsStore.onMessage(
    MessageType.LIST_UPDATES,
    (message: WebSocketMessage) => {
      const updateData = message.data as ListUpdateData

      updateSnippetInCache(updateData.snippet_id, {
        title: updateData.title,
        content: updateData.content,
        language: updateData.language,
      })

      queryClient.invalidateQueries({ queryKey: queryKeys.lists() })
    },
  )

  const cleanupSnippetUpdatesHandler = wsStore.onMessage(
    MessageType.SNIPPET_UPDATES,
    (message: WebSocketMessage) => {
      const snippetData = message.data as SnippetUpdateData
      const updates: Partial<Snippet> = {}

      if (snippetData.update_type === 'content' || snippetData.update_type === 'both') {
        if (snippetData.title) updates.title = snippetData.title
        if (snippetData.content) updates.content = snippetData.content
        if (snippetData.language) updates.language = snippetData.language
      }

      if (snippetData.update_type === 'stats' || snippetData.update_type === 'both') {
        if (snippetData.view_count !== undefined) updates.views = snippetData.view_count
        if (snippetData.like_count !== undefined) updates.likes = snippetData.like_count
      }

      updateSnippetInCache(snippetData.snippet_id, updates)
      queryClient.invalidateQueries({ queryKey: queryKeys.lists() })
    },
  )

  onUnmounted(() => {
    cleanupUserActionsHandler()
    cleanupListUpdatesHandler()
    cleanupSnippetUpdatesHandler()
  })

  return {
    isConnected: wsStore.isConnected,
    isConnecting: wsStore.isConnecting,
    connectionState: wsStore.connectionState,
    subscriptions: wsStore.subscriptions,
    connect: wsStore.connect,
    subscribe: wsStore.subscribe,
    unsubscribe: wsStore.unsubscribe,
    onMessage: wsStore.onMessage,
    cleanupAuthenticatedSubscriptions: wsStore.cleanupAuthenticatedSubscriptions,
  }
}
