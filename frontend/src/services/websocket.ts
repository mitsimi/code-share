// composables/useWebSocket.ts

export interface WebSocketMessage {
  type: string
  data: any
  snippet_id?: string
  user_id?: string
  timestamp: number
}

export interface Subscription {
  type: 'user_actions' | 'post_updates' | 'post_stats'
  post_id?: string
}

class WebSocketService {
  private ws: WebSocket | null = null
  private subscriptions = new Set<string>()
  private messageHandlers = new Map<string, Function[]>()

  connect() {
    const wsUrl = `ws://localhost:8080/ws`
    this.ws = new WebSocket(wsUrl)

    this.ws.onopen = () => {
      console.log('WebSocket connected')
    }

    this.ws.onmessage = (event) => {
      const message: WebSocketMessage = JSON.parse(event.data)
      console.log('Received message', message)
      //this.handleMessage(message)
    }

    this.ws.onclose = () => {
      console.log('WebSocket disconnected')
    }
  }

  subscribe(subscription: Subscription) {
    if (!this.ws) {
      this.connect()
    }

    const subKey = this.getSubscriptionKey(subscription)
    if (this.subscriptions.has(subKey)) return

    this.subscriptions.add(subKey)
    this.send({
      type: 'subscribe',
      data: subscription,
    })
  }

  unsubscribe(subscription: Subscription) {
    if (!this.ws) {
      this.connect()
    }

    const subKey = this.getSubscriptionKey(subscription)
    this.subscriptions.delete(subKey)
    this.send({
      type: 'unsubscribe',
      data: subscription,
    })
  }

  private handleMessage(message: WebSocketMessage) {
    const handlers = this.messageHandlers.get(message.type) || []
    handlers.forEach((handler) => handler(message))
  }

  onMessage(type: string, handler: Function) {
    if (!this.messageHandlers.has(type)) {
      this.messageHandlers.set(type, [])
    }
    this.messageHandlers.get(type)!.push(handler)
  }

  private send(data: any) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    }
  }

  private getSubscriptionKey(sub: Subscription): string {
    return `${sub.type}:${sub.post_id || 'global'}`
  }
}

export const wsService = new WebSocketService()
