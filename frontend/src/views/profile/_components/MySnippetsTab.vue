<template>
  <SnippetGrid
    :snippets="snippetsToShow"
    :is-loading="isLoadingMySnippets"
    :is-empty="!isLoadingMySnippets && snippetsToShow.length === 0"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { usersService } from '@/services/users'
import SnippetGrid from '@/components/snippets/SnippetGrid.vue'

// Fetch user's snippets
const { data: mySnippets, isLoading: isLoadingMySnippets } = useQuery({
  queryKey: ['my-snippets'],
  queryFn: usersService.getMySnippets,
})

// Provide default empty array when data is undefined
const snippetsToShow = computed(() => mySnippets.value || [])
</script>
