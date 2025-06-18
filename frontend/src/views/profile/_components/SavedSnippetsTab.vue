<template>
  <SnippetsList
    :snippets="snippetsToShow"
    :is-loading="isLoadingSaved"
    empty-message="You haven't saved any snippets yet."
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { useFetch } from '@/composables/useCustomFetch'
import { useAuthStore } from '@/stores/auth'
import type { Snippet } from '@/types'
import SnippetsList from './SnippetsList.vue'

const authStore = useAuthStore()

// Fetch saved snippets
const { data: savedSnippets, isLoading: isLoadingSaved } = useQuery({
  queryKey: ['saved-snippets'],
  queryFn: async () => {
    const { data, error } = await useFetch<Snippet[]>(`/users/me/saved`).json()
    if (error.value) throw new Error('Failed to fetch saved snippets')
    return data.value.data || []
  },
})

// Provide default empty array when data is undefined
const snippetsToShow = computed(() => savedSnippets.value || [])
</script>
