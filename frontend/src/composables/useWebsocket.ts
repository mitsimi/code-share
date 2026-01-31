import { useWebSocketStore } from '@/stores/websocket'
import { useSnippetsStore } from '@/stores/snippets'
import {
  MessageType,
  type WebSocketMessage,
  type UserActionData,
  type ListUpdateData,
  type SnippetUpdateData,
} from '@/services/websocket'
import { queryKeys } from '@/composables/queryKeys'
import { useQueryClient } from '@tanstack/vue-query'
import { onUnmounted } from 'vue'

export function useWebSocket() {
  const wsStore = useWebSocketStore()
  const snippetsStore = useSnippetsStore()
  const queryClient = useQueryClient()

  const mySnippetsQueryKey = queryKeys.mySnippets()
  const likedSnippetsQueryKey = queryKeys.likedSnippets()
  const savedSnippetsQueryKey = queryKeys.savedSnippets()

  const cleanupUserActionsHandler = wsStore.onMessage(
    MessageType.USER_ACTIONS,
    (message: WebSocketMessage) => {
      const actionData = message.data as UserActionData
      snippetsStore.handleUserAction(actionData)

      queryClient.invalidateQueries({ queryKey: mySnippetsQueryKey })
      queryClient.invalidateQueries({ queryKey: likedSnippetsQueryKey })
      queryClient.invalidateQueries({ queryKey: savedSnippetsQueryKey })
    },
  )

  const cleanupListUpdatesHandler = wsStore.onMessage(
    MessageType.LIST_UPDATES,
    (message: WebSocketMessage) => {
      const updateData = message.data as ListUpdateData
      snippetsStore.handleContentUpdate(updateData.snippet_id, updateData)

      queryClient.invalidateQueries({ queryKey: mySnippetsQueryKey })
      queryClient.invalidateQueries({ queryKey: likedSnippetsQueryKey })
      queryClient.invalidateQueries({ queryKey: savedSnippetsQueryKey })
    },
  )

  const cleanupSnippetUpdatesHandler = wsStore.onMessage(
    MessageType.SNIPPET_UPDATES,
    (message: WebSocketMessage) => {
      const snippetData = message.data as SnippetUpdateData

      if (snippetData.update_type === 'stats') {
        snippetsStore.handleStatsUpdate(snippetData.snippet_id, {
          views: snippetData.view_count,
          likes: snippetData.like_count,
        })
      } else if (snippetData.update_type === 'content') {
        snippetsStore.handleContentUpdate(snippetData.snippet_id, snippetData)
      } else if (snippetData.update_type === 'both') {
        snippetsStore.handleContentUpdate(snippetData.snippet_id, snippetData)
        snippetsStore.handleStatsUpdate(snippetData.snippet_id, {
          views: snippetData.view_count,
          likes: snippetData.like_count,
        })
      }

      queryClient.invalidateQueries({ queryKey: mySnippetsQueryKey })
      queryClient.invalidateQueries({ queryKey: likedSnippetsQueryKey })
      queryClient.invalidateQueries({ queryKey: savedSnippetsQueryKey })
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
