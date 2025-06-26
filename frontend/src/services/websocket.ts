// composables/useWebSocket.ts

export interface WebSocketMessage {
  type: string
  data: any
  snippet_id?: string
  user_id?: string
  timestamp: number
}

// Message type constants to match backend
export const MessageType = {
  ERROR: 'error',
  SUCCESS: 'success',
  SUBSCRIBE: 'subscribe',
  UNSUBSCRIBE: 'unsubscribe',
  USER_ACTIONS: 'user_actions',
  SNIPPET_UPDATES: 'snippet_updates',
  LIST_UPDATES: 'list_updates',
} as const

export enum SubscriptionType {
  USER_ACTIONS = 'user_actions',
  SNIPPET_UPDATES = 'snippet_updates',
  LIST_UPDATES = 'list_updates',
}

export interface Subscription {
  type: SubscriptionType
  snippet_id?: string
}

// Backend-compatible subscription request structure
export interface SubscriptionRequest {
  type: SubscriptionType
  snippet_id?: string
}

// User action data structure from backend
export interface UserActionData {
  action: string // "like", "unlike", "save", "unsave"
  snippet_id: string
  value: boolean // true for like/save, false for unlike/unsave
  like_count?: number
}

// Snippet update data structure from backend
export interface SnippetUpdateData {
  snippet_id: string
  update_type: string // "content", "stats", "both"
  // Content changes (optional)
  title?: string
  content?: string
  language?: string
  // Stats changes (optional)
  view_count?: number
  like_count?: number
}

// List update data structure from backend
export interface ListUpdateData {
  snippet_id: string
  title?: string
  content?: string
  language?: string
}

export type ConnectionState = 'disconnected' | 'connecting' | 'connected' | 'error'

export interface WebSocketConfig {
  url: string
  maxReconnectAttempts?: number
  reconnectDelay?: (attempt: number) => number
}

export type MessageHandler = (message: WebSocketMessage) => void
export type ConnectionStateHandler = (state: ConnectionState) => void

export class WebSocketService {
  private ws: WebSocket | null = null
  private config: Required<WebSocketConfig>
  private messageHandlers = new Map<string, MessageHandler[]>()
  private connectionStateHandlers: ConnectionStateHandler[] = []
  private messageQueue: any[] = []
  private reconnectAttempts = 0
  private reconnectTimer: number | null = null
  private _connectionState: ConnectionState = 'disconnected'

  constructor(config: WebSocketConfig) {
    this.config = {
      url: config.url,
      maxReconnectAttempts: config.maxReconnectAttempts ?? 5,
      reconnectDelay:
        config.reconnectDelay ?? ((attempt) => Math.min(1000 * Math.pow(2, attempt), 30000)),
    }
  }

  get connectionState(): ConnectionState {
    return this._connectionState
  }

  get isConnected(): boolean {
    return this._connectionState === 'connected'
  }

  get isConnecting(): boolean {
    return this._connectionState === 'connecting'
  }

