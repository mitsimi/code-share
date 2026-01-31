import { onMounted, onUnmounted, watch } from 'vue'
import { useWebSocketStore } from '@/stores/websocket'
import { useAuthStore } from '@/stores/auth'
import { queryKeys } from '@/composables/queryKeys'
import { useQueryClient } from '@tanstack/vue-query'
import {
  MessageType,
  SubscriptionType,
  type WebSocketMessage,
  type UserActionData,
  type ListUpdateData,
  type SnippetUpdateData,
} from '@/services/websocket'
import type { Snippet } from '@/types'

/**
 * Helper to update snippet in TanStack Query cache
 */
function useSnippetCacheUpdater() {
  const queryClient = useQueryClient()

  return (snippetId: string, updates: Partial<Snippet>) => {
    const detailKey = queryKeys.detail(snippetId)
    const current = queryClient.getQueryData<Snippet>(detailKey)
    if (current) {
      queryClient.setQueryData(detailKey, { ...current, ...updates })
    }
  }
}

/**
 * Subscribe to list updates - for all users viewing snippet lists
 * Handles: LIST_UPDATES (content changes broadcasted to all list viewers)
 */
export function useSnippetList() {
  const wsStore = useWebSocketStore()
  const queryClient = useQueryClient()
  const updateSnippetInCache = useSnippetCacheUpdater()

  const cleanupHandler = wsStore.onMessage(
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

  onMounted(() => {
    wsStore.subscribe({ type: SubscriptionType.LIST_UPDATES })
  })

  onUnmounted(() => {
    cleanupHandler()
    wsStore.unsubscribe({ type: SubscriptionType.LIST_UPDATES })
  })
}

/**
 * Subscribe to snippet detail updates - for users viewing a specific snippet
 * Handles: SNIPPET_UPDATES (stats + content for specific snippet)
 */
export function useSnippetDetails(snippetId: string) {
  const wsStore = useWebSocketStore()
  const queryClient = useQueryClient()
  const updateSnippetInCache = useSnippetCacheUpdater()

  const cleanupHandler = wsStore.onMessage(
    MessageType.SNIPPET_UPDATES,
    (message: WebSocketMessage) => {
      const snippetData = message.data as SnippetUpdateData

      // Only process if this update is for our snippet
      if (snippetData.snippet_id !== snippetId) return

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

      updateSnippetInCache(snippetId, updates)
      queryClient.invalidateQueries({ queryKey: queryKeys.detail(snippetId) })
    },
  )

  onMounted(() => {
    wsStore.subscribe({ type: SubscriptionType.SNIPPET_UPDATES, snippet_id: snippetId })
  })

  onUnmounted(() => {
    cleanupHandler()
    wsStore.unsubscribe({ type: SubscriptionType.SNIPPET_UPDATES, snippet_id: snippetId })
  })
}

/**
 * Subscribe to user actions - ONLY for authenticated users
 * Handles: USER_ACTIONS (like/save sync across devices/tabs)
 * Uses watch to auto-subscribe/unsubscribe based on auth state
 */
export function useUserActions() {
  const wsStore = useWebSocketStore()
  const authStore = useAuthStore()
  const queryClient = useQueryClient()
  const updateSnippetInCache = useSnippetCacheUpdater()

  const cleanupHandler = wsStore.onMessage(
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

      // Invalidate lists that depend on user state
      queryClient.invalidateQueries({ queryKey: queryKeys.lists() })
      queryClient.invalidateQueries({ queryKey: queryKeys.my() })
      queryClient.invalidateQueries({ queryKey: queryKeys.liked() })
      queryClient.invalidateQueries({ queryKey: queryKeys.saved() })
    },
  )

  // Watch auth state - subscribe when authenticated, unsubscribe when not
  watch(
    () => authStore.isAuthenticated(),
    (isAuthenticated) => {
      if (isAuthenticated) {
        wsStore.subscribe({ type: SubscriptionType.USER_ACTIONS })
      } else {
        wsStore.unsubscribe({ type: SubscriptionType.USER_ACTIONS })
      }
    },
    { immediate: true },
  )

  onUnmounted(() => {
    cleanupHandler()
    wsStore.unsubscribe({ type: SubscriptionType.USER_ACTIONS })
  })
}
