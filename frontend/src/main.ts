import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { VueQueryPlugin } from '@tanstack/vue-query'

const app = createApp(App)
const pinia = createPinia()

app.use(router).use(pinia).use(VueQueryPlugin)

// Initialize app
const init = async () => {
  // Initialize Pinia stores
  const { useAuthStore } = await import('@/stores/auth')
  const authStore = useAuthStore()
  
  try {
    // Try to restore authentication state
    await authStore.initializeAuth()
  } catch (error) {
    console.error('Failed to initialize auth:', error)
  }

  // Mount the app
  app.mount('#app')
}

init()
