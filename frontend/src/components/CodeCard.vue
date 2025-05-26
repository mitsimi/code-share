<template>
  <div
    class="card flex max-h-[22rem] min-h-[22rem] cursor-pointer flex-col p-6"
    @click="$emit('click')"
  >
    <h2 class="mb-3 border-b-4 border-black pb-2 text-xl font-bold break-words">
      {{ snippet.title }}
    </h2>
    <div class="bg-muted grow overflow-hidden border-4 border-black p-4">
      <pre class="m-0 p-0"><code class="text-sm leading-relaxed">{{ snippet.content }}</code></pre>
    </div>
    <div class="mt-4 flex items-center justify-between border-t-4 border-black pt-2 font-bold">
      <span>{{ snippet.author }}</span>
      <LikeButton 
        :likes="snippet.likes" 
        :is_liked="snippet.is_liked" 
        :snippetId="snippet.id" 
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import LikeButton from '@/components/LikeButton.vue'
import type { Snippet } from '@/types'

defineProps<{
  snippet: Snippet
}>()

defineEmits<{
  (e: 'click'): void
}>()
</script>
