<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { type Snippet } from '@/types'
import { useQueryClient, useQuery, useMutation } from '@tanstack/vue-query'
import { snippetsService } from '@/services/snippets'
import { ref } from 'vue'
import { toast } from 'vue-sonner'
import SnippetModal from './_components/SnippetModal.vue'

const authStore = useAuthStore()

const showModal = ref(false)
const queryClient = useQueryClient()

const { isPending, isError, data, error, refetch } = useQuery({
  queryKey: ['snippets'],
  queryFn: snippetsService.getSnippets,
  staleTime: 1000 * 60, // Consider data fresh for 1 minute
})

const { mutate: submitSnippet, isPending: isSubmitting } = useMutation({
  mutationFn: (formData: { title: string; code: string }) =>
    snippetsService.createSnippet({
      title: formData.title,
      content: formData.code,
    }),
  onSuccess: (newSnippet) => {
    // Update the cache with the new snippet
    queryClient.setQueryData(['snippets'], (oldData: Snippet[] | undefined) => {
      if (!oldData) return [newSnippet]
      return [newSnippet, ...oldData]
    })

    showModal.value = false
    toast.success('Snippet added successfully', {
      description: `"${newSnippet.title}" has been added successfully!`,
    })
  },
  onError: (error) => {
    toast.error(
      error instanceof Error ? error.message : 'Failed to create snippet. Please try again.',
    )
  },
})
</script>

<template>
  <main class="container mx-auto max-w-7xl px-8">
    <SnippetGrid
      :snippets="data || []"
      :is-loading="isPending"
      :is-empty="!isPending && (!data || data.length === 0)"
      :is-error="isError"
      :error-message="error?.message || 'An unexpected error occurred'"
      @retry="refetch"
      @create-snippet="showModal = true"
    />
  </main>

  <Authenticated>
    <FloatingActionButton @click="showModal = true" />
  </Authenticated>

  <SnippetModal
    :show="showModal"
    :is-loading="isSubmitting"
    @close="showModal = false"
    @submit="submitSnippet"
  />
</template>
