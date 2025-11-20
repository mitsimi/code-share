import './assets/main.css'
import './assets/prism.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { useAuthStore } from '@/stores/auth'
import { useWebSocketStore } from '@/stores/websocket'

const app = createApp(App)
const pinia = createPinia()

app.use(router).use(pinia).use(VueQueryPlugin)

// Initialize app
const init = async () => {
  // Initialize Pinia stores
  const authStore = useAuthStore()
  const wsStore = useWebSocketStore()

  try {
    // Try to restore authentication state
    await authStore.initializeAuth()
  } catch (error) {
    console.error('Failed to initialize auth:', error)
  }

  // Initialize WebSocket connection (no auth required)
  wsStore.connect()

  // Mount the app
  app.mount('#app')
}

init()
