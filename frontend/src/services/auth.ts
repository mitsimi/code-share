import { useCustomFetch } from '@/composables/useCustomFetch'
import type { AuthResponse, LoginRequest, RegisterRequest } from '@/types'

export const authService = {
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const useFetch = useCustomFetch()
    const { data, error } = await useFetch<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    }).json()

    // Handle fetch-level errors first
    if (error.value) {
      const errorMessage = error.value.message || 'Failed to login'
      throw new Error(errorMessage)
    }

    if (!data.value) {
      throw new Error('Invalid response from server')
    }

    if (!data.value.data) {
      throw new Error('No authentication data received')
    }

    return data.value.data
  },

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    const useFetch = useCustomFetch()
    const { data, error } = await useFetch<AuthResponse>('/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    }).json()

    if (error.value) {
      const errorMessage = error.value.message || 'Failed to register'
      throw new Error(errorMessage)
    }

    if (!data.value) {
      throw new Error('Invalid response from server')
    }

    if (!data.value.data) {
      throw new Error('No authentication data received')
    }

    return data.value.data
  },

  async logout(): Promise<void> {
    const useFetch = useCustomFetch()
    const { error } = await useFetch<void>('/auth/logout', {
      method: 'POST',
    }).json()

    if (error.value) {
      console.error('Logout error:', error.value)
      const errorMessage = error.value.error || error.value.message || 'Failed to logout'
      throw new Error(errorMessage)
    }
  },

  async refreshToken(refreshToken: string): Promise<AuthResponse> {
    const useFetch = useCustomFetch()
    const { data, error } = await useFetch<AuthResponse>('/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refreshToken: refreshToken }),
    }).json()

    if (error.value) {
      const errorMessage = error.value.message || 'Failed to refresh token'
      throw new Error(errorMessage)
    }

    if (!data.value) {
      throw new Error('Invalid response from server')
    }

    if (!data.value.data) {
      throw new Error('No authentication data received')
    }

    return data.value.data
  },
}
