import { useWebSocketStore } from '@/stores/websocket'
import { type WebSocketMessage } from '@/services/websocket'
import { useQueryClient } from '@tanstack/vue-query'
import { onUnmounted } from 'vue'

export function useWebSocket() {
  const wsStore = useWebSocketStore()
  const queryClient = useQueryClient()

  // Setup message handlers for TanStack Query integration
  const unsubscribeLike = wsStore.onMessage('user_post_like', (message: WebSocketMessage) => {
    // Update user's like status in cache
    console.log('user_post_like', message)
  })

  const unsubscribeStats = wsStore.onMessage('post_stats', (message: WebSocketMessage) => {
    console.log('post_stats', message)
  })

  const unsubscribeUpdate = wsStore.onMessage('post_update', (message: WebSocketMessage) => {
    // Invalidate and refetch post data
    console.log('post_update', message)
    queryClient.setQueryData(['snippet', message.data.snippet_id], (old: any) => {
      if (!old) return old
      return {
        ...old,
        title: message.data.title,
        content: message.data.content,
        language: message.data.language,
      }
    })
  })

  // Cleanup handlers when component unmounts
  onUnmounted(() => {
    unsubscribeLike()
    unsubscribeStats()
    unsubscribeUpdate()
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
