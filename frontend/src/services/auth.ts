import { useCustomFetch } from '@/composables/useCustomFetch'
import type { AuthResponse, LoginRequest, SignupRequest, User } from '@/types'

export const authService = {
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const { data, error } = await useCustomFetch<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    }).json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to login')
    }

    return data.value!
  },

  async signup(userData: SignupRequest): Promise<AuthResponse> {
    const { data, error } = await useCustomFetch<AuthResponse>('/auth/signup', {
      method: 'POST',
      body: JSON.stringify(userData),
    }).json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to sign up')
    }

    return data.value!
  },

  async logout(): Promise<void> {
    const { error } = await useCustomFetch('/auth/logout', {
      method: 'POST',
    }).json()

    if (error.value) {
      throw new Error('Failed to logout')
    }
  },

  async refreshToken(): Promise<AuthResponse> {
    const { data, error } = await useCustomFetch<AuthResponse>('/auth/refresh', {
      method: 'POST',
    }).json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to refresh token')
    }

    return data.value!
  },

  async getCurrentUser(): Promise<User> {
    const { data, error } = await useCustomFetch<User>('/auth/me').json()

    if (error.value) {
      throw new Error('Not authenticated')
    }

    return data.value!
  },
}
