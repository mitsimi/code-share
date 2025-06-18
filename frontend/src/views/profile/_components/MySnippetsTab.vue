<template>
  <SnippetsList
    :snippets="snippetsToShow"
    :is-loading="isLoadingMySnippets"
    empty-message="You haven't created any snippets yet."
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

// Fetch user's snippets
const { data: mySnippets, isLoading: isLoadingMySnippets } = useQuery({
  queryKey: ['my-snippets'],
  queryFn: async () => {
    const { data, error } = await useFetch<Snippet[]>(`/users/me/snippets`).json()
    if (error.value) throw new Error('Failed to fetch your snippets')
    return data.value.data || []
  },
})

// Provide default empty array when data is undefined
const snippetsToShow = computed(() => mySnippets.value || [])
</script>
