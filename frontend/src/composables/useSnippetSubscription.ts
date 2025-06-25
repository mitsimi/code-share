import { onMounted, onUnmounted, watch } from 'vue'
import { useWebSocket } from './useWebsocket'
import { useAuthStore } from '@/stores/auth'

export function useSnippetList() {
  const { subscribe, unsubscribe } = useWebSocket()

  onMounted(() => {
    // Subscribe to general post updates
    subscribe({ type: 'post_updates' })
  })

  onUnmounted(() => {
    unsubscribe({ type: 'post_updates' })
  })
}

export function useSnippetDetails(snippetId: string) {
  const { subscribe, unsubscribe } = useWebSocket()

  onMounted(() => {
    // Subscribe to specific post stats
    subscribe({
      type: 'post_stats',
      post_id: snippetId,
    })

    // Subscribe to post content updates
    subscribe({
      type: 'post_updates',
      post_id: snippetId,
    })
  })

  onUnmounted(() => {
    unsubscribe({ type: 'post_stats', post_id: snippetId })
    unsubscribe({ type: 'post_updates', post_id: snippetId })
  })
}

export function useUserActions() {
  const { subscribe, unsubscribe } = useWebSocket()
  const { isAuthenticated } = useAuthStore()

  watch(
    isAuthenticated,
    (authStatus) => {
      if (authStatus) {
        subscribe({ type: 'user_actions' })
      } else {
        unsubscribe({ type: 'user_actions' })
      }
    },
    { immediate: true },
  )
}
