import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { useAuthStore } from './auth'
import {
  wsService,
  type Subscription,
  type WebSocketMessage,
  type ConnectionState,
} from '@/services/websocket'

// Define which subscription types require authentication
const AUTHENTICATED_SUBSCRIPTION_TYPES = ['user_actions'] as const

export const useWebSocketStore = defineStore('websocket', () => {
  const authStore = useAuthStore()

  // State
  const connectionState = ref<ConnectionState>('disconnected')
  const subscriptions = ref(new Set<string>())

  // Computed
  const isConnected = computed(() => connectionState.value === 'connected')
  const isConnecting = computed(() => connectionState.value === 'connecting')

  // Initialize service connection state synchronization
  let connectionStateUnsubscribe: (() => void) | null = null

  const initializeService = () => {
    // Sync connection state with service
    connectionStateUnsubscribe = wsService.onConnectionStateChange((state) => {
      connectionState.value = state
    })

    // Set initial state
    connectionState.value = wsService.connectionState
  }

  // Connection management (delegates to service)
  const connect = () => {
    if (!connectionStateUnsubscribe) {
      initializeService()
    }
    wsService.connect()
  }

  const disconnect = () => {
    wsService.disconnect()
    subscriptions.value.clear()

    if (connectionStateUnsubscribe) {
      connectionStateUnsubscribe()
      connectionStateUnsubscribe = null
    }
  }

  // Subscription management
  const subscribe = (subscription: Subscription) => {
    // Check if this subscription type requires authentication
    if (AUTHENTICATED_SUBSCRIPTION_TYPES.includes(subscription.type as any)) {
      if (!authStore.isAuthenticated()) {
        console.warn(`Cannot subscribe to ${subscription.type}: authentication required`)
        return
      }
    }

    const subKey = getSubscriptionKey(subscription)

    if (subscriptions.value.has(subKey)) {
      return
    }

    subscriptions.value.add(subKey)

    wsService.send({
      type: 'subscribe',
      data: subscription,
    })
  }

  const unsubscribe = (subscription: Subscription) => {
    const subKey = getSubscriptionKey(subscription)
    subscriptions.value.delete(subKey)

    wsService.send({
      type: 'unsubscribe',
      data: subscription,
    })
  }

  // Clean up authenticated subscriptions when user logs out
  const cleanupAuthenticatedSubscriptions = () => {
    const subsToRemove: string[] = []

    subscriptions.value.forEach((subKey) => {
      const [type] = subKey.split(':')
      if (AUTHENTICATED_SUBSCRIPTION_TYPES.includes(type as any)) {
        subsToRemove.push(subKey)
      }
    })

    subsToRemove.forEach((subKey) => {
      subscriptions.value.delete(subKey)
      const [type, postId] = subKey.split(':')

      // Send unsubscribe message to server
      wsService.send({
        type: 'unsubscribe',
        data: {
          type: type as any,
          post_id: postId === 'global' ? undefined : postId,
        },
      })
    })
  }

  // Message handling (delegates to service)
  const onMessage = (type: string, handler: (message: WebSocketMessage) => void) => {
    return wsService.onMessage(type, handler)
  }

  // Utility functions
  const getSubscriptionKey = (sub: Subscription): string => {
    return `${sub.type}:${sub.post_id || 'global'}`
  }

  return {
    // State
    connectionState: readonly(connectionState),
    isConnected,
    isConnecting,
    subscriptions: readonly(subscriptions),

    // Methods
    connect,
    disconnect,
    subscribe,
    unsubscribe,
    onMessage,
    cleanupAuthenticatedSubscriptions,
  }
})
