import { createFetch } from '@vueuse/core'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

export const useFetch = createFetch({
  fetchOptions: {
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json',
    },
  },
  baseUrl: import.meta.env.VITE_API_URL || '/api',
  options: {
    async beforeFetch({ options }) {
      // Ensure headers object exists
      options.headers = options.headers || {}

      const authStore = useAuthStore()

      // We include the JWT token with every request
      if (authStore.token) {
        options.headers = {
          ...options.headers,
          Authorization: `Bearer ${authStore.token}`,
        }
      }

      return { options }
    },
    async onFetchError({ data, error, response }) {
      // Try to parse as JSON first, if that fails, use the text response
      let errorMessage = error
      if (response) {
        try {
          const jsonError = await response.text()
          errorMessage = jsonError
        } catch {
          // If parsing fails, use the original error
          errorMessage = error
        }
      }
      return { data, error: errorMessage }
    },
    async afterFetch(ctx) {
      const { data, response } = ctx

      if (response.status === 401) {
        const authStore = useAuthStore()
        const router = useRouter()

        // Clear auth and redirect to login
        authStore.clearAuth()
        router.push({
          name: 'login',
          query: { redirect: router.currentRoute.value.fullPath },
        })
        throw new Error('Authentication failed')
      }

      return ctx
    },
  },
})
