import { useFetch } from '@/composables/useCustomFetch'
import type { Snippet } from '@/types'

export const snippetsService = {
  async getSnippets(): Promise<Snippet[]> {
    const { data, error } = await useFetch<Snippet[]>('/snippets', {
      timeout: 1000,
    }).json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to fetch snippets')
    }

    return data.value?.data || []
  },

  async getSnippet(snippetId: string): Promise<Snippet> {
    const { data, error } = await useFetch<Snippet>(`/snippets/${snippetId}`, {
      timeout: 1000,
    }).json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to fetch snippet')
    }

    if (!data.value?.data) {
      throw new Error('Snippet not found')
    }

    return data.value.data
  },

  async createSnippet(formData: {
    title: string
    content: string
    language?: string
  }): Promise<Snippet> {
    const { data, error } = await useFetch<Snippet>('/snippets', {
      method: 'POST',
      body: JSON.stringify({
        title: formData.title,
        content: formData.content,
        language: formData.language || 'plaintext',
      }),
    }).json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to create snippet')
    }

    if (!data.value?.data) {
      throw new Error('No data received from server')
    }

    return data.value.data
  },

  async updateSnippet(
    snippetId: string,
    formData: { title?: string; content?: string; language?: string },
  ): Promise<Snippet> {
    const { data, error } = await useFetch<Snippet>(`/snippets/${snippetId}`, {
      method: 'PUT',
      body: JSON.stringify(formData),
    }).json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to update snippet')
    }

    if (!data.value?.data) {
      throw new Error('No data received from server')
    }

    return data.value.data
  },

  async deleteSnippet(snippetId: string): Promise<void> {
    const { error } = await useFetch<void>(`/snippets/${snippetId}`, {
      method: 'DELETE',
    }).json()

    if (error.value) {
      throw new Error(error.value.message || 'Failed to delete snippet')
    }
  },

  async toggleLike(snippetId: string, action: 'like' | 'unlike'): Promise<Snippet> {
    const { data, error } = await useFetch<Snippet>(
      `/snippets/${snippetId}/like?action=${action}`,
      {
        method: 'PATCH',
      },
    ).json()

    if (error.value) {
      throw new Error(error.value.message || `Failed to ${action} snippet`)
    }

    if (!data.value?.data) {
      throw new Error('No data received from server')
    }

    return data.value.data
  },

  async toggleSave(snippetId: string, action: 'save' | 'unsave'): Promise<Snippet> {
    const { data, error } = await useFetch<Snippet>(
      `/snippets/${snippetId}/save?action=${action}`,
      {
        method: 'PATCH',
      },
    ).json()

    if (error.value) {
      throw new Error(error.value.message || `Failed to ${action} snippet`)
    }

    if (!data.value?.data) {
      throw new Error('No data received from server')
    }

    return data.value.data
  },
}
