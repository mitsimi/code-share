<template>
  <Card class="max-h-[22rem] min-h-[22rem] cursor-pointer gap-0 p-6" @click="$emit('click')">
    <CardHeader class="mb-3 border-b-2 px-0 pb-2">
      <CardTitle class="text-xl font-bold break-words">
        {{ snippet.title }}
      </CardTitle>
    </CardHeader>
    <CardContent class="bg-muted grow overflow-hidden border-2 p-4">
      <pre class="m-0 p-0"><code class="text-sm leading-relaxed">{{ snippet.content }}</code></pre>
    </CardContent>
    <CardFooter class="mt-4 flex justify-between border-t-2 px-0 pt-2 font-bold">
      <span>{{ snippet.author }}</span>
      <LikeButton :likes="snippet.likes" :isLiked="snippet.isLiked" :snippetId="snippet.id" />
    </CardFooter>
  </Card>
</template>

<script setup lang="ts">
import LikeButton from '@/components/LikeButton.vue'
import type { Snippet } from '@/types'
import Card from './ui/card/Card.vue'
import CardHeader from './ui/card/CardHeader.vue'
import CardTitle from './ui/card/CardTitle.vue'
import CardContent from './ui/card/CardContent.vue'
import CardFooter from './ui/card/CardFooter.vue'

defineProps<{
  snippet: Snippet
}>()

defineEmits<{
  (e: 'click'): void
}>()
</script>
