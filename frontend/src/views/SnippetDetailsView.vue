<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { Heart, ArrowLeft } from 'lucide-vue-next'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { useCustomFetch } from '@/composables/useCustomFetch'
import { onMounted } from 'vue'

import { Button } from '../components/ui/button'
import type { Card } from '@/models'
import { useLikeSnippet } from '@/composables/useLikeSnippet'

const route = useRoute()
const router = useRouter()
const queryClient = useQueryClient()

const getSnippet = async (): Promise<Card> => {
  const snippetId = route.params.snippetId as string
  const { data, error } = await useCustomFetch<Card>(`/snippets/${snippetId}`).json()

  if (error.value) {
    throw new Error('Failed to fetch snippet')
  }

  if (!data.value) {
    throw new Error('Snippet not found')
  }

  return data.value
}

const {
  data: snippet,
  isPending,
  isError,
  error,
} = useQuery({
  queryKey: ['snippet', route.params.snippetId],
  queryFn: getSnippet,
})

const { updateLike } = useLikeSnippet()

const toggleLike = () => {
  if (!snippet.value) {
    console.error('Cannot toggle like: snippet is null')
    return
  }
  const action = snippet.value.likes > 0 ? 'unlike' : 'like'
  console.log(`Toggling ${action} for snippet:`, snippet.value.id)
  updateLike({ snippetId: snippet.value.id, action })
}

onMounted(() => {
  window.scrollTo({ top: 0, behavior: 'auto' })
})
</script>

<template>
  <main class="mx-auto my-12 w-fit max-w-7xl min-w-4xl px-4">
    <!-- Back button -->
    <Button variant="outline" @click="router.back()">
      <ArrowLeft class="size-5" />
      Back
    </Button>

    <!-- Loading state -->
    <div v-if="isPending" class="mt-6 space-y-6">
      <div class="h-32 animate-pulse rounded-lg border-4 border-black bg-gray-100"></div>
      <div class="h-64 animate-pulse rounded-lg border-4 border-black bg-gray-100"></div>
    </div>

    <!-- Error state -->
    <div
      v-else-if="isError"
      class="border-destructive bg-destructive/10 mt-6 rounded-lg border-4 p-6"
    >
      <h2 class="text-destructive text-xl font-bold">Error</h2>
      <p class="text-destructive mt-2">{{ error?.message || 'An unexpected error occurred' }}</p>
      <Button class="mt-4" @click="router.back()">Go Back</Button>
    </div>

    <!-- Snippet content -->
    <div v-else-if="snippet" class="mt-6 space-y-6">
      <!-- Snippet header -->
      <div class="rounded-lg border-4 border-black bg-white p-6 shadow-[8px_8px_0_0_#000]">
        <div class="mb-4 flex items-center justify-between">
          <h1 class="text-3xl font-bold">{{ snippet.title }}</h1>
          <Button variant="outline" @click="toggleLike">
            <Heart class="size-5" :class="{ 'fill-current': snippet.likes > 0 }" />
            {{ snippet.likes }}
          </Button>
        </div>
        <p class="text-lg text-gray-600">By {{ snippet.author }}</p>
      </div>

      <!-- Code block -->
      <div class="rounded-lg border-4 border-black bg-white p-6 shadow-[8px_8px_0_0_#000]">
        <pre
          class="overflow-x-auto rounded-lg bg-gray-100 p-4 font-mono text-sm"
        ><code>{{ snippet.content }}</code></pre>
      </div>
    </div>
  </main>
</template>
