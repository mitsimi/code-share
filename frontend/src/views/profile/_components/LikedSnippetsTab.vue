<template>
  <SnippetsList
    :snippets="snippetsToShow"
    :is-loading="isLoadingLiked"
    empty-message="You haven't liked any snippets yet."
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { usersService } from '@/services/users'
import { useAuthStore } from '@/stores/auth'
import SnippetsList from './SnippetsList.vue'

const authStore = useAuthStore()

// Fetch liked snippets
const { data: likedSnippets, isLoading: isLoadingLiked } = useQuery({
  queryKey: ['liked-snippets'],
  queryFn: usersService.getLikedSnippets,
})

// Provide default empty array when data is undefined
const snippetsToShow = computed(() => likedSnippets.value || [])
</script>
