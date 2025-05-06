<template>
  <div
    class="flex max-h-[22rem] min-h-[22rem] cursor-pointer flex-col rounded-lg border-4 border-black bg-white p-6 shadow-[8px_8px_0_0_#000] transition-all duration-200 hover:-translate-y-1 hover:shadow-[12px_12px_0_0_#000]"
  >
    <h2 class="mb-3 border-b-4 border-black pb-2 text-xl font-bold break-words text-black">
      {{ title }}
    </h2>
    <div class="grow overflow-hidden border-4 border-black bg-gray-100 p-4">
      <pre class="m-0 p-0"><code class="text-sm leading-relaxed text-black">{{ code }}</code></pre>
    </div>
    <div
      class="mt-4 flex items-center justify-between border-t-4 border-black pt-2 font-bold text-black"
    >
      <span>{{ author }}</span>
      <button
        class="flex items-center gap-2 rounded-lg border-4 border-black px-4 py-2 font-bold shadow-[4px_4px_0_0_#000] transition-all hover:translate-x-1 hover:translate-y-1 hover:shadow-none"
        :class="{
          'bg-accent text-white': isLiked,
          'bg-black text-white': likes < 0,
          'bg-white': !isLiked && likes >= 0,
        }"
        @click="$emit('toggle-like')"
      >
        <span>{{ likes }}</span>
        <Heart class="h-5 w-5" :class="{ 'fill-current': isLiked }" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Heart } from 'lucide-vue-next'

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
