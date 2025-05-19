<template>
  <div
    class="card flex max-h-[22rem] min-h-[22rem] cursor-pointer flex-col p-6"
    @click="$emit('click')"
  >
    <h2 class="mb-3 border-b-4 border-black pb-2 text-xl font-bold break-words">
      {{ title }}
    </h2>
    <div class="bg-muted grow overflow-hidden border-4 border-black p-4">
      <pre class="m-0 p-0"><code class="text-sm leading-relaxed">{{ content }}</code></pre>
    </div>
    <div class="mt-4 flex items-center justify-between border-t-4 border-black pt-2 font-bold">
      <span>{{ author }}</span>
      <Button
        variant="outline"
        :class="{
          'bg-black text-white': likes < 0,
        }"
        @click.stop="$emit('toggle-like')"
      >
        <span>{{ likes }}</span>
        <Heart class="size-5" />
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Heart } from 'lucide-vue-next'
import { Button } from './ui/button'
defineProps<{
  title: string
  content: string
  author: string
  likes: number
}>()

defineEmits<{
  (e: 'click'): void
  (e: 'toggle-like'): void
}>()
</script>
