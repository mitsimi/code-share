import { createFetch, type UseFetchOptions } from '@vueuse/core'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import type { Ref } from 'vue'
import type { APIResponse } from '@/types'

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

// Type-safe fetch response for successful requests
export interface TypedFetchResponse<T> {
  data: Ref<APIResponse<T> | null>
  error: Ref<APIResponse | null>
  execute: () => Promise<void>
  canAbort: Ref<boolean>
  abort: () => void
  isFetching: Ref<boolean>
  isFinished: Ref<boolean>
}

// Create a composable that adds auth handling and type safety
export function useCustomFetch() {
  const authStore = useAuthStore()
  const router = useRouter()

  const createTypedFetch = createFetch({
    fetchOptions: {
      mode: 'cors',
      headers: {
        'Content-Type': 'application/json',
      },
    },
    baseUrl: import.meta.env.VITE_API_URL || '/api',
    options: {
      immediate: true, // Auto-execute the fetch
      async beforeFetch({ options }) {
        // Ensure headers object exists
        options.headers = options.headers || {}

        // Include the JWT token with every request
        if (authStore.token) {
          options.headers = {
            ...options.headers,
            Authorization: `Bearer ${authStore.token}`,
          }
        }

        return { options }
      },
      async onFetchError({ data, error, response }) {
        // Parse the error response as APIResponse
        let apiError: APIResponse | null = null

        if (response) {
          try {
            const errorText = await response.text()
            apiError = JSON.parse(errorText) as APIResponse
          } catch {
            // If parsing fails, create a basic APIResponse structure
            apiError = {
              statusCode: response.status,
              message: response.statusText || 'An error occurred',
              error: error || 'Unknown error',
            }
          }
        } else {
          // Network error or other fetch error
          apiError = {
            statusCode: 0,
            message: 'Network error',
            error: error || 'Failed to connect to server',
          }
        }

        return { data, error: apiError }
      },
      async afterFetch(ctx) {
        const { response } = ctx

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

  // Return a typed version of the fetch function
  return function typedFetch<T = unknown>(url: string, options?: unknown) {
    const response = createTypedFetch(url, options as UseFetchOptions)

    // Return with proper typing
    return {
      ...response,
      // Override the json method to return properly typed response
      json: () => {
        const jsonResponse = response.json()
        return {
          ...jsonResponse,
          data: jsonResponse.data as Ref<APIResponse<T> | null>,
          error: jsonResponse.error as Ref<APIResponse | null>,
        } as TypedFetchResponse<T>
      },
    }
  }
}

// Export a default instance for backward compatibility
export const useFetch = baseFetch
