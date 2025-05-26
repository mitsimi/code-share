import { createFetch } from '@vueuse/core'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

export const useCustomFetch = createFetch({
  fetchOptions: {
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json',
    },
  },
  baseUrl: '/api',
  options: {
    async beforeFetch({ options }) {
      // Ensure headers object exists
      options.headers = options.headers || {}

      const authStore = useAuthStore()
      if (authStore.token) {
        options.headers = {
          ...options.headers,
          Authorization: `Bearer ${authStore.token}`,
        }
      }

      return { options }
    },
    async onFetchError({ data, error }) {
      return { data, error }
    },
    async afterFetch(ctx) {
      const { data, response } = ctx

      if (response.status === 401) {
        const authStore = useAuthStore()
        const router = useRouter()

        try {
          // Try to refresh the token
          await authStore.refreshAccessToken()

          // Retry the original request with new token
          const retryResponse = await fetch(response.url, {
            headers: {
              'Content-Type': 'application/json',
              Authorization: `Bearer ${authStore.token}`,
            },
          })

          const retryData = await retryResponse.json()
          return { data: retryData, response: retryResponse }
        } catch (refreshError) {
          // If refresh fails, clear auth and redirect to login
          authStore.clearAuth()
          router.push({
            name: 'login',
            query: { redirect: router.currentRoute.value.fullPath },
          })
          throw new Error('Authentication failed')
        }
      }

      return ctx
    },
  },
})
