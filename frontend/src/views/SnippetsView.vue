<script setup lang="ts">
import { ref, watch } from 'vue'
import CardGrid from '@/components/CardGrid.vue'
import FloatingActionButton from '@/components/FloatingActionButton.vue'
import SnippetModal from '@/components/SnippetModal.vue'
import { useToast } from '@/composables/useToast'
import type { Card } from '@/models'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { useCustomFetch } from '@/composables/useCustomFetch'

const showModal = ref(false)
const { showToast } = useToast()
const queryClient = useQueryClient()

const getSnippets = async (): Promise<Card[]> => {
  const { data, error } = await useCustomFetch<Card[]>('/snippets', {
    afterFetch: (ctx) => {
      ctx.data = ctx.data.map((snippet: Card) => ({
        ...snippet,
        id: Number(snippet.id),
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
})

// Show toast notification when an error occurs
watch(isError, (newIsError) => {
  if (newIsError) {
    showToast('Failed to load snippets. Please try again.', 'error')
  }
})

const createSnippet = async (formData: {
  title: string
  code: string
  author: string
}): Promise<Card> => {
  const { data, error } = await useCustomFetch<Card>('/snippets', {
    method: 'POST',
    body: JSON.stringify({
      title: formData.title,
      content: formData.code,
      author: formData.author,
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
    queryClient.setQueryData(['snippets'], (oldData: Card[] | undefined) => {
      if (!oldData) return [newSnippet]
      return [newSnippet, ...oldData]
    })

    showModal.value = false
    showToast('"' + newSnippet.title + '" has been added successfully!', 'success')
  },
  onError: () => {
    showToast('Failed to create snippet. Please try again.', 'error')
  },
})
</script>

<template>
  <main class="mx-auto my-12 max-w-7xl px-4">
    <CardGrid
      :cards="data || []"
      :is-loading="isPending"
      :is-empty="!isPending && (!data || data.length === 0)"
      :is-error="isError"
      :error-message="error?.message || 'An unexpected error occurred'"
      @retry="refetch"
    />
  </main>

  <FloatingActionButton @click="showModal = true" />

  <SnippetModal
    :show="showModal"
    :is-loading="isSubmitting"
    @close="showModal = false"
    @submit="submitSnippet"
  />
</template>
