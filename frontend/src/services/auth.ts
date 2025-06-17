import { useFetch } from '@/composables/useCustomFetch'
import type { AuthResponse, LoginRequest, SignupRequest } from '@/types'

export const authService = {
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const { data, error } = await useFetch<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    }).json()

    if (error.value) {
      const errorMessage =
        typeof error.value === 'string' ? error.value : error.value.message || 'Failed to login'
      throw new Error(errorMessage)
    }

    if (!data.value) {
      throw new Error('Invalid response from server')
    }

    return data.value
  },

  async signup(userData: SignupRequest): Promise<AuthResponse> {
    const { data, error } = await useFetch<AuthResponse>('/auth/signup', {
      method: 'POST',
      body: JSON.stringify(userData),
    }).json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to sign up')
    }

    if (!data.value) {
      throw new Error('No response data received')
    }

    return data.value
  },

  async logout(): Promise<void> {
    const { error } = await useFetch('/auth/logout', {
      method: 'POST',
    }).json()

    if (error.value) {
      throw new Error('Failed to logout')
    }
  },

  async refreshToken(refreshToken: string): Promise<AuthResponse> {
    const { data, error } = await useFetch<AuthResponse>('/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refreshToken: refreshToken }),
    }).json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to refresh token')
    }

    if (!data.value) {
      throw new Error('No response data received')
    }

    return data.value
  },
}
