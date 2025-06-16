import { createFetch } from '@vueuse/core'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

// Create a base fetch instance without auth handling
const baseFetch = createFetch({
  fetchOptions: {
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json',
    },
  },
  baseUrl: import.meta.env.VITE_API_URL || '/api',
})

// Create a composable that adds auth handling
export function useCustomFetch() {
  const authStore = useAuthStore()
  const router = useRouter()

  return createFetch({
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
}

// Export a default instance for backward compatibility
export const useFetch = baseFetch
