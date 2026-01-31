<template>
  <SnippetGrid
    :snippets="snippetsToShow"
    :is-loading="isLoadingSaved"
    :is-empty="!isLoadingSaved && snippetsToShow.length === 0"
    empty-title="No saved snippets yet"
    empty-message="Save interesting code snippets to build your personal collection."
    :show-create-button="false"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { usersService } from '@/services/users'
import { queryKeys } from '@/composables/queryKeys'
import SnippetGrid from '@/components/snippets/SnippetGrid.vue'

const { data: savedSnippets, isLoading: isLoadingSaved } = useQuery({
  queryKey: queryKeys.saved(),
  queryFn: usersService.getSavedSnippets,
})

// Provide default empty array when data is undefined
const snippetsToShow = computed(() => savedSnippets.value || [])
</script>