  // Connection management
  connect(): void {
    if (this._connectionState === 'connecting' || this._connectionState === 'connected') {
      return
    }

    this.setConnectionState('connecting')

    try {
      this.ws = new WebSocket(this.config.url)
      this.setupEventHandlers()
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error)
      this.setConnectionState('error')
    }
  }

  disconnect(): void {
    this.clearReconnectTimer()

    if (this.ws) {
      this.ws.close(1000, 'Client disconnecting')
      this.ws = null
    }

    this.setConnectionState('disconnected')
    this.messageQueue = []
    this.reconnectAttempts = 0
  }

  // Message handling
  send(data: any): void {
    if (this.isConnected && this.ws) {
      this.ws.send(JSON.stringify(data))
    } else {
      this.messageQueue.push(data)

      // Auto-connect if not already connecting
      if (!this.isConnecting) {
        this.connect()
      }
    }
  }

  // Send subscription request in backend-compatible format
  subscribe(subscription: Subscription): void {
    const subscriptionRequest: SubscriptionRequest = {
      type: subscription.type,
      snippet_id: subscription.snippet_id,
    }

    this.send({
      type: MessageType.SUBSCRIBE,
      data: subscriptionRequest,
      timestamp: Date.now(),
    })
  }

  // Send unsubscription request in backend-compatible format
  unsubscribe(subscription: Subscription): void {
    const subscriptionRequest: SubscriptionRequest = {
      type: subscription.type,
      snippet_id: subscription.snippet_id,
    }

    this.send({
      type: MessageType.UNSUBSCRIBE,
      data: subscriptionRequest,
      timestamp: Date.now(),
    })
  }

  onMessage(type: string, handler: MessageHandler): () => void {
    if (!this.messageHandlers.has(type)) {
      this.messageHandlers.set(type, [])
    }
    this.messageHandlers.get(type)!.push(handler)

    // Return unsubscribe function
    return () => {
      const handlers = this.messageHandlers.get(type)
      if (handlers) {
        const index = handlers.indexOf(handler)
        if (index > -1) {
          handlers.splice(index, 1)
        }
      }
    }
  }

  onConnectionStateChange(handler: ConnectionStateHandler): () => void {
    this.connectionStateHandlers.push(handler)

    // Return unsubscribe function
    return () => {
      const index = this.connectionStateHandlers.indexOf(handler)
      if (index > -1) {
        this.connectionStateHandlers.splice(index, 1)
      }
    }
  }

  // Private methods
  private setupEventHandlers(): void {
    if (!this.ws) return

    this.ws.onopen = () => {
      console.log('WebSocket connected')
      this.setConnectionState('connected')
      this.reconnectAttempts = 0
      this.processMessageQueue()
    }

    this.ws.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data)
        this.handleMessage(message)
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    this.ws.onclose = (event) => {
      console.log('WebSocket disconnected', event.code, event.reason)
      this.setConnectionState('disconnected')
      this.ws = null

      // Attempt reconnection if it wasn't a clean close
      if (event.code !== 1000 && this.reconnectAttempts < this.config.maxReconnectAttempts) {
        this.scheduleReconnect()
      }
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
      this.setConnectionState('error')
    }
  }

  private handleMessage(message: WebSocketMessage): void {
    // Handle backend response messages
    if (message.type === MessageType.SUCCESS || message.type === MessageType.ERROR) {
      return
    }

    // Forward other messages to registered handlers
    const handlers = this.messageHandlers.get(message.type) || []
    handlers.forEach((handler) => {
      try {
        handler(message)
      } catch (error) {
        console.error('Error in WebSocket message handler:', error)
      }
    })
  }

  private setConnectionState(state: ConnectionState): void {
    if (this._connectionState !== state) {
      this._connectionState = state
      this.connectionStateHandlers.forEach((handler) => {
        try {
          handler(state)
        } catch (error) {
          console.error('Error in connection state handler:', error)
        }
      })
    }
  }

  private processMessageQueue(): void {
    while (this.messageQueue.length > 0 && this.isConnected && this.ws) {
      const message = this.messageQueue.shift()
      this.ws.send(JSON.stringify(message))
    }
  }

  private scheduleReconnect(): void {
    this.clearReconnectTimer()

    const delay = this.config.reconnectDelay(this.reconnectAttempts)
    this.reconnectAttempts++

    console.log(`Scheduling WebSocket reconnect attempt ${this.reconnectAttempts} in ${delay}ms`)

    this.reconnectTimer = window.setTimeout(() => {
      this.connect()
    }, delay)
  }

  private clearReconnectTimer(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
  }
}

// Create singleton instance
const getWebSocketUrl = (): string => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  return `${protocol}//${host}/ws`
}

export const wsService = new WebSocketService({
  url: getWebSocketUrl(),
})
