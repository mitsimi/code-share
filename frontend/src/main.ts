import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { VueQueryPlugin } from '@tanstack/vue-query'
import { useAuthStore } from '@/stores/auth'

const app = createApp(App)
const pinia = createPinia()

app.use(router).use(pinia).use(VueQueryPlugin)

// Initialize app
const init = async () => {
  // Initialize Pinia stores
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
