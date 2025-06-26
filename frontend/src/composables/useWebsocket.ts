import { useWebSocketStore } from '@/stores/websocket'
import { useSnippetsStore } from '@/stores/snippets'
import {
  MessageType,
  type WebSocketMessage,
  type UserActionData,
  type ListUpdateData,
  type SnippetUpdateData,
} from '@/services/websocket'
import { useQueryClient } from '@tanstack/vue-query'
import { onUnmounted } from 'vue'

export function useWebSocket() {
  const wsStore = useWebSocketStore()
  const snippetsStore = useSnippetsStore()
  const queryClient = useQueryClient()

  // Setup message handlers for TanStack Query integration

  // Update user's like/save status in cache based on action
  // This will sync the UI when actions happen in other tabs/devices
  const cleanupUserActionsHandler = wsStore.onMessage(
    MessageType.USER_ACTIONS,
    (message: WebSocketMessage) => {
      // Handle user action data with proper typing
      const actionData = message.data as UserActionData
      snippetsStore.handleUserAction(actionData)

      // Invalidate related queries to ensure consistency
      queryClient.invalidateQueries({ queryKey: ['my-snippets'] })
      queryClient.invalidateQueries({ queryKey: ['liked-snippets'] })
      queryClient.invalidateQueries({ queryKey: ['saved-snippets'] })
    },
  )

  const cleanupListUpdatesHandler = wsStore.onMessage(
    MessageType.LIST_UPDATES,
    (message: WebSocketMessage) => {
      const updateData = message.data as ListUpdateData

      snippetsStore.handleContentUpdate(updateData.snippet_id, updateData)

      // Invalidate related queries to ensure consistency
      queryClient.invalidateQueries({ queryKey: ['my-snippets'] })
      queryClient.invalidateQueries({ queryKey: ['liked-snippets'] })
      queryClient.invalidateQueries({ queryKey: ['saved-snippets'] })
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

      // Invalidate related queries to ensure consistency
      queryClient.invalidateQueries({ queryKey: ['my-snippets'] })
      queryClient.invalidateQueries({ queryKey: ['liked-snippets'] })
      queryClient.invalidateQueries({ queryKey: ['saved-snippets'] })
    },
  )

  // Cleanup handlers when component unmounts
  onUnmounted(() => {
    cleanupUserActionsHandler()
    cleanupListUpdatesHandler()
    cleanupSnippetUpdatesHandler()
  })

  return {
    // State
    isConnected: wsStore.isConnected,
    isConnecting: wsStore.isConnecting,
    connectionState: wsStore.connectionState,
    subscriptions: wsStore.subscriptions,

    // Methods
    connect: wsStore.connect,
    subscribe: wsStore.subscribe,
    unsubscribe: wsStore.unsubscribe,
    onMessage: wsStore.onMessage,
    cleanupAuthenticatedSubscriptions: wsStore.cleanupAuthenticatedSubscriptions,
  }
}
