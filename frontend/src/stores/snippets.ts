import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Snippet } from '@/types'
import type { UserActionData } from '@/services/websocket'

export const useSnippetsStore = defineStore('snippets', () => {
  // State - using Map for O(1) lookups
  const snippetsMap = ref<Map<string, Snippet>>(new Map())
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const snippets = computed(() => {
    // Convert Map to Array sorted by creation date (newest first)
    return Array.from(snippetsMap.value.values()).sort(
      (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
    )
  })

  const getSnippetById = computed(() => {
    return (id: string) => snippetsMap.value.get(id)
  })

  const snippetsCount = computed(() => snippetsMap.value.size)

  const hasSnippet = computed(() => {
    return (id: string) => snippetsMap.value.has(id)
  })

  // Actions
  const setSnippets = (newSnippets: Snippet[]) => {
    // Clear existing and populate Map
    snippetsMap.value.clear()
    newSnippets.forEach((snippet) => {
      snippetsMap.value.set(snippet.id, snippet)
    })
    error.value = null
  }

  const addSnippet = (snippet: Snippet) => {
    snippetsMap.value.set(snippet.id, snippet)
  }

  const updateSnippet = (snippetId: string, updates: Partial<Snippet>) => {
    const existing = snippetsMap.value.get(snippetId)
    if (existing) {
      const updated = { ...existing, ...updates }
      snippetsMap.value.set(snippetId, updated)
    }
  }

  const removeSnippet = (snippetId: string) => {
    snippetsMap.value.delete(snippetId)
  }

  // WebSocket-specific actions
  const handleUserAction = (actionData: UserActionData) => {
    const snippet = snippetsMap.value.get(actionData.snippet_id)
    if (!snippet) return

    const updated = { ...snippet }

    switch (actionData.action) {
      case 'like':
        updated.isLiked = actionData.value
        updated.likes = actionData.like_count ?? updated.likes + 1
      case 'unlike':
        updated.isLiked = actionData.value
        updated.likes = actionData.like_count
          ? actionData.like_count
          : updated.likes - 1 <= 0
            ? 0
            : updated.likes - 1
        break
      case 'save':
      case 'unsave':
        updated.isSaved = actionData.value
        break
    }

    snippetsMap.value.set(actionData.snippet_id, updated)
  }

  const handleContentUpdate = (
    snippetId: string,
    updates: { title?: string; content?: string; language?: string },
  ) => {
    updateSnippet(snippetId, updates)
  }

  const handleStatsUpdate = (snippetId: string, stats: { views?: number; likes?: number }) => {
    updateSnippet(snippetId, stats)
  }

  const setLoading = (loading: boolean) => {
    isLoading.value = loading
  }

  const setError = (errorMessage: string | null) => {
    error.value = errorMessage
  }

  const clearError = () => {
    error.value = null
  }

  const clear = () => {
    snippetsMap.value.clear()
    error.value = null
  }

  return {
    // State
    snippetsMap,
    isLoading,
    error,

    // Getters
    snippets,
    getSnippetById,
    snippetsCount,
    hasSnippet,

    // Actions
    setSnippets,
    addSnippet,
    updateSnippet,
    removeSnippet,
    handleUserAction,
    handleContentUpdate,
    handleStatsUpdate,
    setLoading,
    setError,
    clearError,
    clear,
  }
})
