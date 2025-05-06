<template>
  <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
    <div
      v-for="(card, index) in cards"
      :key="card.id"
      @click="handleCardClick(card)"
      class="cursor-pointer transition-all duration-200 hover:-translate-y-1 hover:rotate-0"
      :class="{ 'rotate-1': index % 2 === 0, '-rotate-1': index % 2 !== 0 }"
    >
      <CodeCard
        :title="card.title"
        :code="card.code"
        :author="card.author"
        :likes="card.likes"
        :is-liked="card.isLiked"
        @toggle-like="() => handleLikeClick(card)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import CodeCard from './CodeCard.vue'

interface Card {
  id: number
  title: string
  code: string
  author: string
  likes: number
  isLiked: boolean
}

defineProps<{
  cards: Card[]
}>()

const emit = defineEmits<{
  (e: 'toggleLike', card: Card): void
}>()

const router = useRouter()

const handleCardClick = (card: Card) => {
  router.push(`/snippets/${card.id}`)
}

const handleLikeClick = (card: Card) => {
  emit('toggleLike', card)
}
</script>
