<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { ref } from 'vue'
import SnippetFormModal from './_components/SnippetFormModal.vue'
import { useSnippets } from '@/composables/useSnippets'

const authStore = useAuthStore()
const showModal = ref(false)

// Use the new hybrid composable
const { snippets, isLoading, isError, error, createSnippet, isCreating, refetch } = useSnippets()

const handleCreateSnippet = (formData: { title: string; content: string; language?: string }) => {
  createSnippet({
    title: formData.title,
    content: formData.content,
    language: formData.language,
  })
  showModal.value = false
}
</script>

<template>
  <main class="container mx-auto max-w-7xl px-8">
    <SnippetGrid
      :snippets="snippets || []"
      :is-loading="isLoading"
      :is-empty="!isLoading && (!snippets || snippets.length === 0)"
      :is-error="isError"
      :error-message="
        (error instanceof Error ? error.message : error) || 'An unexpected error occurred'
      "
      @retry="refetch"
      @create-snippet="showModal = true"
    />
  </main>

  <Authenticated>
    <FloatingActionButton @click="showModal = true" />
  </Authenticated>

  <SnippetFormModal
    :show="showModal"
    :is-loading="isCreating"
    @close="showModal = false"
    @submit="handleCreateSnippet"
  />
</template>
