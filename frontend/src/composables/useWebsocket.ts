import { wsService, type WebSocketMessage } from '@/services/websocket'
import { useQueryClient } from '@tanstack/vue-query'

export function useWebSocket() {
  const queryClient = useQueryClient()

  // Setup message handlers for TanStack Query integration
  wsService.onMessage('user_post_like', (message: WebSocketMessage) => {
    // Update user's like status in cache
    console.log('user_post_like', message)
  })

  wsService.onMessage('post_stats', (message: WebSocketMessage) => {
    console.log('post_stats', message)
  })

  wsService.onMessage('post_update', (message: WebSocketMessage) => {
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

  return {
    connect: wsService.connect.bind(wsService),
    subscribe: wsService.subscribe.bind(wsService),
    unsubscribe: wsService.unsubscribe.bind(wsService),
  }
}
