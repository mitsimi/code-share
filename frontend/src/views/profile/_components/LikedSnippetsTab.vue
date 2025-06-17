<template>
  <SnippetsList
    :snippets="likedSnippets"
    :is-loading="isLoadingLiked"
    empty-message="You haven't liked any snippets yet."
  />
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { useFetch } from '@/composables/useCustomFetch'
import { useAuthStore } from '@/stores/auth'
import type { Snippet } from '@/types'
import SnippetsList from './SnippetsList.vue'

const authStore = useAuthStore()

// Fetch liked snippets
const { data: likedSnippets, isLoading: isLoadingLiked } = useQuery({
  queryKey: ['liked-snippets'],
  queryFn: async () => {
    const { data, error } = await useFetch<Snippet[]>(`/users/${authStore.user?.id}/liked`).json()
    if (error.value) throw new Error('Failed to fetch liked snippets')
    return data.value || []
  },
})
</script>
