<template>
  <SnippetsList
    :snippets="mySnippets"
    :is-loading="isLoadingMySnippets"
    empty-message="You haven't created any snippets yet."
  />
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { useFetch } from '@/composables/useCustomFetch'
import { useAuthStore } from '@/stores/auth'
import type { Snippet } from '@/types'
import SnippetsList from './SnippetsList.vue'

const authStore = useAuthStore()

// Fetch user's snippets
const { data: mySnippets, isLoading: isLoadingMySnippets } = useQuery({
  queryKey: ['my-snippets'],
  queryFn: async () => {
    const { data, error } = await useFetch<Snippet[]>(
      `/users/${authStore.user?.id}/snippets`,
    ).json()
    if (error.value) throw new Error('Failed to fetch your snippets')
    return data.value || []
  },
})
</script>
