<script setup lang="ts">
import { ref } from 'vue'
import SnippetFormModal from './_components/SnippetFormModal.vue'
import { useSnippets } from '@/composables/useSnippets'

const showModal = ref(false)

const {
  getList: { data: snippets, isLoading, isError, error, refetch },
  createSnippet,
} = useSnippets()

const handleCreateSnippet = (formData: { title: string; content: string; language?: string }) => {
  createSnippet.mutate({
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
    :is-loading="createSnippet.isPending.value"
    @close="showModal = false"
    @submit="handleCreateSnippet"
  />
</template>
