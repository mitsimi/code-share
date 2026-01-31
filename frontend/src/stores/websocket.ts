import { defineStore } from 'pinia'
import { ref, computed, readonly } from 'vue'
import { useAuthStore } from './auth'
import {
  wsService,
  type Subscription,
  type WebSocketMessage,
  type ConnectionState,
  SubscriptionType,
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
    if (
      AUTHENTICATED_SUBSCRIPTION_TYPES.includes(
        subscription.type as (typeof AUTHENTICATED_SUBSCRIPTION_TYPES)[number],
      )
    ) {
      if (!authStore.isAuthenticated()) {
        console.warn(`Cannot subscribe to ${subscription.type}: authentication required`)
        return
      }
    }

    const subKey = getSubscriptionKey(subscription)

    if (subscriptions.value.has(subKey)) {
      return
    }

    // Use the service's subscribe method instead of raw send
    wsService.subscribe(subscription)
    subscriptions.value.add(subKey)
  }

  const unsubscribe = (subscription: Subscription) => {
    const subKey = getSubscriptionKey(subscription)
    subscriptions.value.delete(subKey)

    // Use the service's unsubscribe method instead of raw send
    wsService.unsubscribe(subscription)
  }

  // Clean up authenticated subscriptions when user logs out
  const cleanupAuthenticatedSubscriptions = () => {
    const subsToRemove: string[] = []

    subscriptions.value.forEach((subKey) => {
      const [type] = subKey.split(':')
      if (
        AUTHENTICATED_SUBSCRIPTION_TYPES.includes(
          type as (typeof AUTHENTICATED_SUBSCRIPTION_TYPES)[number],
        )
      ) {
        subsToRemove.push(subKey)
      }
    })

    subsToRemove.forEach((subKey) => {
      subscriptions.value.delete(subKey)
      const [type, snippetId] = subKey.split(':')

      // Send unsubscribe message to server using the service method
      const subscription: Subscription = {
        type: type as SubscriptionType,
        snippet_id: snippetId === 'global' ? undefined : snippetId,
      }
      wsService.unsubscribe(subscription)
    })
  }

  // Message handling (delegates to service)
  const onMessage = (type: string, handler: (message: WebSocketMessage) => void) => {
    console.log('Registering WebSocket message handler for type:', type)
    return wsService.onMessage(type, handler)
  }

  // Utility functions
  const getSubscriptionKey = (sub: Subscription): string => {
    switch (sub.type) {
      case SubscriptionType.SNIPPET_UPDATES:
        return `${sub.type}:${sub.snippet_id}`
      case SubscriptionType.LIST_UPDATES:
        return `${sub.type}`
      case SubscriptionType.USER_ACTIONS:
        return `${sub.type}`
      default:
        throw new Error(`Unknown subscription type: ${sub.type}`)
    }
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
