<template>
  <main class="mx-auto my-12 max-w-7xl px-4 lg:w-fit lg:min-w-4xl">
    <!-- Back button -->
    <Button variant="outline" @click="router.back()">
      <ArrowLeftIcon class="size-5" />
      Back
    </Button>

    <!-- Loading state -->
    <div v-if="isPending" class="mt-6 space-y-6">
      <!-- Snippet header skeleton -->
      <div class="rounded-lg border-2 bg-white p-6 shadow">
        <div class="mb-4 flex items-center justify-between">
          <div class="h-8 w-3/4 animate-pulse rounded-lg bg-gray-200"></div>
          <div class="h-10 w-20 animate-pulse rounded-lg bg-gray-200"></div>
        </div>
        <div class="h-6 w-1/4 animate-pulse rounded-lg bg-gray-200"></div>
      </div>

      <!-- Code block skeleton -->
      <div class="rounded-lg border-2 bg-white p-6 shadow">
        <div class="space-y-2">
          <div class="h-80 w-full animate-pulse rounded-lg bg-gray-200"></div>
        </div>
      </div>
    </div>

    <!-- Error state -->
    <div
      v-else-if="isError"
      class="border-destructive bg-destructive/10 mt-6 rounded-lg border-2 p-6"
    >
      <h2 class="text-destructive text-xl font-bold">Error</h2>
      <p class="text-destructive mt-2">{{ error?.message || 'An unexpected error occurred' }}</p>
    </div>

    <!-- Snippet content -->
    <div v-else-if="snippet" class="mt-6 space-y-6">
      <!-- Snippet header -->
      <div class="bg-card text-card-foreground rounded-lg border-2 p-6 shadow">
        <div class="mb-4 flex items-center justify-between">
          <h1 class="text-3xl font-bold">{{ snippet.title }}</h1>
          <LikeButton :likes="snippet.likes" :isLiked="snippet.isLiked" :snippetId="snippet.id" />
        </div>
        <p class="text-accent-foreground text-lg">By {{ snippet.author }}</p>
      </div>

      <!-- Code block -->
      <div class="bg-card rounded-lg border-2 p-6 shadow">
        <pre
          class="bg-muted overflow-x-auto rounded-lg p-4 font-mono text-sm"
        ><code>{{ snippet.content }}</code></pre>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import type { Snippet } from '@/types'
import { useQuery } from '@tanstack/vue-query'
import { useFetch } from '@/composables/useCustomFetch'
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import LikeButton from './_components/LikeButton.vue'
import { ArrowLeftIcon } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()

const getSnippet = async (): Promise<Snippet> => {
  const snippetId = route.params.snippetId as string
  const { data, error } = await useFetch<Snippet>(`/snippets/${snippetId}`, {
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

onMounted(() => {
  window.scrollTo({ top: 0, behavior: 'auto' })
})
</script>
