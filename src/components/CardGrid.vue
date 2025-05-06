<template>
  <div class="grid grid-cols-1 gap-10 md:grid-cols-2 lg:grid-cols-3">
    <CodeCard
      v-for="(card, index) in cards"
      :key="card.id"
      :title="card.title"
      :code="card.code"
      :author="card.author"
      :likes="card.likes"
      :is-liked="card.isLiked"
      :is-odd="index % 2 === 0"
      @toggle-like="toggleLike(card)"
    />
  </div>
</template>

<script setup lang="ts">
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
  (e: 'toggle-like', card: Card): void
}>()

const toggleLike = (card: Card) => {
  emit('toggle-like', card)
}
</script>
