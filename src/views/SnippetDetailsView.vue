<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import AppHeader from '../components/AppHeader.vue'
import AppFooter from '../components/AppFooter.vue'
import { useToast } from '../composables/useToast'
import { Heart, ArrowLeft } from 'lucide-vue-next'
import { snippets } from '../data'

interface Snippet {
  id: number
  title: string
  code: string
  author: string
  likes: number
  isLiked: boolean
}

const route = useRoute()
const { showToast } = useToast()
const snippet = ref<Snippet | null>(null)

onMounted(() => {
  const snippetId = Number(route.params.snippetId)
  const foundSnippet = snippets.find((s) => s.id === snippetId)
  if (foundSnippet) {
    snippet.value = foundSnippet
  }
})

const toggleLike = () => {
  if (!snippet.value) return
  snippet.value.isLiked = !snippet.value.isLiked
  snippet.value.likes += snippet.value.isLiked ? 1 : -1
  showToast(
    snippet.value.isLiked ? 'Added to favorites' : 'Removed from favorites',
    snippet.value.isLiked ? 'success' : 'info',
  )
}
</script>

<template>
  <div class="bg-background min-h-screen text-black">
    <AppHeader />

    <main class="mx-auto my-12 max-w-4xl px-4">
      <div v-if="snippet" class="space-y-6">
        <!-- Back button -->
        <button
          class="inline-flex items-center gap-2 rounded-lg border-4 border-black bg-white px-4 py-2 font-bold shadow-[4px_4px_0_0_#000] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none"
          @click="$router.back()"
        >
          <ArrowLeft class="h-5 w-5" />
          Back
        </button>

        <!-- Snippet header -->
        <div class="rounded-lg border-4 border-black bg-white p-6 shadow-[8px_8px_0_0_#000]">
          <div class="mb-4 flex items-center justify-between">
            <h1 class="text-3xl font-bold">{{ snippet.title }}</h1>
            <button
              class="flex items-center gap-2 rounded-lg border-4 border-black px-4 py-2 font-bold shadow-[4px_4px_0_0_#000] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none"
              :class="{ 'bg-accent text-white': snippet.isLiked }"
              @click="toggleLike"
            >
              <Heart class="h-5 w-5" :class="{ 'fill-current': snippet.isLiked }" />
              {{ snippet.likes }}
            </button>
          </div>
          <p class="text-lg text-gray-600">By {{ snippet.author }}</p>
        </div>

        <!-- Code block -->
        <div class="rounded-lg border-4 border-black bg-white p-6 shadow-[8px_8px_0_0_#000]">
          <pre
            class="overflow-x-auto rounded-lg bg-gray-100 p-4 font-mono text-sm"
          ><code>{{ snippet.code }}</code></pre>
        </div>
      </div>
    </main>

    <AppFooter />
  </div>
</template>
