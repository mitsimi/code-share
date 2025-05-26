import { createFetch } from '@vueuse/core'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

export const useCustomFetch = createFetch({
  fetchOptions: {
    mode: 'cors',
  },
  baseUrl: '/api',
  options: {
    async beforeFetch({ options }) {
      const authStore = useAuthStore()
      if (authStore.token) {
        options.headers = {
          ...options.headers,
          Authorization: `Bearer ${authStore.token}`,
        }
      }
      return { options }
    },
    async afterFetch({ response }) {
      if (response.status === 401) {
        const authStore = useAuthStore()
        const router = useRouter()

        // Clear auth state
        authStore.clearAuth()

        // Redirect to login with current route
        router.push({
          name: 'login',
          query: { redirect: router.currentRoute.value.fullPath },
        })
      }
      return { response }
    },
  },
})
