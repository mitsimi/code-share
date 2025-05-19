<script setup lang="ts">
import { ref } from 'vue'
import CardGrid from '../components/CardGrid.vue'
import FloatingActionButton from '../components/FloatingActionButton.vue'
import SnippetModal from '../components/SnippetModal.vue'
import { useToast } from '../composables/useToast'
import { snippets } from '../data'
interface Card {
  id: number
  title: string
  code: string
  author: string
  likes: number
  isLiked: boolean
}

const showModal = ref(false)
const { showToast } = useToast()

const cards = ref<Card[]>(snippets)

const toggleLike = (card: Card) => {
  card.isLiked = !card.isLiked
  card.likes += card.isLiked ? 1 : -1
  showToast(
    card.isLiked ? 'Added to favorites' : 'Removed from favorites',
    card.isLiked ? 'success' : 'info',
  )
}

const submitSnippet = (formData: { title: string; code: string; author: string }) => {
  const newId = Math.max(...cards.value.map((card) => card.id)) + 1

  const newCard: Card = {
    id: newId,
    title: formData.title,
    code: formData.code,
    author: formData.author,
    likes: 0,
    isLiked: false,
  }

  cards.value.unshift(newCard)
  showModal.value = false
  showToast('"' + newCard.title + '" has been added successfully!', 'success')
}
</script>

<template>
  <main class="mx-auto my-12 max-w-7xl px-4">
    <CardGrid :cards="cards" @toggle-like="toggleLike" />
  </main>

  <FloatingActionButton @click="showModal = true" />

  <SnippetModal :show="showModal" @close="showModal = false" @submit="submitSnippet" />
</template>
