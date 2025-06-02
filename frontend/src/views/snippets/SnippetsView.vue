<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { type Snippet } from '@/types'
import { useQueryClient, useQuery, useMutation } from '@tanstack/vue-query'
import { useFetch } from '@/composables/useCustomFetch'
import { ref } from 'vue'
import { toast } from 'vue-sonner'
import SnippetGrid from './_components/SnippetGrid.vue'
import SnippetModal from './_components/SnippetModal.vue'

const authStore = useAuthStore()

const showModal = ref(false)
const queryClient = useQueryClient()

const getSnippets = async (): Promise<Snippet[]> => {
  const { data, error } = await useFetch<Snippet[]>('/snippets', {
    timeout: 1000,
    afterFetch: (ctx) => {
      ctx.data = ctx.data.map((snippet: Snippet) => ({
        ...snippet,
        id: snippet.id,
      }))
      return ctx
    },
  }).json()

  if (error.value) {
    throw new Error('Failed to fetch snippets')
  }

  return data.value || []
}

const { isPending, isError, data, error, refetch } = useQuery({
  queryKey: ['snippets'],
  queryFn: getSnippets,
  staleTime: 1000 * 60, // Consider data fresh for 1 minute
})

// Show toast notification when an error occurs
/* watch(isError, (newIsError) => {
  if (newIsError) {
    toast.error('Failed to load snippets. Please try again.')
  }
}) */

const createSnippet = async (formData: { title: string; code: string }): Promise<Snippet> => {
  const { data, error } = await useFetch<Snippet>('/snippets', {
    method: 'POST',
    body: JSON.stringify({
      title: formData.title,
      content: formData.code,
    }),
  }).json()

  if (error.value) {
    throw new Error('Failed to create snippet')
  }

  return data.value!
}

const { mutate: submitSnippet, isPending: isSubmitting } = useMutation({
  mutationFn: createSnippet,
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
  onError: () => {
    toast.error('Failed to create snippet. Please try again.')
  },
})
</script>

<template>
  <main class="mx-auto max-w-7xl px-4">
    <SnippetGrid
      :cards="data || []"
      :is-loading="isPending"
      :is-empty="!isPending && (!data || data.length === 0)"
      :is-error="isError"
      :error-message="error?.message || 'An unexpected error occurred'"
      @retry="refetch"
      @create-snippet="showModal = true"
    />
  </main>

  <FloatingActionButton v-show="authStore.isAuthenticated()" @click="showModal = true" />

  <SnippetModal
    :show="showModal"
    :is-loading="isSubmitting"
    @close="showModal = false"
    @submit="submitSnippet"
  />
</template>
