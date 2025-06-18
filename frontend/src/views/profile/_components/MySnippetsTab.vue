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
import { usersService } from '@/services/users'
import { useAuthStore } from '@/stores/auth'
import SnippetsList from './SnippetsList.vue'

const authStore = useAuthStore()

// Fetch user's snippets
const { data: mySnippets, isLoading: isLoadingMySnippets } = useQuery({
  queryKey: ['my-snippets'],
  queryFn: usersService.getMySnippets,
})

// Provide default empty array when data is undefined
const snippetsToShow = computed(() => mySnippets.value || [])
</script>
