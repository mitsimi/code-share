<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { Heart, ArrowLeft } from 'lucide-vue-next'
import { useQuery } from '@tanstack/vue-query'
import { useCustomFetch } from '@/composables/useCustomFetch'
import { onMounted } from 'vue'

import { Button } from '../components/ui/button'
import { useLikeSnippet } from '@/composables/useLikeSnippet'
import type { Snippet } from '@/models'

const route = useRoute()
const router = useRouter()

const getSnippet = async (): Promise<Snippet> => {
  const snippetId = route.params.snippetId as string
  const { data, error } = await useCustomFetch<Snippet>(`/snippets/${snippetId}`, {
    timeout: 1000,
  }).json()

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
  const action = snippet.value.is_liked ? 'unlike' : 'like'
  console.log(`Toggling ${action} (${snippet.value.is_liked}) for snippet:`, snippet.value.id)
  updateLike({ snippetId: snippet.value.id, action })
}

onMounted(() => {
  window.scrollTo({ top: 0, behavior: 'auto' })
})
</script>

<template>
  <main class="mx-auto my-12 max-w-7xl px-4 lg:w-fit lg:min-w-4xl">
    <!-- Back button -->
    <Button variant="outline" @click="router.back()">
      <ArrowLeft class="size-5" />
      Back
    </Button>

    <!-- Loading state -->
    <div v-if="isPending" class="mt-6 space-y-6">
      <!-- Snippet header skeleton -->
      <div class="rounded-lg border-4 border-black bg-white p-6 shadow-[8px_8px_0_0_#000]">
        <div class="mb-4 flex items-center justify-between">
          <div class="h-8 w-3/4 animate-pulse rounded-lg bg-gray-200"></div>
          <div class="h-10 w-20 animate-pulse rounded-lg bg-gray-200"></div>
        </div>
        <div class="h-6 w-1/4 animate-pulse rounded-lg bg-gray-200"></div>
      </div>

      <!-- Code block skeleton -->
      <div class="rounded-lg border-4 border-black bg-white p-6 shadow-[8px_8px_0_0_#000]">
        <div class="space-y-2">
          <div class="h-80 w-full animate-pulse rounded-lg bg-gray-200"></div>
        </div>
      </div>
    </div>

    <!-- Error state -->
    <div
      v-else-if="isError"
      class="border-destructive bg-destructive/10 mt-6 rounded-lg border-4 p-6"
    >
      <h2 class="text-destructive text-xl font-bold">Error</h2>
      <p class="text-destructive mt-2">{{ error?.message || 'An unexpected error occurred' }}</p>
    </div>

    <!-- Snippet content -->
    <div v-else-if="snippet" class="mt-6 space-y-6">
      <!-- Snippet header -->
      <div class="rounded-lg border-4 border-black bg-white p-6 shadow-[8px_8px_0_0_#000]">
        <div class="mb-4 flex items-center justify-between">
          <h1 class="text-3xl font-bold">{{ snippet.title }}</h1>
          <Button variant="outline" @click="toggleLike">
            <Heart class="size-5" :fill="snippet.is_liked ? 'red' : 'none'" />
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
