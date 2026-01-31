import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { VueQueryPlugin, type VueQueryPluginOptions } from '@tanstack/vue-query'
import { useAuthStore } from '@/stores/auth'
import { useWebSocketStore } from '@/stores/websocket'

const app = createApp(App)
const pinia = createPinia()

const vueQueryOptions: VueQueryPluginOptions = {
  queryClientConfig: {
    defaultOptions: {
      queries: {
        staleTime: 5 * 60 * 1000,
        gcTime: 10 * 60 * 1000,
        retry: 3,
        refetchOnWindowFocus: true,
        refetchOnReconnect: true,
      },
      mutations: {
        retry: 1,
      },
    },
  },
}

app.use(router).use(pinia).use(VueQueryPlugin, vueQueryOptions)

const init = async () => {
  const authStore = useAuthStore()
  const wsStore = useWebSocketStore()

  try {
    await authStore.initializeAuth()
  } catch (error) {
    console.error('Failed to initialize auth:', error)
  }

  wsStore.connect()
  app.mount('#app')
}

init()
