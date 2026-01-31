import { useWebSocketStore } from '@/stores/websocket'
import { queryKeys } from '@/composables/queryKeys'
import { useQueryClient } from '@tanstack/vue-query'
import { onUnmounted } from 'vue'
import {
  MessageType,
  type WebSocketMessage,
  type SnippetUpdateData,
  SubscriptionType,
} from '@/services/websocket'

export function useSnippetSubscription(snippetId: string) {
  const wsStore = useWebSocketStore()
  const queryClient = useQueryClient()

  const cleanupSnippetUpdatesHandler = wsStore.onMessage(
    MessageType.SNIPPET_UPDATES,
    (message: WebSocketMessage) => {
      const snippetData = message.data as SnippetUpdateData
      if (snippetData.snippet_id === snippetId) {
        queryClient.invalidateQueries({ queryKey: queryKeys.detail(snippetId) })
        queryClient.invalidateQueries({ queryKey: queryKeys.lists() })
      }
    },
  )
  wsStore.subscribe({ type: SubscriptionType.SNIPPET_UPDATES, snippet_id: snippetId })

  onUnmounted(() => {
    cleanupSnippetUpdatesHandler()
    wsStore.unsubscribe({ type: SubscriptionType.SNIPPET_UPDATES, snippet_id: snippetId })
  })
}

export function useSnippetList() {
  const wsStore = useWebSocketStore()
  const queryClient = useQueryClient()

  const cleanupUserActionsHandler = wsStore.onMessage(MessageType.USER_ACTIONS, () => {
    queryClient.invalidateQueries({ queryKey: queryKeys.lists() })
  })
  wsStore.subscribe({ type: SubscriptionType.USER_ACTIONS })

  const cleanupListUpdatesHandler = wsStore.onMessage(MessageType.LIST_UPDATES, () => {
    queryClient.invalidateQueries({ queryKey: queryKeys.lists() })
  })
  wsStore.subscribe({ type: SubscriptionType.LIST_UPDATES })

  const cleanupSnippetUpdatesHandler = wsStore.onMessage(MessageType.SNIPPET_UPDATES, () => {
    queryClient.invalidateQueries({ queryKey: queryKeys.lists() })
  })
  wsStore.subscribe({ type: SubscriptionType.SNIPPET_UPDATES })

  onUnmounted(() => {
    cleanupUserActionsHandler()
    wsStore.unsubscribe({ type: SubscriptionType.USER_ACTIONS })
    cleanupListUpdatesHandler()
    wsStore.unsubscribe({ type: SubscriptionType.LIST_UPDATES })
    cleanupSnippetUpdatesHandler()
    wsStore.unsubscribe({ type: SubscriptionType.SNIPPET_UPDATES })
  })
}

export const useSnippetDetails = useSnippetSubscription
export const useUserActions = useSnippetList
