import { useWebSocketStore } from '@/stores/websocket'

/**
 * Thin wrapper around WebSocket store for easy access to connection state and methods.
 * Message handlers and subscriptions are managed in useSnippetSubscription.ts
 */
export function useWebSocket() {
  const wsStore = useWebSocketStore()

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
