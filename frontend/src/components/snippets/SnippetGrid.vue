<template>
  <!-- Loading State with Enhanced Skeletons -->
  <div v-if="isLoading" class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
    <div v-for="i in 6" :key="i" class="animate-pulse">
      <div class="bg-card rounded-lg border-2 p-6 shadow-sm">
        <!-- Header skeleton -->
        <div class="mb-4">
          <div class="mb-3 flex items-start justify-between gap-3">
            <div class="bg-muted h-6 w-3/4 rounded"></div>
            <div class="bg-muted h-5 w-16 rounded-full"></div>
          </div>
          <div class="flex items-center gap-2">
            <div class="bg-muted h-8 w-8 rounded-full"></div>
            <div class="space-y-1">
              <div class="bg-muted h-4 w-24 rounded"></div>
              <div class="bg-muted h-3 w-16 rounded"></div>
            </div>
          </div>
        </div>

        <!-- Code content skeleton -->
        <div class="mb-4">
          <div class="bg-muted/30 rounded-md border">
            <div class="bg-muted/50 flex items-center justify-between border-b px-3 py-2">
              <div class="flex items-center gap-2">
                <div class="flex gap-1">
                  <div class="bg-muted h-2 w-2 rounded-full"></div>
                  <div class="bg-muted h-2 w-2 rounded-full"></div>
                  <div class="bg-muted h-2 w-2 rounded-full"></div>
                </div>
                <div class="bg-muted h-3 w-20 rounded"></div>
              </div>
            </div>
            <div class="space-y-2 p-3">
              <div class="bg-muted h-3 w-full rounded"></div>
              <div class="bg-muted h-3 w-4/5 rounded"></div>
              <div class="bg-muted h-3 w-3/5 rounded"></div>
            </div>
          </div>
        </div>

        <!-- Footer skeleton -->
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <div class="bg-muted h-3 w-8 rounded"></div>
            <div class="bg-muted h-3 w-12 rounded"></div>
          </div>
          <div class="flex items-center gap-2">
            <div class="bg-muted h-8 w-8 rounded"></div>
            <div class="bg-muted h-8 w-16 rounded"></div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Error State -->
  <div v-else-if="isError" class="flex flex-col items-center justify-center py-16 text-center">
    <div class="bg-destructive/10 mb-6 rounded-full p-6">
      <AlertTriangle class="text-destructive h-16 w-16" />
    </div>
    <h3 class="text-foreground mb-3 text-2xl font-bold">Something went wrong</h3>
    <p class="text-muted-foreground mb-6 max-w-md">
      {{ errorMessage }}
    </p>
    <Button @click="$emit('retry')" variant="outline" size="lg" class="gap-2">
      <RotateCcw class="h-4 w-4" />
      Try Again
    </Button>
  </div>

  <!-- Empty State -->
  <div v-else-if="isEmpty" class="flex flex-col items-center justify-center py-16 text-center">
    <div class="bg-muted/50 mb-6 rounded-full p-6">
      <FileCode class="text-muted-foreground h-16 w-16" />
    </div>

    <h3 class="text-foreground mb-3 text-2xl font-bold">
      {{ emptyTitle }}
    </h3>

    <p class="text-muted-foreground mb-6 max-w-md">
      {{ emptyMessage }}
    </p>

    <Authenticated>
      <Button v-if="showCreateButton" @click="$emit('create-snippet')" size="lg" class="gap-2">
        <Plus class="h-4 w-4" />
        Share Your First Snippet
      </Button>
    </Authenticated>
  </div>

  <!-- Grid with Code Snippets -->
  <div v-else class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
    <div
      v-for="(snippet, index) in snippets"
      :key="snippet.id"
      class="cursor-pointer transition-all duration-300 hover:-translate-y-1 hover:scale-105 hover:rotate-0"
      :class="{
        'rotate-1': index % 3 === 0,
        '-rotate-1': index % 3 === 1,
        'rotate-0': index % 3 === 2,
      }"
    >
      <SnippetCard :snippet="snippet" @click="handleCardClick(snippet)" />
    </div>
  </div>
</template>

<script setup lang="ts">
import SnippetCard from './SnippetCard.vue'
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { AlertTriangle, RotateCcw, FileCode, Plus } from 'lucide-vue-next'
import type { Snippet } from '@/types'

withDefaults(
  defineProps<{
    snippets: Snippet[]
    isLoading?: boolean
    isEmpty?: boolean
    isError?: boolean
    errorMessage?: string
    emptyTitle?: string
    emptyMessage?: string
    showCreateButton?: boolean
  }>(),
  {
    emptyTitle: 'No code snippets found',
    emptyMessage:
      'Be the first to share your code snippet with the community! Start building something amazing.',
    errorMessage: 'Failed to load code snippets. Please check your connection and try again.',
    showCreateButton: true,
  },
)

defineEmits<{
  (e: 'retry'): void
  (e: 'create-snippet'): void
}>()

const router = useRouter()

const handleCardClick = (snippet: Snippet) => {
  router.push({ name: 'snippet-details', params: { snippetId: snippet.id } })
}
</script>
