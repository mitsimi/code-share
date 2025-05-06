<template>
  <div
    class="flex max-h-[22rem] cursor-pointer flex-col border-4 border-black bg-white p-6 shadow-[5px_5px_0_0_#000] transition-all duration-200 hover:-translate-y-1 hover:shadow-[8px_8px_0_0_#000]"
  >
    <h2
      class="overflow-wrap mb-3 border-b-3 border-black pb-2 text-xl font-black break-words text-black"
    >
      {{ title }}
    </h2>
    <div class="grow overflow-hidden border-3 border-black bg-white p-3">
      <pre class="m-0 p-0"><code class="text-sm leading-relaxed text-black">{{ code }}</code></pre>
    </div>
    <div
      class="mt-4 flex items-center justify-between border-t-3 border-black pt-2 font-bold text-black"
    >
      <span>{{ author }}</span>
      <div
        class="flex cursor-pointer items-center gap-2 border-3 border-black px-2.5 py-1 shadow-[2px_2px_0_0_#000]"
        :class="{
          'bg-accent text-white': isLiked,
          'bg-black text-white': likes < 0,
        }"
        @click="$emit('toggle-like')"
      >
        <span>{{ likes }}</span>
        <LucideHeart :fill="isLiked ? 'currentColor' : 'none'" class="size-4" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { LucideHeart } from 'lucide-vue-next'

defineProps<{
  title: string
  code: string
  author: string
  likes: number
  isLiked: boolean
  isOdd: boolean
}>()

defineEmits<{
  (e: 'toggle-like'): void
}>()
</script>
