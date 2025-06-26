import { onMounted, onUnmounted, watch } from 'vue'
import { useWebSocket } from './useWebsocket'
import { useAuthStore } from '@/stores/auth'
import { SubscriptionType } from '@/services/websocket'

export function useSnippetList() {
  const { subscribe, unsubscribe } = useWebSocket()

  onMounted(() => {
    subscribe({ type: SubscriptionType.LIST_UPDATES })
  })

  onUnmounted(() => {
    unsubscribe({ type: SubscriptionType.LIST_UPDATES })
  })
}

export function useSnippetDetails(snippetId: string) {
  const { subscribe, unsubscribe } = useWebSocket()

  onMounted(() => {
    subscribe({
      type: SubscriptionType.SNIPPET_UPDATES,
      snippet_id: snippetId,
    })
  })

  onUnmounted(() => {
    unsubscribe({ type: SubscriptionType.SNIPPET_UPDATES, snippet_id: snippetId })
  })
}

export function useUserActions() {
  const { subscribe, unsubscribe } = useWebSocket()
  const authStore = useAuthStore()

  watch(
    () => authStore.isAuthenticated(),
    (isAuthenticated) => {
      if (isAuthenticated) {
        subscribe({ type: SubscriptionType.USER_ACTIONS })
      } else {
        unsubscribe({ type: SubscriptionType.USER_ACTIONS })
      }
    },
    { immediate: true },
  )

  // Cleanup on unmount
  onUnmounted(() => {
    unsubscribe({ type: SubscriptionType.USER_ACTIONS })
  })
}
