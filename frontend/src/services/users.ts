import { useFetch } from '@/composables/useCustomFetch'
import type { Snippet, User } from '@/types'

export const usersService = {
  async getMySnippets(): Promise<Snippet[]> {
    const { data, error } = await useFetch<Snippet[]>('/users/me/snippets').json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to fetch your snippets')
    }

    return data.value?.data || []
  },

  async getLikedSnippets(): Promise<Snippet[]> {
    const { data, error } = await useFetch<Snippet[]>('/users/me/liked').json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to fetch liked snippets')
    }

    return data.value?.data || []
  },

  async getSavedSnippets(): Promise<Snippet[]> {
    const { data, error } = await useFetch<Snippet[]>('/users/me/saved').json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to fetch saved snippets')
    }

    return data.value?.data || []
  },

  async getMe(): Promise<User> {
    const { data, error } = await useFetch<User>('/users/me').json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to fetch user profile')
    }

    if (!data.value?.data) {
      throw new Error('No user data received from server')
    }

    return data.value.data
  },

  async updateProfile(profileData: { username: string; email: string }): Promise<User> {
    const { data, error } = await useFetch<User>('/users/me', {
      method: 'PATCH',
      body: JSON.stringify(profileData),
    }).json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to update profile')
    }

    if (!data.value?.data) {
      throw new Error('No data received from server')
    }

    return data.value.data
  },

  async updatePassword(passwordData: {
    currentPassword: string
    newPassword: string
  }): Promise<void> {
    const { error } = await useFetch<void>('/users/me/password', {
      method: 'PATCH',
      body: JSON.stringify(passwordData),
    }).json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to update password')
    }
  },

  async updateAvatar(avatarUrl: string): Promise<string> {
    const { data, error } = await useFetch<{ avatar: string }>('/users/me/avatar', {
      method: 'PATCH',
      body: JSON.stringify({ avatarUrl }),
    }).json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to update avatar')
    }

    if (!data.value?.data) {
      throw new Error('No data received from server')
    }

    return data.value.data.avatar
  },
}
