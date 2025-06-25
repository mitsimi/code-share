import { onMounted, onUnmounted, watch } from 'vue'
import { useWebSocket } from './useWebsocket'
import { useAuthStore } from '@/stores/auth'
import { SubscriptionType } from '@/services/websocket'

export function useSnippetList() {
  const { subscribe, unsubscribe } = useWebSocket()

  onMounted(() => {
    // Subscribe to general post updates
    subscribe({ type: SubscriptionType.SNIPPET_UPDATES })
  })

  onUnmounted(() => {
    unsubscribe({ type: SubscriptionType.SNIPPET_UPDATES })
  })
}

export function useSnippetDetails(snippetId: string) {
  const { subscribe, unsubscribe } = useWebSocket()

  onMounted(() => {
    // Subscribe to specific post stats
    subscribe({
      type: SubscriptionType.SNIPPET_STATS,
      post_id: snippetId,
    })

    // Subscribe to post content updates
    subscribe({
      type: SubscriptionType.SNIPPET_UPDATES,
      post_id: snippetId,
    })
  })

  onUnmounted(() => {
    unsubscribe({ type: SubscriptionType.SNIPPET_STATS, post_id: snippetId })
    unsubscribe({ type: SubscriptionType.SNIPPET_UPDATES, post_id: snippetId })
  })
}

export function useUserActions() {
  const { subscribe, unsubscribe } = useWebSocket()
  const { isAuthenticated } = useAuthStore()

  watch(
    isAuthenticated,
    (authStatus) => {
      if (authStatus) {
        subscribe({ type: SubscriptionType.USER_ACTIONS })
      } else {
        unsubscribe({ type: SubscriptionType.USER_ACTIONS })
      }
    },
    { immediate: true },
  )
}
