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
import { usersService } from '@/services/users'
import { useAuthStore } from '@/stores/auth'
import SnippetsList from './SnippetsList.vue'

const authStore = useAuthStore()

// Fetch saved snippets
const { data: savedSnippets, isLoading: isLoadingSaved } = useQuery({
  queryKey: ['saved-snippets'],
  queryFn: usersService.getSavedSnippets,
})

// Provide default empty array when data is undefined
const snippetsToShow = computed(() => savedSnippets.value || [])
</script>
