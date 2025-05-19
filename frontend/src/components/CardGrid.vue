<template>
  <div v-if="isLoading" class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
    <div v-for="i in 6" :key="i" class="animate-pulse">
      <div class="card flex max-h-[22rem] min-h-[22rem] flex-col p-6">
        <div class="mb-3 h-8 w-3/4 rounded bg-gray-200"></div>
        <div class="bg-muted grow overflow-hidden border-4 border-black p-4">
          <div class="h-full w-full rounded bg-gray-200"></div>
        </div>
        <div class="mt-4 flex items-center justify-between border-t-4 border-black pt-2">
          <div class="h-6 w-24 rounded bg-gray-200"></div>
          <div class="h-10 w-20 rounded bg-gray-200"></div>
        </div>
      </div>
    </div>
  </div>

  <div v-else-if="isError" class="flex flex-col items-center justify-center py-12 text-center">
    <div class="mb-4 rounded-full bg-red-100 p-4">
      <svg class="h-12 w-12 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
        />
      </svg>
    </div>
    <h3 class="mb-2 text-xl font-bold text-red-600">Something went wrong</h3>
    <p class="mb-4 text-gray-600">{{ errorMessage }}</p>
    <Button @click="$emit('retry')" variant="outline"> Try Again </Button>
  </div>

  <div v-else-if="isEmpty" class="flex flex-col items-center justify-center py-12 text-center">
    <div class="mb-4 rounded-full bg-gray-100 p-4">
      <svg class="h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
        />
      </svg>
    </div>
    <h3 class="mb-2 text-xl font-bold">No snippets found</h3>
    <p class="text-gray-600">Be the first to share your code snippet!</p>
  </div>

  <div v-else class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
    <div
      v-for="(card, index) in cards"
      :key="card.id"
      class="cursor-pointer transition-all duration-200 hover:-translate-y-1 hover:rotate-0"
      :class="{ 'rotate-1': index % 2 === 0, '-rotate-1': index % 2 !== 0 }"
    >
      <CodeCard
        :title="card.title"
        :content="card.content"
        :author="card.author"
        :likes="card.likes"
        @click="handleCardClick(card)"
        @toggle-like="() => handleLikeClick(card)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import CodeCard from './CodeCard.vue'
import { Button } from './ui/button'
import type { Card } from '@/models'

defineProps<{
  cards: Card[]
  isLoading?: boolean
  isEmpty?: boolean
  isError?: boolean
  errorMessage?: string
}>()

const emit = defineEmits<{
  (e: 'toggleLike', card: Card): void
  (e: 'retry'): void
}>()

const router = useRouter()

const handleCardClick = (card: Card) => {
  router.push(`/snippets/${card.id}`)
}

const handleLikeClick = (card: Card) => {
  emit('toggleLike', card)
}
</script>
