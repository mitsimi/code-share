<template>
  <SnippetGrid
    :snippets="snippetsToShow"
    :is-loading="isLoadingLiked"
    :is-empty="!isLoadingLiked && snippetsToShow.length === 0"
    empty-title="No liked snippets yet"
    empty-message="When you like code snippets, they will appear here for easy access."
    :show-create-button="false"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { usersService } from '@/services/users'
import SnippetGrid from '@/components/snippets/SnippetGrid.vue'

// Fetch liked snippets
const { data: likedSnippets, isLoading: isLoadingLiked } = useQuery({
  queryKey: ['liked-snippets'],
  queryFn: usersService.getLikedSnippets,
})

// Provide default empty array when data is undefined
const snippetsToShow = computed(() => likedSnippets.value || [])
</script>
